package routes

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/daffafauzan4445/kebcodeaja/backend/internal/handlers"
	"github.com/daffafauzan4445/kebcodeaja/backend/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			if origin == os.Getenv("FRONTEND_URL") {
				return true
			}
			return strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:")
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", "./uploads")

	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"status": "ok", "app": "keBCodeAja API"})
		})

		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/auth/me", handlers.Me)

			teacher := protected.Group("/teacher")
			teacher.Use(middleware.RoleRequired("teacher"))
			{
				teacher.GET("/dashboard", handlers.TeacherDashboard)
				teacher.POST("/classes", handlers.CreateClass)
			}

			student := protected.Group("/student")
			student.Use(middleware.RoleRequired("student"))
			{
				student.GET("/dashboard", handlers.StudentDashboard)
				student.POST("/classes/join", handlers.JoinClass)
			}
		}
	}

	return r
}
