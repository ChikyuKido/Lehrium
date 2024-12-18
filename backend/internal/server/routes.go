package server

import (
	"lehrium-backend/internal/controller"
	middlewares "lehrium-backend/internal/middleware"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8081"}
	config.AllowCredentials = true

	r := gin.Default()
	r.Use(cors.New(config))

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	user := api.Group("/user").Use(middlewares.Auth())
	teacher := api.Group("/teacher").Use(middlewares.Auth())

	r.GET("/health", s.healthHandler)

	auth.POST("/login", controller.LoginUser)
	auth.POST("/register", controller.RegisterUser)
	auth.POST("/verifyEmail", controller.VerifyEmail)
	auth.POST("/requestEmailVerification", middlewares.Auth(), controller.SendVerificationEmail)

	user.GET("/comment", nil)
	user.GET("/rate", nil)

	teacher.GET("/list", nil)
	teacher.GET("/:id", nil)

	return r
}
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}
