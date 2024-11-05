package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetProfileReq struct{}

type GetProfileRes struct{}

func (api *API) getProfile(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

type UpdateProfileReq struct{}

type UpdateProfileRes struct{}

func (api *API) updateProfile(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

type CreateBlogReq struct{}

type CreateBlogRes struct{}

func (api *API) createBlog(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

type UpdateBlogReq struct{}

type UpdateBlogRes struct{}

func (api *API) updateBlog(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}

type DeleteBlogReq struct{}

type DeleteBlogRes struct{}

func (api *API) deleteBlog(c echo.Context) error {
	return c.NoContent(http.StatusNotImplemented)
}
