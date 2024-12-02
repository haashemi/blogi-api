package api

import (
	"blogi/internal/postgres"
	"blogi/pkg/argon2id"
	"context"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type GetDashboardUsersRes []postgres.ListUsersRow

func (api *API) getDashboardUsers(c echo.Context) error {
	data, err := api.DB.ListUsers(c.Request().Context())
	if err != nil {
		c.Logger().Error("api.DB.ListUsers", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	count, err := api.DB.ListUsersCount(c.Request().Context())
	if err != nil {
		c.Logger().Error("api.DB.ListUsersCount", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	c.Response().Header().Set("x-total-count", strconv.FormatInt(count, 10))
	return c.JSON(http.StatusOK, GetDashboardUsersRes(data))
}

type GetDashboardUserReq struct {
	ID int64 `param:"id" validate:"required"`
}

type GetDashboardUserRes postgres.GetUserRow

func (api *API) getDashboardUser(c echo.Context) error {
	var body GetDashboardUserReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	data, err := api.DB.GetUser(c.Request().Context(), body.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		c.Logger().Error("api.DB.GetUser", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	return c.JSON(http.StatusOK, GetDashboardUserRes(data))
}

type UpdateDashboardUserReq struct {
	ID       int64   `param:"id" validate:"required"`
	FullName string  `json:"fullName" validate:"required,max=256"`
	Username string  `json:"username" validate:"required,max=32"`
	AboutMe  *string `json:"aboutMe" validate:"omitempty,max=1024"`
	Password *string `json:"password" validate:"required,min=6"`
	IsBanned bool    `json:"isBanned"`
}

func (api *API) updateDashboardUser(c echo.Context) error {
	var body UpdateDashboardUserReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	user, err := api.DB.GetUserCriticalAuthData(context.Background(), body.ID)
	if err != nil {
		c.Logger().Error("api.DB.GetUserPassword", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	if body.Password != nil {
		user.Password, err = argon2id.CreateHash(*body.Password, argon2id.DefaultParams)
		if err != nil {
			return echo.ErrInternalServerError.SetInternal(err)
		}
	}

	err = api.DB.UpdateUserFull(context.Background(), postgres.UpdateUserFullParams{
		ID:       body.ID,
		FullName: body.FullName,
		Username: body.Username,
		AboutMe:  body.AboutMe,
		Password: user.Password,
		IsBanned: body.IsBanned,
	})
	if err != nil {
		c.Logger().Error("api.DB.UpdateUserFull", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	return c.NoContent(http.StatusNoContent)
}
