package main

import (
	"github.com/gin-gonic/gin"
	"github.com/thisgirlElan/jwt_auth/initializers"
	"github.com/thisgirlElan/jwt_auth/routes"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DbConnect()
	initializers.SyncDatabase()
}

func main() {
	r := gin.Default()

	// setup routes
	routes.UserRoutes(r)

	r.Run()
}
