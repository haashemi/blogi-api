package api

import (
	"blogi/internal/postgres"
	"blogi/pkg/argon2id"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	TokenCookieName string        = "blogi-token"
	TokenDuration   time.Duration = 365 * 24 * time.Hour
)

var (
	ClearTokenCookie *http.Cookie = &http.Cookie{Name: TokenCookieName, Value: "", Path: "/", HttpOnly: true}
)

type SignUpReq struct {
	FullName string `json:"full_name" validate:"required,max=265"`
	Username string `json:"username" validate:"required,max=32"`
	Password string `json:"password" validate:"required"`
}

type SignUpRes postgres.CreateUserRow

func (api *API) signUp(c echo.Context) error {
	var body SignUpReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	hashedPassword, err := argon2id.CreateHash(body.Password, argon2id.DefaultParams)
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	user, err := api.DB.CreateUser(c.Request().Context(), postgres.CreateUserParams{
		FullName: body.FullName,
		Username: body.Username,
		Password: hashedPassword,
	})
	if err != nil {
		return echo.ErrInternalServerError.WithInternal(err)
	}

	token, err := api.createAuthToken(user.ID, false)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}
	c.SetCookie(api.createAuthCookie(token))

	return c.JSON(http.StatusOK, SignUpRes(user))
}

type SignInReq struct {
	Username string `json:"username" validate:"required,max=32"`
	Password string `json:"password" validate:"required"`
}

type SignInRes struct {
	FullName string `json:"fullName"`
}

func (api *API) signIn(c echo.Context) error {
	var body SignInReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	user, err := api.DB.GetUserAuthData(c.Request().Context(), body.Username)
	if err != nil {
		return echo.ErrNotFound.SetInternal(err)
	}

	isMatch, err := argon2id.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	if isMatch {
		return echo.ErrNotFound
	}

	if user.IsBanned {
		return echo.NewHTTPError(http.StatusForbidden, "You're banned.")
	}

	token, err := api.createAuthToken(user.ID, user.IsAdmin)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}
	c.SetCookie(api.createAuthCookie(token))

	return c.JSON(http.StatusOK, SignInRes{FullName: user.FullName})
}

func (api *API) signOut(c echo.Context) error {
	c.SetCookie(ClearTokenCookie)
	return c.NoContent(http.StatusNoContent)
}

type ChangePasswordReq struct {
	Password    string `json:"password" validate:"required"`
	NewPassword string `json:"newPassword" validate:"required"`
}

func (api *API) changePassword(c echo.Context) error {
	var body ChangePasswordReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	claims := getClaims(c)

	if body.Password == body.NewPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "New password shouldn't be same as the old password.")
	}

	user, err := api.DB.GetUserCriticalAuthData(c.Request().Context(), claims.UserID)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	isMatch, err := argon2id.ComparePasswordAndHash(body.Password, user.Password)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	} else if !isMatch {
		return echo.NewHTTPError(http.StatusForbidden, "Password is incorrect.")
	}

	if user.IsBanned {
		c.SetCookie(ClearTokenCookie)
		return echo.NewHTTPError(http.StatusForbidden, "You're banned.")
	}

	newHashedPassword, err := argon2id.CreateHash(body.NewPassword, argon2id.DefaultParams)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	err = api.DB.UpdateUserPassword(c.Request().Context(), postgres.UpdateUserPasswordParams{
		ID:       claims.UserID,
		Password: newHashedPassword,
	})
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	token, err := api.createAuthToken(claims.UserID, user.IsAdmin)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}
	c.SetCookie(api.createAuthCookie(token))

	return echo.ErrNotImplemented
}

type JWTClaims struct {
	UserID  int64 `json:"userId"`
	IsAdmin bool  `json:"isAdmin"`
	jwt.RegisteredClaims
}

func (api *API) createAuthToken(userID int64, isAdmin bool) (string, error) {
	currentTime := time.Now()

	// Create the Claims
	claims := &JWTClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "blogi-api",
			Audience:  []string{"blogi-frontend"},
			ExpiresAt: jwt.NewNumericDate(currentTime.Add(TokenDuration)),
			NotBefore: jwt.NewNumericDate(currentTime),
			IssuedAt:  jwt.NewNumericDate(currentTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(api.HMAC))
}

func (api *API) createAuthCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:     TokenCookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   !api.IsDevBuild,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(TokenDuration.Seconds()),
		Expires:  time.Now().Add(TokenDuration),
	}
}

func (api *API) createAuthMiddleware() echo.MiddlewareFunc {
	config := echojwt.Config{
		ErrorHandler: func(c echo.Context, err error) error {
			if err == echojwt.ErrJWTMissing {
				return echo.ErrUnauthorized
			}
			return err
		},
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(JWTClaims) },
		SigningKey:    []byte(api.HMAC),
		TokenLookup:   "cookie:" + TokenCookieName,
	}

	return echojwt.WithConfig(config)
}

func (api *API) adminCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := getClaims(c)

		if !claims.IsAdmin {
			return echo.ErrForbidden
		}
		return next(c)
	}
}
