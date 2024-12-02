package api

import (
	"blogi/internal/postgres"
	"net/http"

	"github.com/labstack/echo/v4"
)

type GetProfileRes postgres.GetUserRow

func (api *API) getProfile(c echo.Context) error {
	claims := getClaims(c)

	user, err := api.DB.GetUser(c.Request().Context(), claims.UserID)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.JSON(http.StatusOK, GetProfileRes(user))
}

type UpdateProfileReq struct {
	FullName string  `json:"fullName" validate:"required,max=256"`
	Username string  `json:"username" validate:"required,max=32"`
	AboutMe  *string `json:"aboutMe" validate:"omitempty,max=1024"`
}

func (api *API) updateProfile(c echo.Context) error {
	var body UpdateProfileReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	claims := getClaims(c)

	err := api.DB.UpdateUser(c.Request().Context(), postgres.UpdateUserParams{
		ID:       claims.UserID,
		FullName: body.FullName,
		Username: body.Username,
		AboutMe:  body.AboutMe,
	})
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

type CreateBlogReq struct {
	Title   string `json:"title" validate:"required,max=1024"`
	Summary string `json:"summary" validate:"required,max=2048"`
	Content string `json:"content" validate:"required"`
}

type CreateBlogRes struct {
	ID int64 `json:"id"`
}

func (api *API) createBlog(c echo.Context) error {
	var body CreateBlogReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	claims := getClaims(c)

	id, err := api.DB.CreateBlog(c.Request().Context(), postgres.CreateBlogParams{
		AuthorID: claims.UserID,
		Title:    body.Title,
		Summary:  body.Summary,
		Content:  body.Content,
	})
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.JSON(http.StatusOK, CreateBlogRes{ID: id})
}

type UpdateBlogReq struct {
	ID      int64  `param:"id" validate:"required"`
	Title   string `json:"title" validate:"required,max=1024"`
	Summary string `json:"summary" validate:"required,max=2048"`
	Content string `json:"content" validate:"required"`
}

func (api *API) updateBlog(c echo.Context) error {
	var body UpdateBlogReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	err := api.DB.UpdateBlog(c.Request().Context(), postgres.UpdateBlogParams{
		ID:      body.ID,
		Title:   body.Title,
		Summary: body.Summary,
		Content: body.Content,
	})
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}

type DeleteBlogReq struct {
	ID int64 `param:"id" validate:"required"`
}

func (api *API) deleteBlog(c echo.Context) error {
	var body DeleteBlogReq
	if err := c.Bind(&body); err != nil {
		return echo.ErrBadRequest
	} else if err = c.Validate(body); err != nil {
		return echo.ErrBadRequest
	}

	err := api.DB.RemoveBlog(c.Request().Context(), body.ID)
	if err != nil {
		return echo.ErrInternalServerError.SetInternal(err)
	}

	return c.NoContent(http.StatusNoContent)
}
