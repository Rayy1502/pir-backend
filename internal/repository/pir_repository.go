package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Rayy1502/pir-backend/internal/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type PIRRepository struct {
	DB *sql.DB
}

func NewPIRRepository(db *sql.DB) *PIRRepository {
	return &PIRRepository{DB: db}
}

func (r *PIRRepository) Create(ctx context.Context, req *model.CreatePIRRequest, createdBy string) (*model.PIRForm, error) {
	status := req.Status
	if status == "" {
		status = "draft"
	}

	// 1. Insert ke pir_form
	query := `
		INSERT INTO pir_form (
			proyek_id, tujuan_awal, metode, catatan_kondisi,
			nama_narasumber, jabatan_narasumber, daftar_pertanyaan,
			kesimpulan_wawancara, responden, status, created_by
		)
		OUTPUT INSERTED.id, INSERTED.created_at, INSERTED.updated_at
		VALUES (
			@proyek_id, @tujuan_awal, @metode, @catatan_kondisi,
			@nama_narasumber, @jabatan_narasumber, @daftar_pertanyaan,
			@kesimpulan_wawancara, @responden, @status, @created_by
		)
	`

	var form model.PIRForm
	var id mssql.UniqueIdentifier

	var createdByVal any = nil
	if createdBy != "" {
		createdByVal = createdBy
	}

	err := r.DB.QueryRowContext(ctx, query,
		sql.Named("proyek_id", req.ProyekID),
		sql.Named("tujuan_awal", req.TujuanAwal),
		sql.Named("metode", req.Metode),
		sql.Named("catatan_kondisi", req.CatatanKondisi),
		sql.Named("nama_narasumber", req.NamaNarasumber),
		sql.Named("jabatan_narasumber", req.JabatanNarasumber),
		sql.Named("daftar_pertanyaan", req.DaftarPertanyaan),
		sql.Named("kesimpulan_wawancara", req.KesimpulanWawancara),
		sql.Named("responden", req.Responden),
		sql.Named("status", status),
		sql.Named("created_by", createdByVal),
	).Scan(&id, &form.CreatedAt, &form.UpdatedAt)

	if err != nil {
		return nil, err
	}

	form.ID = id.String()
	form.ProyekID = req.ProyekID
	form.TujuanAwal = req.TujuanAwal
	form.Metode = req.Metode
	form.CatatanKondisi = req.CatatanKondisi
	form.NamaNarasumber = req.NamaNarasumber
	form.JabatanNarasumber = req.JabatanNarasumber
	form.DaftarPertanyaan = req.DaftarPertanyaan
	form.KesimpulanWawancara = req.KesimpulanWawancara
	form.Responden = req.Responden
	form.Status = status
	if createdBy != "" {
		form.CreatedBy = &createdBy
	}

	// 2. Update Vendor pada Proyek jika diisi
	if req.Vendor != "" {
		_, _ = r.DB.ExecContext(ctx, `
			UPDATE proyek
			SET vendor = @vendor, updated_at = SYSUTCDATETIME()
			WHERE id = @proyek_id
		`, sql.Named("vendor", req.Vendor), sql.Named("proyek_id", req.ProyekID))
	}

	return &form, nil
}

func (r *PIRRepository) GetByID(ctx context.Context, id string) (*model.PIRForm, error) {
	query := `
		SELECT 
			f.id, f.proyek_id, f.tujuan_awal, f.metode, f.catatan_kondisi,
			f.nama_narasumber, f.jabatan_narasumber, f.daftar_pertanyaan,
			f.kesimpulan_wawancara, f.responden, f.status, f.created_by,
			f.created_at, f.updated_at, f.approved_by, f.approved_at,
			p.id, p.nama_kegiatan, p.kategori, p.unit_kerja, p.vendor, p.tanggal_go_live, p.status_pir
		FROM pir_form f
		JOIN proyek p ON f.proyek_id = p.id
		WHERE f.id = @id
	`
	var form model.PIRForm
	var formID, proyekID, createdByUID, pID mssql.UniqueIdentifier
	var crBy, appBy sql.NullString
	var appAt sql.NullTime
	var p model.Proyek

	err := r.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&formID, &proyekID, &form.TujuanAwal, &form.Metode, &form.CatatanKondisi,
		&form.NamaNarasumber, &form.JabatanNarasumber, &form.DaftarPertanyaan,
		&form.KesimpulanWawancara, &form.Responden, &form.Status, &createdByUID,
		&form.CreatedAt, &form.UpdatedAt, &appBy, &appAt,
		&pID, &p.NamaKegiatan, &p.Kategori, &p.UnitKerja, &p.Vendor, &p.TanggalGoLive, &p.StatusPIR,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		// Coba scan jika created_by NULL
		return r.getByIDWithNullCreatedBy(ctx, id)
	}

	form.ID = formID.String()
	form.ProyekID = proyekID.String()
	p.ID = pID.String()
	form.Proyek = &p
	if crBy.Valid {
		form.CreatedBy = &crBy.String
	}
	if appBy.Valid {
		form.ApprovedBy = &appBy.String
	}
	if appAt.Valid {
		form.ApprovedAt = &appAt.Time
	}

	return &form, nil
}

func (r *PIRRepository) getByIDWithNullCreatedBy(ctx context.Context, id string) (*model.PIRForm, error) {
	query := `
		SELECT 
			f.id, f.proyek_id, ISNULL(f.tujuan_awal, ''), ISNULL(f.metode, ''), ISNULL(f.catatan_kondisi, ''),
			ISNULL(f.nama_narasumber, ''), ISNULL(f.jabatan_narasumber, ''), ISNULL(f.daftar_pertanyaan, ''),
			ISNULL(f.kesimpulan_wawancara, ''), ISNULL(f.responden, ''), ISNULL(f.status, 'draft'),
			f.created_at, f.updated_at, f.approved_by, f.approved_at,
			p.id, p.nama_kegiatan, p.kategori, p.unit_kerja, ISNULL(p.vendor, ''), p.tanggal_go_live, p.status_pir
		FROM pir_form f
		JOIN proyek p ON f.proyek_id = p.id
		WHERE f.id = @id
	`
	var form model.PIRForm
	var formID, proyekID, pID mssql.UniqueIdentifier
	var appBy sql.NullString
	var appAt sql.NullTime
	var p model.Proyek

	err := r.DB.QueryRowContext(ctx, query, sql.Named("id", id)).Scan(
		&formID, &proyekID, &form.TujuanAwal, &form.Metode, &form.CatatanKondisi,
		&form.NamaNarasumber, &form.JabatanNarasumber, &form.DaftarPertanyaan,
		&form.KesimpulanWawancara, &form.Responden, &form.Status,
		&form.CreatedAt, &form.UpdatedAt, &appBy, &appAt,
		&pID, &p.NamaKegiatan, &p.Kategori, &p.UnitKerja, &p.Vendor, &p.TanggalGoLive, &p.StatusPIR,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	form.ID = formID.String()
	form.ProyekID = proyekID.String()
	p.ID = pID.String()
	form.Proyek = &p
	if appBy.Valid {
		form.ApprovedBy = &appBy.String
	}
	if appAt.Valid {
		form.ApprovedAt = &appAt.Time
	}

	return &form, nil
}

func (r *PIRRepository) Approve(ctx context.Context, id string, approvedBy string) error {
	query := `
		UPDATE pir_form
		SET status = 'approved',
		    approved_by = @approved_by,
		    approved_at = SYSUTCDATETIME(),
		    updated_at = SYSUTCDATETIME()
		WHERE id = @id
	`
	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("id", id),
		sql.Named("approved_by", approvedBy),
	)
	return err
}

func (r *PIRRepository) GetByProyekID(ctx context.Context, proyekID string) (*model.PIRForm, error) {
	query := `
		SELECT TOP 1 id FROM pir_form WHERE proyek_id = @proyek_id ORDER BY created_at DESC
	`
	var formID mssql.UniqueIdentifier
	err := r.DB.QueryRowContext(ctx, query, sql.Named("proyek_id", proyekID)).Scan(&formID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return r.GetByID(ctx, formID.String())
}

func (r *PIRRepository) Update(ctx context.Context, id string, req *model.UpdatePIRRequest) error {
	status := req.Status
	if status == "" {
		status = "draft"
	}

	query := `
		UPDATE pir_form
		SET tujuan_awal = @tujuan_awal,
		    metode = @metode,
		    catatan_kondisi = @catatan_kondisi,
		    nama_narasumber = @nama_narasumber,
		    jabatan_narasumber = @jabatan_narasumber,
		    daftar_pertanyaan = @daftar_pertanyaan,
		    kesimpulan_wawancara = @kesimpulan_wawancara,
		    responden = @responden,
		    status = @status,
		    updated_at = SYSUTCDATETIME()
		WHERE id = @id
	`
	_, err := r.DB.ExecContext(ctx, query,
		sql.Named("id", id),
		sql.Named("tujuan_awal", req.TujuanAwal),
		sql.Named("metode", req.Metode),
		sql.Named("catatan_kondisi", req.CatatanKondisi),
		sql.Named("nama_narasumber", req.NamaNarasumber),
		sql.Named("jabatan_narasumber", req.JabatanNarasumber),
		sql.Named("daftar_pertanyaan", req.DaftarPertanyaan),
		sql.Named("kesimpulan_wawancara", req.KesimpulanWawancara),
		sql.Named("responden", req.Responden),
		sql.Named("status", status),
	)
	if err != nil {
		return err
	}

	// Update vendor proyek juga jika ada
	if req.Vendor != "" {
		_, _ = r.DB.ExecContext(ctx, `
			UPDATE p
			SET p.vendor = @vendor, p.updated_at = SYSUTCDATETIME()
			FROM proyek p
			JOIN pir_form f ON f.proyek_id = p.id
			WHERE f.id = @id
		`, sql.Named("vendor", req.Vendor), sql.Named("id", id))
	}

	return nil
}
