package main

import (
	"context"
	"log"

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

	tables := []string{
		"tindak_lanjut_bukti",
		"tindak_lanjut",
		"laporan_temuan",
		"laporan_pir",
		"survey_jawaban",
		"survey_pertanyaan",
		"pir_form",
		"proyek",
	}

	for _, tbl := range tables {
		res, err := database.ExecContext(ctx, "DELETE FROM "+tbl)
		if err != nil {
			log.Printf("Gagal hapus data dari %s: %v", tbl, err)
		} else {
			rows, _ := res.RowsAffected()
			log.Printf("Berhasil mengosongkan tabel %s (%d baris dihapus)", tbl, rows)
		}
	}

	log.Println("Pembersihan data selesai. Database proyek & PIR sekarang kosong.")
}
