package scheduler

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/Rayy1502/pir-backend/internal/repository"
	mssql "github.com/microsoft/go-mssqldb"
)

type NotificationScheduler struct {
	db          *sql.DB
	notifRepo   *repository.NotifikasiRepository
	stopChan    chan struct{}
	interval    time.Duration
}

func NewNotificationScheduler(db *sql.DB, notifRepo *repository.NotifikasiRepository, interval time.Duration) *NotificationScheduler {
	if interval <= 0 {
		interval = 1 * time.Hour // default evaluasi setiap 1 jam
	}
	return &NotificationScheduler{
		db:        db,
		notifRepo: notifRepo,
		stopChan:  make(chan struct{}),
		interval:  interval,
	}
}

// Start menjalankan background scheduler
func (s *NotificationScheduler) Start() {
	log.Println("🕒 Background Notification Scheduler started (interval:", s.interval, ")")
	
	// Jalankan evaluasi langsung pada saat startup
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		if err := s.Evaluate(ctx); err != nil {
			log.Printf("⚠️ Scheduler startup evaluation error: %v", err)
		}
	}()

	// Loop periodik
	go func() {
		ticker := time.NewTicker(s.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
				if err := s.Evaluate(ctx); err != nil {
					log.Printf("⚠️ Scheduler evaluation error: %v", err)
				}
				cancel()
			case <-s.stopChan:
				log.Println("🛑 Background Notification Scheduler stopped")
				return
			}
		}
	}()
}

func (s *NotificationScheduler) Stop() {
	close(s.stopChan)
}

// Evaluate memeriksa seluruh proyek dan menghasilkan notifikasi otomatis dual-level serta status keterlambatan
func (s *NotificationScheduler) Evaluate(ctx context.Context) error {
	// 1. Ambil batas waktu konfigurasi per unit kerja
	unitLimits := make(map[string]int)
	rowsPengaturan, err := s.db.QueryContext(ctx, `SELECT unit_kerja, batas_waktu_bulan FROM pengaturan_unit_kerja`)
	if err == nil {
		defer rowsPengaturan.Close()
		for rowsPengaturan.Next() {
			var unit string
			var limit int
			if err := rowsPengaturan.Scan(&unit, &limit); err == nil {
				unitLimits[strings.ToLower(strings.TrimSpace(unit))] = limit
			}
		}
		if err := rowsPengaturan.Err(); err != nil {
			log.Printf("⚠️ Scheduler: error saat iterasi pengaturan_unit_kerja: %v", err)
		}
	}

	// 2. Ambil seluruh proyek yang belum berstatus 'selesai'
	query := `
		SELECT id, nama_kegiatan, unit_kerja, tanggal_go_live, status_pir
		FROM proyek
		WHERE status_pir != 'selesai'
	`
	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer rows.Close()

	type ProyekCheck struct {
		ID            string
		NamaKegiatan  string
		UnitKerja     string
		TanggalGoLive time.Time
		StatusPIR     string
	}

	proyekList := make([]ProyekCheck, 0)
	for rows.Next() {
		var p ProyekCheck
		var pID mssql.UniqueIdentifier
		if err := rows.Scan(&pID, &p.NamaKegiatan, &p.UnitKerja, &p.TanggalGoLive, &p.StatusPIR); err == nil {
			p.ID = pID.String()
			proyekList = append(proyekList, p)
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	now := time.Now()
	for _, p := range proyekList {
		// Default masa percobaan 2 bulan pasca rilis
		tglPercobaanSelesai := p.TanggalGoLive.AddDate(0, 2, 0)

		// Batas toleransi dari unit kerja (default 1 bulan toleransi)
		toleransiBulan := 1
		unitKey := strings.ToLower(strings.TrimSpace(p.UnitKerja))
		if custom, ok := unitLimits[unitKey]; ok && custom > 2 {
			toleransiBulan = custom - 2
		}
		tglBatasAkhir := tglPercobaanSelesai.AddDate(0, toleransiBulan, 0)
		hMinus7 := tglBatasAkhir.AddDate(0, 0, -7)

		tglBatasStr := tglBatasAkhir.Format("02 January 2006")

		if now.After(tglBatasAkhir) {
			// Kasus A: Melewati batas akhir toleransi -> status otomatis 'terlambat'
			if p.StatusPIR != "terlambat" {
				_, _ = s.db.ExecContext(ctx, `UPDATE proyek SET status_pir = 'terlambat', updated_at = SYSUTCDATETIME() WHERE id = @id`, sql.Named("id", p.ID))
			}

			notif := &model.Notifikasi{
				ProyekID:  p.ID,
				Tipe:      "terlambat",
				Judul:     fmt.Sprintf("PIR Terlambat: %s", p.NamaKegiatan),
				Deskripsi: fmt.Sprintf("Proyek '%s' telah melewati batas waktu toleransi pelaksanaan PIR sejak %s. Segera buat dan lengkapi laporan PIR.", p.NamaKegiatan, tglBatasStr),
			}
			_ = s.notifRepo.CreateIfNotExists(ctx, notif)

		} else if now.After(hMinus7) {
			// Kasus B: Peringatan Eskalasi Toleransi (H-7 sebelum batas akhir)
			notif := &model.Notifikasi{
				ProyekID:  p.ID,
				Tipe:      "segera",
				Judul:     fmt.Sprintf("Peringatan Eskalasi Toleransi (H-7): %s", p.NamaKegiatan),
				Deskripsi: fmt.Sprintf("Batas waktu toleransi pelaksanaan PIR untuk proyek '%s' akan berakhir dalam 7 hari (%s). Mohon segera memproses formulir PIR.", p.NamaKegiatan, tglBatasStr),
			}
			_ = s.notifRepo.CreateIfNotExists(ctx, notif)

		} else if now.After(tglPercobaanSelesai) {
			// Kasus C: Masa adaptasi pengguna 2 bulan pasca rilis telah selesai -> waktu PIR dimulai (status 'pending')
			if p.StatusPIR == "belum-jatuh-tempo" {
				_, _ = s.db.ExecContext(ctx, `UPDATE proyek SET status_pir = 'pending', updated_at = SYSUTCDATETIME() WHERE id = @id`, sql.Named("id", p.ID))
			}

			notif := &model.Notifikasi{
				ProyekID:  p.ID,
				Tipe:      "waktu_pir",
				Judul:     fmt.Sprintf("Waktu PIR Dimulai: %s", p.NamaKegiatan),
				Deskripsi: fmt.Sprintf("Masa uji coba pengguna 2 bulan untuk proyek '%s' telah selesai. Formulir PIR kini siap diisi dan dievaluasi.", p.NamaKegiatan),
			}
			_ = s.notifRepo.CreateIfNotExists(ctx, notif)
		}
	}

	return nil
}
