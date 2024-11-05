package api

import (
	"blogi/internal/postgres"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type GetAuthorsReq struct {
	// TODO: Implement the pagination
}

type GetAuthorsRes []postgres.ListUsersPublicRow

func (api *API) getAuthors(c echo.Context) error {
	var body GetAuthorsReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request binding failed.")
	} else if err = c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request validation failed.")
	}

	data, err := api.DB.ListUsersPublic(c.Request().Context())
	if err != nil {
		c.Logger().Error("api.DB.ListUsersPublic", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	count, err := api.DB.ListUsersPublicCount(c.Request().Context())
	if err != nil {
		c.Logger().Error("api.DB.ListUsersPublicCount", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	c.Response().Header().Set("x-total-count", strconv.FormatInt(count, 10))
	return c.JSON(http.StatusOK, GetAuthorsRes(data))
}

type GetAuthorReq struct {
	Username string `param:"username" validate:"required"`
}

type GetAuthorRes postgres.GetUserPublicRow

func (api *API) getAuthor(c echo.Context) error {
	var body GetAuthorReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request binding failed.")
	} else if err = c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request validation failed.")
	}

	data, err := api.DB.GetUserPublic(c.Request().Context(), body.Username)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Author not found")
		}
		c.Logger().Error("api.DB.GetUserPublic", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	return c.JSON(http.StatusOK, GetAuthorRes(data))
}
