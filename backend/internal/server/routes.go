package server

import (
    "lehrium-backend/internal/controller"
    "lehrium-backend/internal/middleware"
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

    auth := r.Group("/auth")
    test := r.Group("/test").Use(middlewares.Auth())
    user := r.Group("/user").Use(middlewares.Auth())
    teacher :=r.Group("/teacher").Use(middlewares.Auth())

	r.GET("/health", s.healthHandler)
    
    auth.POST("/login", controllers.LoginUser)
    auth.POST("/register", controllers.RegisterUser)

    user.GET("/comment", nil)
    user.GET("/rate", nil)

    teacher.GET("/list", nil)
    teacher.GET("/:id", nil)

    test.GET("/ping", s.pong)
	return r
}
func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) pong(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{"message": "pong"})
}
