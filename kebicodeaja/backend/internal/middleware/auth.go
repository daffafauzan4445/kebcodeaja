package middleware

import (
	"net/http"
	"strings"

	"github.com/daffafauzan4445/kebcodeaja/backend/internal/config"
	"github.com/daffafauzan4445/kebcodeaja/backend/internal/models"
	"github.com/daffafauzan4445/kebcodeaja/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token tidak ditemukan"})
			return
		}

		claims, err := utils.ParseToken(strings.TrimPrefix(header, "Bearer "))
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Token tidak valid atau kedaluwarsa"})
			return
		}

		var count int64
		if err := config.DB.Model(&models.User{}).Where("id = ?", claims.UserID).Count(&count).Error; err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "Gagal memeriksa akun"})
			return
		}
		if count == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "Akun tidak ditemukan, silakan login ulang"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func RoleRequired(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role := c.GetString("role")
		for _, r := range roles {
			if r == role {
				c.Next()
				return
			}
		}
		c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"message": "Akses ditolak untuk role ini"})
	}
}