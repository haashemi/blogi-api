package api

import (
	"blogi/internal/postgres"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type GetBlogsRes []postgres.ListBlogsPublicRow

func (api *API) getBlogs(c echo.Context) error {
	data, err := api.DB.ListBlogsPublic(c.Request().Context())
	if err != nil {
		c.Logger().Error("api.DB.ListBlogsPublic", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	count, err := api.DB.ListBlogsPublicCount(c.Request().Context())
	if err != nil {
		c.Logger().Error("api.DB.ListBlogsPublicCount", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	c.Response().Header().Set("x-total-count", strconv.FormatInt(count, 10))
	return c.JSON(http.StatusOK, GetBlogsRes(data))
}

type GetBlogReq struct {
	ID int64 `param:"id" validate:"required"`
}

type GetBlogRes postgres.GetBlogPublicRow

func (api *API) getBlog(c echo.Context) error {
	var body GetBlogReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	data, err := api.DB.GetBlogPublic(c.Request().Context(), body.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Blog not found")
		}
		c.Logger().Error("api.DB.GetBlogPublic", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	return c.JSON(http.StatusOK, GetBlogRes(data))
}
