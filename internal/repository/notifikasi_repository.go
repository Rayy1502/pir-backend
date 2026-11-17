package repository

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/Rayy1502/pir-backend/internal/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type NotifikasiRepository struct {
	DB *sql.DB
}

func NewNotifikasiRepository(db *sql.DB) *NotifikasiRepository {
	return &NotifikasiRepository{DB: db}
}

func (r *NotifikasiRepository) GetAll(ctx context.Context, filter model.NotifikasiFilter) ([]model.Notifikasi, error) {
	query := `
		SELECT 
			n.id, n.proyek_id, n.tipe, ISNULL(n.judul, ''), ISNULL(n.deskripsi, ''), n.dibaca, n.created_at,
			ISNULL(p.nama_kegiatan, ''), ISNULL(p.unit_kerja, ''), ISNULL(p.tanggal_go_live, SYSUTCDATETIME())
		FROM notifikasi n
		LEFT JOIN proyek p ON n.proyek_id = p.id
		WHERE 1=1
	`
	var args []any
	if filter.Dibaca != nil {
		query += " AND n.dibaca = @dibaca"
		args = append(args, sql.Named("dibaca", *filter.Dibaca))
	}
	query += " ORDER BY n.dibaca ASC, n.created_at DESC"

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	list := make([]model.Notifikasi, 0)
	for rows.Next() {
		var n model.Notifikasi
		var nID, prID mssql.UniqueIdentifier
		var tglRilis time.Time

		if err := rows.Scan(
			&nID, &prID, &n.Tipe, &n.Judul, &n.Deskripsi, &n.Dibaca, &n.CreatedAt,
			&n.NamaKegiatan, &n.UnitKerja, &tglRilis,
		); err != nil {
			return nil, err
		}

		n.ID = nID.String()
		n.ProyekID = prID.String()
		n.TanggalRilis = tglRilis.Format("2006-01-02")
		list = append(list, n)
	}

	return list, rows.Err()
}

func (r *NotifikasiRepository) GetUnreadCount(ctx context.Context) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM notifikasi WHERE dibaca = 0`
	err := r.DB.QueryRowContext(ctx, query).Scan(&count)
	return count, err
}

func (r *NotifikasiRepository) MarkAsRead(ctx context.Context, id string) error {
	query := `UPDATE notifikasi SET dibaca = 1 WHERE id = @id`
	_, err := r.DB.ExecContext(ctx, query, sql.Named("id", id))
	return err
}

func (r *NotifikasiRepository) MarkAllAsRead(ctx context.Context) error {
	query := `UPDATE notifikasi SET dibaca = 1 WHERE dibaca = 0`
	_, err := r.DB.ExecContext(ctx, query)
	return err
}

func (r *NotifikasiRepository) CreateIfNotExists(ctx context.Context, notif *model.Notifikasi) error {
	// Periksa apakah notifikasi untuk proyek ini dan tipe yang sama sudah ada dan belum dibaca
	var existingID mssql.UniqueIdentifier
	checkQuery := `
		SELECT TOP 1 id FROM notifikasi 
		WHERE proyek_id = @proyek_id AND tipe = @tipe AND dibaca = 0
	`
	err := r.DB.QueryRowContext(ctx, checkQuery,
		sql.Named("proyek_id", notif.ProyekID),
		sql.Named("tipe", notif.Tipe),
	).Scan(&existingID)

	if err == nil {
		// Sudah ada notifikasi aktif untuk proyek dan tipe ini, tidak perlu insert duplikat
		return nil
	} else if err != sql.ErrNoRows {
		return err
	}

	insertQuery := `
		INSERT INTO notifikasi (proyek_id, tipe, judul, deskripsi, dibaca)
		VALUES (@proyek_id, @tipe, @judul, @deskripsi, 0)
	`
	_, err = r.DB.ExecContext(ctx, insertQuery,
		sql.Named("proyek_id", notif.ProyekID),
		sql.Named("tipe", notif.Tipe),
		sql.Named("judul", strings.TrimSpace(notif.Judul)),
		sql.Named("deskripsi", strings.TrimSpace(notif.Deskripsi)),
	)
	return err
}
