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

	// TODO-> Get the user info from JWT
	var USER_ID int64

	if body.Password == body.NewPassword {
		return echo.NewHTTPError(http.StatusBadRequest, "New password shouldn't be same as the old password.")
	}

	user, err := api.DB.GetUserCriticalAuthData(c.Request().Context(), USER_ID)
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
		ID:       USER_ID,
		Password: newHashedPassword,
	})
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	token, err := api.createAuthToken(USER_ID, user.IsAdmin)
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

// TODO: Fill the fields and use config if necessary.
func (api *API) createAuthToken(userID int64, isAdmin bool) (string, error) {
	// Create the Claims
	claims := &JWTClaims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "",
			Subject:   "",
			Audience:  []string{},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenDuration)),
			NotBefore: &jwt.NumericDate{},
			IssuedAt:  &jwt.NumericDate{},
			ID:        "",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString(api.HMAC)
}

// TODO: Fill the fields and use config if necessary.
func (api *API) createAuthCookie(token string) *http.Cookie {
	return &http.Cookie{
		Name:  TokenCookieName,
		Value: token,
		// Quoted:      false,
		Path: "/",
		// Domain:      "",
		Expires: time.Now().Add(TokenDuration),
		// RawExpires:  "",
		// MaxAge:      0,
		// Secure:      false,
		HttpOnly: true,
		// SameSite:    0,
		// Partitioned: false,
		// Raw:         "",
		// Unparsed:    []string{},
	}
}

func (api *API) createAuthMiddleware() echo.MiddlewareFunc {
	// Configure middleware with the custom claims type
	config := echojwt.Config{
		SigningKey:    api.HMAC,
		TokenLookup:   "cookie:" + TokenCookieName,
		NewClaimsFunc: func(c echo.Context) jwt.Claims { return new(JWTClaims) },
	}

	return echojwt.WithConfig(config)
}

func (api *API) adminCheckMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(*JWTClaims)

		if !claims.IsAdmin {
			return echo.ErrForbidden
		}
		return next(c)
	}
}
