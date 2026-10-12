package repository

import (
	"context"
	"database/sql"
	"strings"

	"github.com/Rayy1502/pir-backend/internal/dto"
	"github.com/Rayy1502/pir-backend/internal/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type ProyekRepository struct {
	DB *sql.DB
}

func NewProyekRepository(db *sql.DB) *ProyekRepository {
	return &ProyekRepository{DB: db}
}

func scanProyek(rows *sql.Rows) ([]model.Proyek, error) {
	var list []model.Proyek
	for rows.Next() {
		var p model.Proyek
		var uid mssql.UniqueIdentifier
		var vendor sql.NullString

		err := rows.Scan(
			&uid,
			&p.NamaKegiatan,
			&p.Kategori,
			&p.UnitKerja,
			&vendor,
			&p.TanggalGoLive,
			&p.StatusPIR,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		p.ID = uid.String()
		if vendor.Valid {
			p.Vendor = vendor.String
		}

		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *ProyekRepository) GetAll(ctx context.Context, f dto.ProyekFilter) ([]model.Proyek, error) {
	query := `
		SELECT id, nama_kegiatan, kategori, unit_kerja, vendor, tanggal_go_live,
		       status_pir, created_at, updated_at
		FROM proyek
		WHERE 1=1
	`

	var args []any

	if f.Query != "" {
		query += ` AND LOWER(nama_kegiatan) LIKE @q`
		args = append(args, sql.Named("q", "%"+strings.ToLower(f.Query)+"%"))
	}
	if f.Kategori != "" {
		query += ` AND kategori = @kategori`
		args = append(args, sql.Named("kategori", f.Kategori))
	}
	if f.Status != "" {
		query += ` AND status_pir = @status`
		args = append(args, sql.Named("status", f.Status))
	}

	query += ` ORDER BY tanggal_go_live DESC`

	rows, err := r.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanProyek(rows)
}

func (r *ProyekRepository) GetByID(ctx context.Context, id string) (*model.Proyek, error) {
	query := `
		SELECT id, nama_kegiatan, kategori, unit_kerja, vendor, tanggal_go_live,
		       status_pir, created_at, updated_at
		FROM proyek
		WHERE id = @id
	`
	row := r.DB.QueryRowContext(ctx, query, sql.Named("id", id))

	var p model.Proyek
	var uid mssql.UniqueIdentifier
	var vendor sql.NullString

	err := row.Scan(
		&uid,
		&p.NamaKegiatan,
		&p.Kategori,
		&p.UnitKerja,
		&vendor,
		&p.TanggalGoLive,
		&p.StatusPIR,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	p.ID = uid.String()
	if vendor.Valid {
		p.Vendor = vendor.String
	}

	return &p, nil
}

func (r *ProyekRepository) GetStatistik(ctx context.Context) (*model.StatistikPIR, error) {
	query := `
		SELECT status_pir, COUNT(*) as total
		FROM proyek
		GROUP BY status_pir
	`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stat := &model.StatistikPIR{}
	for rows.Next() {
		var status string
		var total int
		if err := rows.Scan(&status, &total); err != nil {
			return nil, err
		}
		switch status {
		case "selesai":
			stat.Selesai = total
		case "pending":
			stat.Pending = total
		case "terlambat":
			stat.Terlambat = total
		case "belum-jatuh-tempo":
			stat.BelumJatuhTempo = total
		}
		stat.Total += total
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stat, nil
}
