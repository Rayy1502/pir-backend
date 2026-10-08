package repository

import (
	"context"
	"database/sql"

	"github.com/Rayy1502/pir-backend/internal/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type PengaturanRepository struct {
	DB *sql.DB
}

func NewPengaturanRepository(db *sql.DB) *PengaturanRepository {
	return &PengaturanRepository{DB: db}
}

func (r *PengaturanRepository) GetAll(ctx context.Context) ([]model.PengaturanUnitKerja, error) {
	query := `
		SELECT id, unit_kerja, batas_waktu_bulan, aktif, updated_at
		FROM pengaturan_unit_kerja
		WHERE aktif = 1
		ORDER BY unit_kerja
	`
	rows, err := r.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []model.PengaturanUnitKerja
	for rows.Next() {
		var p model.PengaturanUnitKerja
		var uid mssql.UniqueIdentifier
		if err := rows.Scan(&uid, &p.UnitKerja, &p.BatasWaktuBulan, &p.Aktif, &p.UpdatedAt); err != nil {
			return nil, err
		}
		p.ID = uid.String()
		list = append(list, p)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *PengaturanRepository) GetByUnitKerja(ctx context.Context, unit string) (*model.PengaturanUnitKerja, error) {
	query := `
		SELECT id, unit_kerja, batas_waktu_bulan, aktif, updated_at
		FROM pengaturan_unit_kerja
		WHERE unit_kerja = @unit AND aktif = 1
	`
	var p model.PengaturanUnitKerja
	var uid mssql.UniqueIdentifier
	err := r.DB.QueryRowContext(ctx, query, sql.Named("unit", unit)).Scan(
		&uid, &p.UnitKerja, &p.BatasWaktuBulan, &p.Aktif, &p.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	p.ID = uid.String()
	return &p, nil
}

func (r *PengaturanRepository) Update(ctx context.Context, unit string, batasBulan int) error {
	query := `
		UPDATE pengaturan_unit_kerja
		SET batas_waktu_bulan = @batas, updated_at = SYSUTCDATETIME()
		WHERE unit_kerja = @unit
	`
	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("batas", batasBulan),
		sql.Named("unit", unit),
	)
	return err
}
