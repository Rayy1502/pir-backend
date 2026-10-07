package main

import (
	"context"
	"log"

	"github.com/Rayy1502/pir-backend/internal/config"
	"github.com/Rayy1502/pir-backend/internal/db"
	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
	"github.com/Rayy1502/pir-backend/pkg/hash"
)

func main() {
	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database)

	// Cek apakah user sudah ada
	existing, _ := userRepo.FindByNIP(context.Background(), "19850312")
	if existing != nil {
		log.Println("ℹ️  User Budi Santoso sudah ada, skip seed")
		return
	}

	// Hash password
	hashed, err := hash.HashPassword("password123")
	if err != nil {
		log.Fatalf("Hash password gagal: %v", err)
	}

	user := &model.User{
		NIP:          "19850312",
		Nama:         "Budi Santoso",
		PasswordHash: hashed,
		Jabatan:      "Staff Divisi TI",
		Divisi:       "Teknologi Informasi",
		Unit:         "PT PUSRI Palembang",
		Role:         "staff_ti",
		IsActive:     true,
	}

	if err := userRepo.Create(context.Background(), user); err != nil {
		log.Fatalf("Insert user gagal: %v", err)
	}

	log.Printf("✅ Seed user berhasil: %s (NIP: %s), ID: %s\n", user.Nama, user.NIP, user.ID)
	log.Println("   Password: password123")
}
