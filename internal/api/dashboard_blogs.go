package api

import (
	"blogi/internal/postgres"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
)

type GetDashboardBlogsReq struct {
	// TODO: Implement the pagination
}

type GetDashboardBlogsRes []postgres.ListBlogsRow

func (api *API) getDashboardBlogs(c echo.Context) error {
	var body GetDashboardBlogsReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request binding failed.")
	} else if err = c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request validation failed.")
	}

	data, err := api.DB.ListBlogs(c.Request().Context())
	if err != nil {
		c.Logger().Error("api.DB.ListBlogs", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	count, err := api.DB.ListBlogsCount(c.Request().Context())
	if err != nil {
		c.Logger().Error("api.DB.ListBlogsCount", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	c.Response().Header().Set("x-total-count", strconv.FormatInt(count, 10))
	return c.JSON(http.StatusOK, GetDashboardBlogsRes(data))
}

type GetDashboardBlogReq struct {
	ID int64 `param:"id" validate:"required"`
}

type GetDashboardBlogRes postgres.Blog

func (api *API) getDashboardBlog(c echo.Context) error {
	var body GetDashboardBlogReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request binding failed.")
	} else if err = c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request validation failed.")
	}

	data, err := api.DB.GetBlog(c.Request().Context(), body.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Blog not found")
		}
		c.Logger().Error("api.DB.GetBlog", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	return c.JSON(http.StatusOK, GetDashboardBlogRes(data))
}

type UpdateDashboardBlogReq struct {
	ID      int64  `json:"id" validate:"required"`
	Title   string `json:"title" validate:"required,max=1024"`
	Summary string `json:"summary" validate:"required,max=2048"`
	Content string `json:"content" validate:"required"`
}

func (api *API) updateDashboardBlog(c echo.Context) error {
	var body UpdateDashboardBlogReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request binding failed.")
	} else if err = c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request validation failed.")
	}

	err := api.DB.UpdateBlog(c.Request().Context(), postgres.UpdateBlogParams(body))
	if err != nil {
		c.Logger().Error("api.DB.UpdateBlog", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	return c.NoContent(http.StatusNotImplemented)
}

type DeleteDashboardBlogReq struct {
	ID int64 `param:"id" validate:"required"`
}

func (api *API) deleteDashboardBlog(c echo.Context) error {
	var body DeleteDashboardBlogReq
	if err := c.Bind(&body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request binding failed.")
	} else if err = c.Validate(body); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Request validation failed.")
	}

	err := api.DB.RemoveBlog(c.Request().Context(), body.ID)
	if err != nil {
		c.Logger().Error("api.DB.RemoveBlog", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Database error.")
	}

	return c.NoContent(http.StatusNotImplemented)
}
