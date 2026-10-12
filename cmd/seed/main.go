package main

import (
	"context"
	"database/sql"
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

	ctx := context.Background()

	// ============================================
	// Seed User
	// ============================================
	userRepo := repository.NewUserRepository(database)

	existing, _ := userRepo.FindByNIP(ctx, "19850312")
	if existing != nil {
		log.Println("ℹ️  User Budi Santoso sudah ada, skip seed")
	} else {
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

		if err := userRepo.Create(ctx, user); err != nil {
			log.Fatalf("Insert user gagal: %v", err)
		}

		log.Printf("✅ Seed user: %s (NIP: %s)\n", user.Nama, user.NIP)
	}

	// ============================================
	// Seed Proyek
	// ============================================
	type proyekSeed struct {
		Nama     string
		Kategori string
		Unit     string
		Vendor   string
		GoLive   string
		Status   string
	}

	proyekList := []proyekSeed{
		{"Sistem e-Procurement v2", "Aplikasi", "Infrastruktur", "Tim Internal Aplikasi", "2026-05-12", "terlambat"},
		{"Upgrade Jaringan WAN Pabrik", "Infrastruktur", "Infrastruktur", "PT Sigma Cipta Caraka", "2026-05-18", "terlambat"},
		{"Migrasi Server Virtualisasi", "Infrastruktur", "Infrastruktur", "Tim Internal Infrastruktur", "2026-07-15", "pending"},
		{"Portal HSE Mobile", "Aplikasi", "Digitalisasi", "PT Digital Nusantara", "2026-04-10", "selesai"},
		{"Aplikasi Monitoring Produksi", "Aplikasi", "Digitalisasi", "Tim Internal Aplikasi", "2026-03-03", "selesai"},
		{"Pengadaan Laptop & Thin Client", "EUC", "Layanan", "PT Lenovo Indonesia", "2026-08-28", "belum-jatuh-tempo"},
	}

	for _, p := range proyekList {
		var count int
		err := database.QueryRowContext(ctx,
			"SELECT COUNT(*) FROM proyek WHERE nama_kegiatan = @nama",
			sql.Named("nama", p.Nama),
		).Scan(&count)
		if err != nil {
			log.Printf("❌ Cek proyek '%s' gagal: %v", p.Nama, err)
			continue
		}
		if count > 0 {
			log.Printf("ℹ️  Proyek '%s' sudah ada, skip", p.Nama)
			continue
		}

		query := `
			INSERT INTO proyek (nama_kegiatan, kategori, unit_kerja, vendor, tanggal_go_live, status_pir)
			VALUES (@nama, @kategori, @unit, @vendor, @golive, @status)
		`
		_, err = database.ExecContext(ctx, query,
			sql.Named("nama", p.Nama),
			sql.Named("kategori", p.Kategori),
			sql.Named("unit", p.Unit),
			sql.Named("vendor", p.Vendor),
			sql.Named("golive", p.GoLive),
			sql.Named("status", p.Status),
		)
		if err != nil {
			log.Printf("❌ Insert proyek '%s' gagal: %v", p.Nama, err)
			continue
		}
		log.Printf("✅ Seed proyek: %s", p.Nama)
	}

	log.Println("🎉 Seed selesai")
}
