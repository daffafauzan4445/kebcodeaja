package handlers

import (
	"net/http"
	"strings"

	"github.com/daffafauzan4445/kebcodeaja/backend/internal/config"
	"github.com/daffafauzan4445/kebcodeaja/backend/internal/dto"
	"github.com/daffafauzan4445/kebcodeaja/backend/internal/models"
	"github.com/daffafauzan4445/kebcodeaja/backend/internal/utils"
	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data tidak valid. Pastikan semua kolom terisi dan password minimal 8 karakter"})
		return
	}

	email := strings.ToLower(strings.TrimSpace(req.Email))

	var count int64
	config.DB.Model(&models.User{}).Where("email = ?", email).Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "Email sudah terdaftar"})
		return
	}

	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal memproses password"})
		return
	}

	user := models.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		PasswordHash: hash,
		Role:         models.UserRole(req.Role),
	}
	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat akun"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Registrasi berhasil", "user": user})
}

func Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Email dan password wajib diisi"})
		return
	}

	var user models.User
	err := config.DB.Where("email = ?", strings.ToLower(strings.TrimSpace(req.Email))).First(&user).Error
	if err != nil || !utils.CheckPassword(user.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Email atau password salah"})
		return
	}

	token, err := utils.GenerateToken(user.ID, string(user.Role))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login berhasil", "token": token, "user": user})
}

func Me(c *gin.Context) {
	userID := c.GetUint("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": "User tidak ditemukan"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}
