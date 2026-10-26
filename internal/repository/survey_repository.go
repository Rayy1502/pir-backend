package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Rayy1502/pir-backend/internal/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type SurveyRepository struct {
	DB *sql.DB
}

func NewSurveyRepository(db *sql.DB) *SurveyRepository {
	return &SurveyRepository{DB: db}
}

func (r *SurveyRepository) SaveSurvey(ctx context.Context, pirFormID string, req *model.SubmitSurveyRequest) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Hapus pertanyaan & jawaban lama untuk PIR form ini agar bersih
	_, err = tx.ExecContext(ctx, `
		DELETE FROM survey_pertanyaan WHERE pir_form_id = @pir_form_id
	`, sql.Named("pir_form_id", pirFormID))
	if err != nil {
		return err
	}

	// Helper insert pertanyaan & jawaban
	insertQAndA := func(pertanyaan, tipe string, urutan int, jawabanTeks *string, jawabanAngka *int) error {
		var qID mssql.UniqueIdentifier
		err := tx.QueryRowContext(ctx, `
			INSERT INTO survey_pertanyaan (pir_form_id, pertanyaan, tipe, urutan)
			OUTPUT INSERTED.id
			VALUES (@pir_form_id, @pertanyaan, @tipe, @urutan)
		`,
			sql.Named("pir_form_id", pirFormID),
			sql.Named("pertanyaan", pertanyaan),
			sql.Named("tipe", tipe),
			sql.Named("urutan", urutan),
		).Scan(&qID)
		if err != nil {
			return err
		}

		var jTeksVal any = nil
		if jawabanTeks != nil {
			jTeksVal = *jawabanTeks
		}
		var jAngkaVal any = nil
		if jawabanAngka != nil {
			jAngkaVal = *jawabanAngka
		}

		_, err = tx.ExecContext(ctx, `
			INSERT INTO survey_jawaban (pertanyaan_id, jawaban_teks, jawaban_angka)
			VALUES (@pertanyaan_id, @jawaban_teks, @jawaban_angka)
		`,
			sql.Named("pertanyaan_id", qID),
			sql.Named("jawaban_teks", jTeksVal),
			sql.Named("jawaban_angka", jAngkaVal),
		)
		return err
	}

	// Simpan 4 Pertanyaan Standar
	if err := insertQAndA("Tingkat Kepuasan Pengguna terhadap Sistem", "skala", 1, nil, &req.Kepuasan); err != nil {
		return err
	}
	if err := insertQAndA("Kesesuaian Hasil Implementasi dengan Tujuan Awal", "skala", 2, nil, &req.KesesuaianTujuan); err != nil {
		return err
	}
	if err := insertQAndA("Kendala Operasional / Teknis yang Dihadapi", "teks", 3, &req.Kendala, nil); err != nil {
		return err
	}
	if err := insertQAndA("Saran Perbaikan dan Tindak Lanjut Sistem", "teks", 4, &req.Saran, nil); err != nil {
		return err
	}

	// Simpan Pertanyaan Kustom jika ada
	for i, custom := range req.CustomPertanyaan {
		urutan := 5 + i
		tipe := custom.Tipe
		if tipe == "" {
			tipe = "teks"
		}
		if err := insertQAndA(custom.Pertanyaan, tipe, urutan, custom.JawabanTeks, custom.JawabanAngka); err != nil {
			return err
		}
	}

	// Update pir_form status jika masih draft
	_, _ = tx.ExecContext(ctx, `
		UPDATE pir_form
		SET status = CASE WHEN status = 'draft' THEN 'in_progress' ELSE status END,
		    updated_at = SYSUTCDATETIME()
		WHERE id = @pir_form_id
	`, sql.Named("pir_form_id", pirFormID))

	return tx.Commit()
}

func (r *SurveyRepository) GetByPIRFormID(ctx context.Context, pirFormID string) (*model.SurveyResponse, error) {
	query := `
		SELECT 
			p.id, p.pir_form_id, p.pertanyaan, p.tipe, p.urutan, p.created_at,
			j.id, j.jawaban_teks, j.jawaban_angka, j.created_at
		FROM survey_pertanyaan p
		LEFT JOIN survey_jawaban j ON p.id = j.pertanyaan_id
		WHERE p.pir_form_id = @pir_form_id
		ORDER BY p.urutan ASC
	`
	rows, err := r.DB.QueryContext(ctx, query, sql.Named("pir_form_id", pirFormID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp := &model.SurveyResponse{
		PIRFormID:  pirFormID,
		Pertanyaan: make([]model.SurveyPertanyaan, 0),
	}

	for rows.Next() {
		var p model.SurveyPertanyaan
		var pID, pfID mssql.UniqueIdentifier
		var jID mssql.UniqueIdentifier
		var jTeks sql.NullString
		var jAngka sql.NullInt64
		var jCreated sql.NullTime

		err := rows.Scan(
			&pID, &pfID, &p.Pertanyaan, &p.Tipe, &p.Urutan, &p.CreatedAt,
			&jID, &jTeks, &jAngka, &jCreated,
		)
		if err != nil {
			return nil, err
		}

		p.ID = pID.String()
		p.PIRFormID = pfID.String()

		var jawaban model.SurveyJawaban
		jawaban.ID = jID.String()
		jawaban.PertanyaanID = p.ID
		if jTeks.Valid {
			jawaban.JawabanTeks = &jTeks.String
		}
		if jAngka.Valid {
			val := int(jAngka.Int64)
			jawaban.JawabanAngka = &val
		}
		if jCreated.Valid {
			jawaban.CreatedAt = jCreated.Time
		}
		p.Jawaban = &jawaban

		// Populate standard summary values based on urutan
		switch p.Urutan {
		case 1:
			if jawaban.JawabanAngka != nil {
				resp.Kepuasan = *jawaban.JawabanAngka
			}
		case 2:
			if jawaban.JawabanAngka != nil {
				resp.KesesuaianTujuan = *jawaban.JawabanAngka
			}
		case 3:
			if jawaban.JawabanTeks != nil {
				resp.Kendala = *jawaban.JawabanTeks
			}
		case 4:
			if jawaban.JawabanTeks != nil {
				resp.Saran = *jawaban.JawabanTeks
			}
		}

		resp.Pertanyaan = append(resp.Pertanyaan, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(resp.Pertanyaan) == 0 {
		return nil, errors.New("survey_not_found")
	}

	return resp, nil
}
