package main

import (
	"log"
	"os"

	"github.com/daffafauzan4445/kebcodeaja/backend/internal/config"
	"github.com/daffafauzan4445/kebcodeaja/backend/internal/routes"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env tidak ditemukan, memakai environment sistem")
	}

	if os.Getenv("JWT_SECRET") == "" {
		log.Println("JWT_SECRET kosong, memakai secret bawaan (hanya untuk development)")
		os.Setenv("JWT_SECRET", "kebicodeaja-dev-secret-ubah-di-produksi-1234567890")
	}

	config.ConnectDatabase()

	r := routes.SetupRouter()

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(r.Run(":" + port))
}
