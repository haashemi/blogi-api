package api

import (
	"blogi/internal/config"
	"blogi/internal/postgres"
	"blogi/pkg/validate"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type API struct{ APIConfig }

type APIConfig struct {
	config.APIConfig
	DB *postgres.Connection
}

func Run(conf APIConfig) {
	e := echo.New()
	e.Validator = validate.NewEchoValidator()

	e.Use(middleware.Recover())
	e.Use(middleware.Logger())
	e.Use(middleware.Secure())
	e.Use(middleware.CORS())
	e.Use(middleware.CSRF())
	e.Use(middleware.Gzip())

	api := API{conf}
	authMiddleware := api.createAuthMiddleware()

	auth := e.Group("/auth")
	{
		auth.POST("/sign-up", api.signUp)
		auth.POST("/sign-in", api.signIn)
		auth.POST("/sign-out", api.signOut)
		auth.POST("/change-password", api.changePassword, authMiddleware)
	}

	dashboard := e.Group("/api/dashboard", authMiddleware, api.adminCheckMiddleware)
	{
		blog := dashboard.Group("/blogs")
		{
			blog.GET("", api.getDashboardBlogs)
			blog.GET("/:id", api.getDashboardBlog)
			blog.PATCH("/:id", api.updateDashboardBlog)
			blog.DELETE("/:id", api.deleteDashboardBlog)
		}

		users := dashboard.Group("/users")
		{
			users.GET("", api.getDashboardUsers)
			users.GET("/:id", api.getDashboardUser)
			users.PATCH("/:id", api.updateDashboardUser)
		}
	}

	profile := e.Group("/profile", authMiddleware)
	{
		profile.GET("", api.getProfile)
		profile.PATCH("", api.updateProfile)

		blog := profile.Group("/blog")
		{
			blog.POST("", api.createBlog)
			blog.PATCH("/:id", api.updateBlog)
			blog.DELETE("/:id", api.deleteBlog)
		}
	}

	public := e.Group("/api/public")
	{
		authors := public.Group("/authors")
		{
			authors.GET("", api.getAuthors)
			authors.GET("/:username", api.getAuthor)
		}

		blogs := public.Group("/blogs")
		{
			blogs.GET("", api.getBlogs)
			blogs.GET("/:id", api.getBlog)
		}
	}

	e.Logger.Fatal(e.Start(conf.APIAddr))
}
