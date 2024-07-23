package routes

import (
	"github.com/Ilya-Semakin/fh/controllers"
	"github.com/Ilya-Semakin/fh/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(r *gin.Engine) {
	r.POST("/register", controllers.Register)
	r.POST("/login", controllers.Login)

	r.GET("/auth/vk/login", controllers.OAuthVKLogin)
	r.GET("/auth/vk/callback", controllers.OAuthVKCallback)
	r.GET("/auth/google/login", controllers.OAuthGoogleLogin)
	r.GET("/auth/google/callback", controllers.OAuthGoogleCallback)

	auth := r.Group("/auth")
	auth.Use(middlewares.JWTAuthMiddleware())
	{
		// тут можно маршруты ебануть
	}
}
