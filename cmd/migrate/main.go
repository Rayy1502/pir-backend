package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Rayy1502/pir-backend/internal/config"
	"github.com/Rayy1502/pir-backend/internal/db"
)

func main() {
	cfg := config.Load()
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB connection failed: %v", err)
	}
	defer database.Close()

	ctx := context.Background()

	// 1. Cek tabel-tabel yang ada saat ini
	rows, err := database.QueryContext(ctx, "SELECT name FROM sys.tables ORDER BY name")
	if err != nil {
		log.Fatalf("Query sys.tables gagal: %v", err)
	}
	defer rows.Close()

	fmt.Println("📋 Tabel saat ini di database:")
	for rows.Next() {
		var name string
		_ = rows.Scan(&name)
		fmt.Printf(" - %s\n", name)
	}
	if err := rows.Err(); err != nil {
		log.Printf("Peringatan saat iterasi tabel: %v", err)
	}

	// 2. Jalankan migrasi 2 jika pengaturan_unit_kerja belum ada
	var countPengaturan int
	_ = database.QueryRowContext(ctx, "SELECT COUNT(*) FROM sys.tables WHERE name = 'pengaturan_unit_kerja'").Scan(&countPengaturan)
	if countPengaturan == 0 {
		fmt.Println("\n⏳ Menjalankan migrasi 000002_pengaturan_per_unit.up.sql...")
		m2Path := filepath.Join("migrations", "000002_pengaturan_per_unit.up.sql")
		sql2, err := os.ReadFile(m2Path)
		if err != nil {
			log.Fatalf("Gagal membaca %s: %v", m2Path, err)
		}
		if _, err := database.ExecContext(ctx, string(sql2)); err != nil {
			log.Fatalf("Gagal eksekusi migrasi 2: %v", err)
		}
		fmt.Println("✅ Migrasi 000002 berhasil dijalankan!")
	} else {
		fmt.Println("ℹ️  Tabel pengaturan_unit_kerja sudah ada.")
	}

	// 3. Jalankan migrasi 3 (workflow approval & signed pdf)
	fmt.Println("\n⏳ Menjalankan migrasi 000003_workflow_approval_signed_pdf.up.sql...")
	m3Path := filepath.Join("migrations", "000003_workflow_approval_signed_pdf.up.sql")
	sql3, err := os.ReadFile(m3Path)
	if err != nil {
		log.Fatalf("Gagal membaca %s: %v", m3Path, err)
	}
	if _, err := database.ExecContext(ctx, string(sql3)); err != nil {
		log.Fatalf("Gagal eksekusi migrasi 3: %v", err)
	}
	fmt.Println("✅ Migrasi 000003 berhasil dijalankan!")

	// 4. Jalankan migrasi 4 (multiple survey respondents)
	fmt.Println("\n⏳ Menjalankan migrasi 000004_multiple_survey_respondents.up.sql...")
	m4Path := filepath.Join("migrations", "000004_multiple_survey_respondents.up.sql")
	sql4, err := os.ReadFile(m4Path)
	if err != nil {
		log.Fatalf("Gagal membaca %s: %v", m4Path, err)
	}
	if _, err := database.ExecContext(ctx, string(sql4)); err != nil {
		log.Fatalf("Gagal eksekusi migrasi 4: %v", err)
	}
	fmt.Println("✅ Migrasi 000004 berhasil dijalankan!")

	// 5. Jalankan migrasi 5 (pir_dokumentasi)
	fmt.Println("\n⏳ Menjalankan migrasi 000005_pir_dokumentasi.up.sql...")
	m5Path := filepath.Join("migrations", "000005_pir_dokumentasi.up.sql")
	sql5, err := os.ReadFile(m5Path)
	if err != nil {
		log.Fatalf("Gagal membaca %s: %v", m5Path, err)
	}
	if _, err := database.ExecContext(ctx, string(sql5)); err != nil {
		log.Fatalf("Gagal eksekusi migrasi 5: %v", err)
	}
	fmt.Println("✅ Migrasi 000005 berhasil dijalankan!")

	// 6. Verifikasi kolom di pir_form dan laporan_pir
	fmt.Println("\n🔍 Kolom di pir_form:")
	cols, err := database.QueryContext(ctx, `
		SELECT c.name, t.name as type_name 
		FROM sys.columns c 
		JOIN sys.types t ON c.user_type_id = t.user_type_id 
		WHERE c.object_id = OBJECT_ID('pir_form')
	`)
	if err == nil {
		defer cols.Close()
		for cols.Next() {
			var colName, typeName string
			_ = cols.Scan(&colName, &typeName)
			fmt.Printf("   - %s (%s)\n", colName, typeName)
		}
		if err := cols.Err(); err != nil {
			log.Printf("Peringatan saat membaca kolom pir_form: %v", err)
		}
	}

	fmt.Println("\n🔍 Kolom di laporan_pir:")
	cols2, err := database.QueryContext(ctx, `
		SELECT c.name, t.name as type_name 
		FROM sys.columns c 
		JOIN sys.types t ON c.user_type_id = t.user_type_id 
		WHERE c.object_id = OBJECT_ID('laporan_pir')
	`)
	if err == nil {
		defer cols2.Close()
		for cols2.Next() {
			var colName, typeName string
			_ = cols2.Scan(&colName, &typeName)
			fmt.Printf("   - %s (%s)\n", colName, typeName)
		}
		if err := cols2.Err(); err != nil {
			log.Printf("Peringatan saat membaca kolom laporan_pir: %v", err)
		}
	}

	fmt.Println("\n🎉 Database migration selesai dengan sukses!")
}
