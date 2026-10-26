package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Rayy1502/pir-backend/internal/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type TindakLanjutRepository struct {
	DB *sql.DB
}

func NewTindakLanjutRepository(db *sql.DB) *TindakLanjutRepository {
	return &TindakLanjutRepository{DB: db}
}

func (r *TindakLanjutRepository) Create(ctx context.Context, req *model.CreateTindakLanjutRequest) (*model.TindakLanjut, error) {
	status := req.Status
	if status == "" {
		status = "open"
	}

	query := `
		INSERT INTO tindak_lanjut (
			laporan_id, temuan, pic, target_selesai, status
		)
		OUTPUT INSERTED.id, INSERTED.created_at, INSERTED.updated_at
		VALUES (
			@laporan_id, @temuan, @pic, @target_selesai, @status
		)
	`
	var id mssql.UniqueIdentifier
	var createdAt, updatedAt time.Time

	var targetVal any = nil
	if req.TargetSelesai != "" {
		targetVal = req.TargetSelesai
	}

	err := r.DB.QueryRowContext(ctx, query,
		sql.Named("laporan_id", req.LaporanID),
		sql.Named("temuan", req.Temuan),
		sql.Named("pic", req.PIC),
		sql.Named("target_selesai", targetVal),
		sql.Named("status", status),
	).Scan(&id, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}

	return r.GetByID(ctx, id.String())
}

func (r *TindakLanjutRepository) GetByID(ctx context.Context, id string) (*model.TindakLanjut, error) {
	query := `
		SELECT 
			t.id, t.laporan_id, t.temuan, ISNULL(t.pic, ''), t.target_selesai,
			t.status, t.created_at, t.updated_at,
			l.nomor_pir, l.proyek_id,
			p.nama_kegiatan, p.unit_kerja
		FROM tindak_lanjut t
		JOIN laporan_pir l ON t.laporan_id = l.id
		JOIN proyek p ON l.proyek_id = p.id
		WHERE t.id = @id
	`
	var item model.TindakLanjut
	var tID, lID, pID mssql.UniqueIdentifier
	var targetTime sql.NullTime

	err := r.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&tID, &lID, &item.Temuan, &item.PIC, &targetTime,
		&item.Status, &item.CreatedAt, &item.UpdatedAt,
		&item.NomorPIR, &pID,
		&item.NamaKegiatan, &item.UnitKerja,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	item.ID = tID.String()
	item.LaporanID = lID.String()
	item.ProyekID = pID.String()
	if targetTime.Valid {
		item.TargetSelesai = targetTime.Time.Format("2006-01-02")
	}

	// Bukti
	item.Bukti = r.loadBukti(ctx, item.ID)

	return &item, nil
}

func (r *TindakLanjutRepository) loadBukti(ctx context.Context, tindakLanjutID string) []model.TindakLanjutBukti {
	buktiList := make([]model.TindakLanjutBukti, 0)
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, tindak_lanjut_id, nama_file, url, ISNULL(ukuran_bytes, 0), uploaded_at
		FROM tindak_lanjut_bukti
		WHERE tindak_lanjut_id = @tl_id
		ORDER BY uploaded_at DESC
	`, sql.Named("tl_id", tindakLanjutID))
	if err != nil {
		return buktiList
	}
	defer rows.Close()

	for rows.Next() {
		var b model.TindakLanjutBukti
		var bID, tlID mssql.UniqueIdentifier
		var urlStr sql.NullString
		if err := rows.Scan(&bID, &tlID, &b.NamaFile, &urlStr, &b.UkuranBytes, &b.UploadedAt); err == nil {
			b.ID = bID.String()
			b.TindakLanjutID = tlID.String()
			if urlStr.Valid {
				b.URL = urlStr.String
			}
			buktiList = append(buktiList, b)
		}
	}
	if err := rows.Err(); err != nil {
		return buktiList
	}
	return buktiList
}

func (r *TindakLanjutRepository) GetAll(ctx context.Context, status, laporanID, search string) ([]model.TindakLanjut, error) {
	query := `
		SELECT 
			t.id, t.laporan_id, t.temuan, ISNULL(t.pic, ''), t.target_selesai,
			t.status, t.created_at, t.updated_at,
			l.nomor_pir, l.proyek_id,
			p.nama_kegiatan, p.unit_kerja
		FROM tindak_lanjut t
		JOIN laporan_pir l ON t.laporan_id = l.id
		JOIN proyek p ON l.proyek_id = p.id
		WHERE 1=1
	`
	var args []any
	if status != "" && status != "all" {
		query += ` AND t.status = @status`
		args = append(args, sql.Named("status", status))
	}
	if laporanID != "" {
		query += ` AND t.laporan_id = @laporan_id`
		args = append(args, sql.Named("laporan_id", laporanID))
	}
	if search != "" {
		query += ` AND (LOWER(t.temuan) LIKE @search OR LOWER(t.pic) LIKE @search OR LOWER(p.nama_kegiatan) LIKE @search)`
		args = append(args, sql.Named("search", "%"+search+"%"))
	}
	query += ` ORDER BY t.created_at DESC`

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.TindakLanjut, 0)
	for rows.Next() {
		var item model.TindakLanjut
		var tID, lID, pID mssql.UniqueIdentifier
		var targetTime sql.NullTime

		err := rows.Scan(
			&tID, &lID, &item.Temuan, &item.PIC, &targetTime,
			&item.Status, &item.CreatedAt, &item.UpdatedAt,
			&item.NomorPIR, &pID,
			&item.NamaKegiatan, &item.UnitKerja,
		)
		if err != nil {
			return nil, err
		}

		item.ID = tID.String()
		item.LaporanID = lID.String()
		item.ProyekID = pID.String()
		if targetTime.Valid {
			item.TargetSelesai = targetTime.Time.Format("2006-01-02")
		}

		item.Bukti = r.loadBukti(ctx, item.ID)
		list = append(list, item)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return list, nil
}

func (r *TindakLanjutRepository) Update(ctx context.Context, id string, req *model.UpdateTindakLanjutRequest) error {
	var targetVal any = nil
	if req.TargetSelesai != "" {
		targetVal = req.TargetSelesai
	}

	query := `
		UPDATE tindak_lanjut
		SET temuan = @temuan,
		    pic = @pic,
		    target_selesai = @target_selesai,
		    status = @status,
		    updated_at = SYSUTCDATETIME()
		WHERE id = @id
	`
	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("id", id),
		sql.Named("temuan", req.Temuan),
		sql.Named("pic", req.PIC),
		sql.Named("target_selesai", targetVal),
		sql.Named("status", req.Status),
	)
	if err != nil {
		return err
	}

	// Cek apakah seluruh tindak lanjut pada laporan ini sudah 'done'
	r.checkAndUpdateProyekStatus(ctx, id)

	return nil
}

func (r *TindakLanjutRepository) checkAndUpdateProyekStatus(ctx context.Context, tindakLanjutID string) {
	// Ambil laporan_id dan proyek_id
	var lapID, prID mssql.UniqueIdentifier
	err := r.DB.QueryRowContext(ctx, `
		SELECT l.id, l.proyek_id
		FROM tindak_lanjut t
		JOIN laporan_pir l ON t.laporan_id = l.id
		WHERE t.id = @id
	`, sql.Named("id", tindakLanjutID)).Scan(&lapID, &prID)
	if err != nil {
		return
	}

	// Hitung total dan open/in_progress
	var total, nonDone int
	err = r.DB.QueryRowContext(ctx, `
		SELECT COUNT(*), ISNULL(SUM(CASE WHEN status != 'done' THEN 1 ELSE 0 END), 0)
		FROM tindak_lanjut
		WHERE laporan_id = @lap_id
	`, sql.Named("lap_id", lapID)).Scan(&total, &nonDone)
	if err != nil {
		return
	}

	if total == 0 || nonDone == 0 {
		// Semua selesai atau tidak ada temuan -> Proyek PIR Selesai
		_, _ = r.DB.ExecContext(ctx, `
			UPDATE proyek
			SET status_pir = 'selesai', updated_at = SYSUTCDATETIME()
			WHERE id = @proyek_id
		`, sql.Named("proyek_id", prID))

		_, _ = r.DB.ExecContext(ctx, `
			UPDATE f
			SET f.status = 'selesai', f.updated_at = SYSUTCDATETIME()
			FROM pir_form f
			WHERE f.proyek_id = @proyek_id
		`, sql.Named("proyek_id", prID))
	}
}

func (r *TindakLanjutRepository) AddBukti(ctx context.Context, tindakLanjutID, namaFile, url string, ukuranBytes int64) (*model.TindakLanjutBukti, error) {
	query := `
		INSERT INTO tindak_lanjut_bukti (tindak_lanjut_id, nama_file, url, ukuran_bytes)
		OUTPUT INSERTED.id, INSERTED.uploaded_at
		VALUES (@tl_id, @nama_file, @url, @ukuran)
	`
	var id mssql.UniqueIdentifier
	var uploadedAt time.Time

	err := r.DB.QueryRowContext(ctx, query,
		sql.Named("tl_id", tindakLanjutID),
		sql.Named("nama_file", namaFile),
		sql.Named("url", url),
		sql.Named("ukuran", ukuranBytes),
	).Scan(&id, &uploadedAt)
	if err != nil {
		return nil, err
	}

	return &model.TindakLanjutBukti{
		ID:             id.String(),
		TindakLanjutID: tindakLanjutID,
		NamaFile:       namaFile,
		URL:            url,
		UkuranBytes:    ukuranBytes,
		UploadedAt:     uploadedAt,
	}, nil
}
