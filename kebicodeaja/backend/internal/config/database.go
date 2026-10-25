package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/daffafauzan4445/kebcodeaja/backend/internal/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func ensureDatabase(user, pass, host, port, name string) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local", user, pass, host, port)
	conn, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Exec("CREATE DATABASE IF NOT EXISTS `" + name + "` CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci")
	return err
}

func ConnectDatabase() {
	user := getenv("DB_USER", "root")
	pass := os.Getenv("DB_PASS")
	host := getenv("DB_HOST", "127.0.0.1")
	port := getenv("DB_PORT", "3306")
	name := getenv("DB_NAME", "kebcodeaja")

	if err := ensureDatabase(user, pass, host, port, name); err != nil {
		log.Fatalf("Gagal koneksi ke MySQL (pastikan MySQL di XAMPP/Laragon sudah menyala): %v", err)
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, pass, host, port, name)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Class{},
		&models.ClassMember{},
		&models.Material{},
		&models.MaterialAttachment{},
		&models.Challenge{},
		&models.Submission{},
	)
	if err != nil {
		log.Fatalf("Gagal auto-migrate: %v", err)
	}

	DB = db
	log.Println("Database terhubung dan migrasi selesai")
}
