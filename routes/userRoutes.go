package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/thisgirlElan/jwt_auth/controllers"
	"github.com/thisgirlElan/jwt_auth/middleware"
)

var UserRoutes = func(r *gin.Engine) {
	auth := r.Group("/auth")
	{
		auth.POST("/register", controllers.RegisterHandler)
		auth.POST("/login", controllers.LoginHandler)
	}

	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/users", controllers.UsersHandler)
		api.POST("/user", controllers.UserByIdHandler)
		api.GET("/profile", controllers.UserHandler)
		api.PUT("/profile/update/:id", controllers.UpdateUserHandler)
		api.DELETE("/profile/delete/:id", controllers.DeleteProfile)
	}
}
