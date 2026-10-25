package handlers

import (
	"crypto/rand"
	"errors"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"

	"github.com/daffafauzan4445/kebcodeaja/backend/internal/config"
	"github.com/daffafauzan4445/kebcodeaja/backend/internal/dto"
	"github.com/daffafauzan4445/kebcodeaja/backend/internal/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const codeCharset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"

func debugError(err error) string {
	if os.Getenv("APP_ENV") == "development" {
		return err.Error()
	}
	return ""
}

func generateCode() (string, error) {
	b := make([]byte, 4)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeCharset))))
		if err != nil {
			return "", err
		}
		b[i] = codeCharset[n.Int64()]
	}
	return "KBA-" + string(b), nil
}

func uniqueCode() (string, error) {
	for i := 0; i < 10; i++ {
		code, err := generateCode()
		if err != nil {
			return "", err
		}
		var count int64
		if err := config.DB.Model(&models.Class{}).Where("class_code = ?", code).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return code, nil
		}
	}
	return "", errors.New("gagal membuat kode kelas yang unik")
}

func CreateClass(c *gin.Context) {
	var req dto.CreateClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Nama kelas wajib diisi, minimal 3 karakter"})
		return
	}

	name := strings.TrimSpace(req.Name)
	if len([]rune(name)) < 3 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Nama kelas wajib diisi, minimal 3 karakter"})
		return
	}

	code, err := uniqueCode()
	if err != nil {
		log.Printf("gagal membuat kode kelas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat kode kelas", "error": debugError(err)})
		return
	}

	class := models.Class{
		TeacherID:   c.GetUint("user_id"),
		Name:        name,
		Description: strings.TrimSpace(req.Description),
		ClassCode:   code,
	}
	if err := config.DB.Create(&class).Error; err != nil {
		log.Printf("gagal menyimpan kelas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal membuat kelas", "error": debugError(err)})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Kelas berhasil dibuat", "class": class})
}

func JoinClass(c *gin.Context) {
	var req dto.JoinClassRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Kode kelas wajib diisi"})
		return
	}

	code := strings.ToUpper(strings.TrimSpace(req.ClassCode))
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Kode kelas wajib diisi"})
		return
	}

	var class models.Class
	if err := config.DB.Where("class_code = ?", code).First(&class).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"message": "Kode kelas tidak ditemukan"})
			return
		}
		log.Printf("gagal mencari kelas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal mencari kelas", "error": debugError(err)})
		return
	}

	studentID := c.GetUint("user_id")

	var count int64
	config.DB.Model(&models.ClassMember{}).
		Where("class_id = ? AND student_id = ?", class.ID, studentID).
		Count(&count)
	if count > 0 {
		c.JSON(http.StatusConflict, gin.H{"message": "Kamu sudah tergabung di kelas ini"})
		return
	}

	member := models.ClassMember{ClassID: class.ID, StudentID: studentID}
	if err := config.DB.Create(&member).Error; err != nil {
		log.Printf("gagal bergabung ke kelas: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal bergabung ke kelas", "error": debugError(err)})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Berhasil bergabung ke kelas " + class.Name, "class": class})
}
