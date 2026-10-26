package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Rayy1502/pir-backend/internal/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type LaporanRepository struct {
	DB *sql.DB
}

func NewLaporanRepository(db *sql.DB) *LaporanRepository {
	return &LaporanRepository{DB: db}
}

func (r *LaporanRepository) GenerateNomorPIR(ctx context.Context) (string, error) {
	year := time.Now().Year()
	pattern := fmt.Sprintf("PIR/TI/%d/%%", year)

	var count int
	err := r.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM laporan_pir WHERE nomor_pir LIKE @pattern
	`, sql.Named("pattern", pattern)).Scan(&count)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("PIR/TI/%d/%03d", year, count+1), nil
}

func (r *LaporanRepository) Create(ctx context.Context, req *model.CreateLaporanRequest, pembuatID string) (*model.LaporanPIR, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	nomorPIR, err := r.GenerateNomorPIR(ctx)
	if err != nil {
		return nil, err
	}

	tglLaporan := req.TanggalLaporan
	if tglLaporan == "" {
		tglLaporan = time.Now().Format("2006-01-02")
	}

	var pembuatVal any = nil
	if pembuatID != "" {
		pembuatVal = pembuatID
	}

	query := `
		INSERT INTO laporan_pir (
			nomor_pir, proyek_id, pir_form_id, ringkasan,
			kepuasan, kesesuaian_tujuan, rekomendasi,
			tanggal_laporan, pembuat_id
		)
		OUTPUT INSERTED.id, INSERTED.created_at
		VALUES (
			@nomor_pir, @proyek_id, @pir_form_id, @ringkasan,
			@kepuasan, @kesesuaian_tujuan, @rekomendasi,
			@tanggal_laporan, @pembuat_id
		)
	`
	var lapID mssql.UniqueIdentifier
	var createdAt time.Time

	err = tx.QueryRowContext(ctx, query,
		sql.Named("nomor_pir", nomorPIR),
		sql.Named("proyek_id", req.ProyekID),
		sql.Named("pir_form_id", req.PIRFormID),
		sql.Named("ringkasan", req.Ringkasan),
		sql.Named("kepuasan", req.Kepuasan),
		sql.Named("kesesuaian_tujuan", req.KesesuaianTujuan),
		sql.Named("rekomendasi", req.Rekomendasi),
		sql.Named("tanggal_laporan", tglLaporan),
		sql.Named("pembuat_id", pembuatVal),
	).Scan(&lapID, &createdAt)
	if err != nil {
		return nil, err
	}

	lapIDStr := lapID.String()

	// Insert Temuan
	temuanList := make([]model.LaporanTemuan, 0, len(req.Temuan))
	for i, t := range req.Temuan {
		if t == "" {
			continue
		}
		var tID mssql.UniqueIdentifier
		err := tx.QueryRowContext(ctx, `
			INSERT INTO laporan_temuan (laporan_id, temuan, urutan)
			OUTPUT INSERTED.id
			VALUES (@laporan_id, @temuan, @urutan)
		`,
			sql.Named("laporan_id", lapIDStr),
			sql.Named("temuan", t),
			sql.Named("urutan", i+1),
		).Scan(&tID)
		if err != nil {
			return nil, err
		}

		temuanList = append(temuanList, model.LaporanTemuan{
			ID:        tID.String(),
			LaporanID: lapIDStr,
			Temuan:    t,
			Urutan:    i + 1,
		})
	}

	// Update pir_form status menjadi menunggu_ttd (belum selesai sampai di-ACC & upload PDF TTD)
	_, _ = tx.ExecContext(ctx, `
		UPDATE pir_form
		SET status = 'menunggu_ttd', updated_at = SYSUTCDATETIME()
		WHERE id = @pir_form_id
	`, sql.Named("pir_form_id", req.PIRFormID))

	// Status proyek tetap 'pending' sampai upload bukti PDF bertanda tangan
	_, _ = tx.ExecContext(ctx, `
		UPDATE proyek
		SET status_pir = 'pending', updated_at = SYSUTCDATETIME()
		WHERE id = @proyek_id
	`, sql.Named("proyek_id", req.ProyekID))

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, lapIDStr)
}

func (r *LaporanRepository) SaveSignedPDF(ctx context.Context, laporanID, fileURL, fileName string) (*model.LaporanPIR, error) {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// 1. Update laporan_pir dengan file PDF bertanda tangan & status approved
	queryLap := `
		UPDATE laporan_pir
		SET file_ttd_url = @url,
		    file_ttd_nama = @nama,
		    status_laporan = 'approved',
		    signed_at = SYSUTCDATETIME()
		WHERE id = @id
	`
	_, err = tx.ExecContext(ctx, queryLap,
		sql.Named("url", fileURL),
		sql.Named("nama", fileName),
		sql.Named("id", laporanID),
	)
	if err != nil {
		return nil, err
	}

	// 2. Ambil proyek_id dan pir_form_id dari laporan_pir
	var proyekID, pirFormID mssql.UniqueIdentifier
	err = tx.QueryRowContext(ctx, `SELECT proyek_id, pir_form_id FROM laporan_pir WHERE id = @id`, sql.Named("id", laporanID)).Scan(&proyekID, &pirFormID)
	if err != nil {
		return nil, err
	}

	// 3. Proyek RESMI SELESAI setelah upload bukti PDF bertanda tangan
	_, err = tx.ExecContext(ctx, `
		UPDATE proyek
		SET status_pir = 'selesai', updated_at = SYSUTCDATETIME()
		WHERE id = @proyek_id
	`, sql.Named("proyek_id", proyekID))
	if err != nil {
		return nil, err
	}

	// 4. Update status pir_form menjadi 'selesai'
	_, _ = tx.ExecContext(ctx, `
		UPDATE pir_form
		SET status = 'selesai', updated_at = SYSUTCDATETIME()
		WHERE id = @pir_form_id
	`, sql.Named("pir_form_id", pirFormID))

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return r.GetByID(ctx, laporanID)
}

func (r *LaporanRepository) GetByID(ctx context.Context, id string) (*model.LaporanPIR, error) {
	query := `
		SELECT 
			l.id, l.nomor_pir, l.proyek_id, l.pir_form_id, l.ringkasan,
			l.kepuasan, l.kesesuaian_tujuan, l.rekomendasi, l.tanggal_laporan,
			l.pembuat_id, ISNULL(u.nama, 'Staff TI') AS pembuat_nama, l.created_at,
			l.file_ttd_url, l.file_ttd_nama, ISNULL(l.status_laporan, 'draft'), l.signed_at,
			p.id, p.nama_kegiatan, p.kategori, p.unit_kerja, ISNULL(p.vendor, ''), p.tanggal_go_live, p.status_pir,
			f.metode, f.tujuan_awal
		FROM laporan_pir l
		JOIN proyek p ON l.proyek_id = p.id
		LEFT JOIN pir_form f ON l.pir_form_id = f.id
		LEFT JOIN users u ON l.pembuat_id = u.id
		WHERE l.id = @id
	`
	var lap model.LaporanPIR
	var lapID, prID, pfID mssql.UniqueIdentifier
	var pembuatUID, fileTTDURL, fileTTDNama sql.NullString
	var statusLaporan string
	var signedAt sql.NullTime
	var pID mssql.UniqueIdentifier
	var p model.Proyek
	var tglLap time.Time
	var metode, tujuanAwal sql.NullString

	err := r.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&lapID, &lap.NomorPIR, &prID, &pfID, &lap.Ringkasan,
		&lap.Kepuasan, &lap.KesesuaianTujuan, &lap.Rekomendasi, &tglLap,
		&pembuatUID, &lap.PembuatNama, &lap.CreatedAt,
		&fileTTDURL, &fileTTDNama, &statusLaporan, &signedAt,
		&pID, &p.NamaKegiatan, &p.Kategori, &p.UnitKerja, &p.Vendor, &p.TanggalGoLive, &p.StatusPIR,
		&metode, &tujuanAwal,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	lap.ID = lapID.String()
	lap.ProyekID = prID.String()
	lap.PIRFormID = pfID.String()
	lap.TanggalLaporan = tglLap.Format("2006-01-02")
	lap.StatusLaporan = statusLaporan
	if pembuatUID.Valid {
		lap.PembuatID = &pembuatUID.String
	}
	if fileTTDURL.Valid && fileTTDURL.String != "" {
		lap.FileTTDURL = &fileTTDURL.String
	}
	if fileTTDNama.Valid && fileTTDNama.String != "" {
		lap.FileTTDNama = &fileTTDNama.String
	}
	if signedAt.Valid {
		lap.SignedAt = &signedAt.Time
	}

	p.ID = pID.String()
	lap.Proyek = &p

	var pirForm model.PIRForm
	pirForm.ID = lap.PIRFormID
	pirForm.ProyekID = lap.ProyekID
	if metode.Valid {
		pirForm.Metode = metode.String
	}
	if tujuanAwal.Valid {
		pirForm.TujuanAwal = tujuanAwal.String
	}
	lap.PIRForm = &pirForm

	// Load temuan
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, temuan, urutan
		FROM laporan_temuan
		WHERE laporan_id = @lap_id
		ORDER BY urutan ASC
	`, sql.Named("lap_id", lap.ID))
	if err == nil {
		defer rows.Close()
		for rows.Next() {
			var t model.LaporanTemuan
			var tID mssql.UniqueIdentifier
			if err := rows.Scan(&tID, &t.Temuan, &t.Urutan); err == nil {
				t.ID = tID.String()
				t.LaporanID = lap.ID
				lap.Temuan = append(lap.Temuan, t)
			}
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}

	return &lap, nil
}

func (r *LaporanRepository) GetByProyekID(ctx context.Context, proyekID string) (*model.LaporanPIR, error) {
	var lapID mssql.UniqueIdentifier
	err := r.DB.QueryRowContext(ctx, `
		SELECT TOP 1 id FROM laporan_pir WHERE proyek_id = @proyek_id ORDER BY created_at DESC
	`, sql.Named("proyek_id", proyekID)).Scan(&lapID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.GetByID(ctx, lapID.String())
}

func (r *LaporanRepository) GetAll(ctx context.Context, proyekID, search string) ([]model.LaporanPIR, error) {
	query := `
		SELECT 
			l.id, l.nomor_pir, l.proyek_id, l.pir_form_id, l.ringkasan,
			l.kepuasan, l.kesesuaian_tujuan, l.rekomendasi, l.tanggal_laporan,
			ISNULL(u.nama, 'Staff TI') AS pembuat_nama, l.created_at,
			l.file_ttd_url, ISNULL(l.status_laporan, 'draft'),
			p.nama_kegiatan, p.kategori, p.unit_kerja, ISNULL(p.vendor, ''), p.tanggal_go_live, p.status_pir
		FROM laporan_pir l
		JOIN proyek p ON l.proyek_id = p.id
		LEFT JOIN users u ON l.pembuat_id = u.id
		WHERE 1=1
	`
	var args []any
	if proyekID != "" {
		query += ` AND l.proyek_id = @proyek_id`
		args = append(args, sql.Named("proyek_id", proyekID))
	}
	if search != "" {
		query += ` AND (LOWER(l.nomor_pir) LIKE @search OR LOWER(p.nama_kegiatan) LIKE @search)`
		args = append(args, sql.Named("search", "%"+search+"%"))
	}
	query += ` ORDER BY l.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.LaporanPIR, 0)
	for rows.Next() {
		var lap model.LaporanPIR
		var lapID, prID, pfID mssql.UniqueIdentifier
		var tglLap time.Time
		var p model.Proyek
		var fileTTDURL sql.NullString
		var statusLaporan string

		err := rows.Scan(
			&lapID, &lap.NomorPIR, &prID, &pfID, &lap.Ringkasan,
			&lap.Kepuasan, &lap.KesesuaianTujuan, &lap.Rekomendasi, &tglLap,
			&lap.PembuatNama, &lap.CreatedAt,
			&fileTTDURL, &statusLaporan,
			&p.NamaKegiatan, &p.Kategori, &p.UnitKerja, &p.Vendor, &p.TanggalGoLive, &p.StatusPIR,
		)
		if err != nil {
			return nil, err
		}

		lap.ID = lapID.String()
		lap.ProyekID = prID.String()
		lap.PIRFormID = pfID.String()
		lap.TanggalLaporan = tglLap.Format("2006-01-02")
		lap.StatusLaporan = statusLaporan
		if fileTTDURL.Valid && fileTTDURL.String != "" {
			lap.FileTTDURL = &fileTTDURL.String
		}
		p.ID = lap.ProyekID
		lap.Proyek = &p

		list = append(list, lap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}
