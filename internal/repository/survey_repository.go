package repository

import (
	"context"
	"database/sql"
	"math"
	"strings"

	"github.com/Rayy1502/pir-backend/internal/model"
	mssql "github.com/microsoft/go-mssqldb"
)

type SurveyRepository struct {
	DB *sql.DB
}

func NewSurveyRepository(db *sql.DB) *SurveyRepository {
	return &SurveyRepository{DB: db}
}

type StandardQuestionItem struct {
	Pertanyaan string
	Tipe       string
	Urutan     int
	Kategori   string
}

func GetStandardQuestions() []StandardQuestionItem {
	return []StandardQuestionItem{
		// 1. Kualitas Sistem & Antarmuka
		{"Saya bisa menemukan menu atau fitur yang saya cari dengan cepat, tanpa coba-coba atau bertanya ke rekan", "skala", 1, "Kualitas Sistem & Antarmuka"},
		{"Proses di aplikasi (misal menyimpan, mencari) berjalan cepat sehingga tidak mengganggu pekerjaan saya", "skala", 2, "Kualitas Sistem & Antarmuka"},
		{"Aplikasi ini jarang error, atau tidak bisa diakses pada jam kerja", "skala", 3, "Keandalan Sistem & Ketersediaan"},
		{"Pada jam kerja, saya dapat tersambung ke aplikasi lewat jaringan kantor setiap kali saya membutuhkannya", "skala", 4, "Akses & Jaringan Kantor"},
		{"Koneksi jaringan jarang terputus saat saya sedang mengisi, menyimpan, atau mengirim data di aplikasi", "skala", 5, "Akses & Jaringan Kantor"},

		// 2. Kualitas Data & Informasi
		{"Data yang saya gunakan dalam pekerjaan (angka, status, atau laporan) akurat dan sesuai dengan kondisi sebenarnya", "skala", 6, "Kualitas Data & Informasi"},
		{"Data yang saya butuhkan selalu terbaru, tanpa keterlambatan pembaruan yang mengganggu pekerjaan", "skala", 7, "Kualitas Data & Informasi"},
		{"Informasinya cukup lengkap sehingga saya tidak perlu mencari ke sumber lain untuk menyelesaikan tugas", "skala", 8, "Kualitas Data & Informasi"},
		{"Informasi di layar bisa langsung saya pahami tanpa perlu penjelasan tambahan", "skala", 9, "Kualitas Data & Informasi"},

		// 3. Layanan & Dukungan TI
		{"Saat saya melaporkan kendala, tim IT/helpdesk menanggapi dalam waktu yang wajar", "skala", 10, "Layanan & Dukungan TI"},
		{"Solusi dari tim IT/helpdesk benar-benar menyelesaikan masalah saya, bukan solusi sementara atau masalah yang berulang", "skala", 11, "Layanan & Dukungan TI"},
		{"Tim IT/helpdesk memahami kebutuhan kerja saya ketika saya meminta bantuan, pelatihan atau perubahan fitur", "skala", 12, "Layanan & Dukungan TI"},

		// 4. Adopsi & Kepuasan Pengguna
		{"Pekerjaan inti saya, saya kerjakan melalui aplikasi ini secara rutin", "skala", 13, "Adopsi & Penggunaan Sistem"},
		{"Saya memakai aplikasi ini karena memang membantu pekerjaan, bukan sekadar karena diwajibkan", "skala", 14, "Adopsi & Penggunaan Sistem"},
		{"Secara keseluruhan, saya puas dengan pengalaman memakai aplikasi ini dalam pekerjaan sehari-hari", "skala", 15, "Kepuasan Pengguna Keseluruhan"},
		{"Saya puas dengan fitur utama yang tersedia", "skala", 16, "Kepuasan Pengguna Keseluruhan"},
		{"Pengalaman saya memakai aplikasi ini sesuai dengan harapan saya sebelum mulai menggunakannya", "skala", 17, "Kesesuaian Ekspektasi Awal"},

		// 5. Dampak Kinerja Pekerjaan
		{"Aplikasi ini mempersingkat waktu penyelesaian tugas saya dibandingkan cara sebelumnya", "skala", 18, "Dampak Kinerja & Efisiensi"},
		{"Aplikasi ini mengurangi pekerjaan berulang atau kesalahan input yang sebelumnya sering terjadi", "skala", 19, "Dampak Kinerja & Efisiensi"},
		{"Informasi dari aplikasi ini membantu saya mengambil keputusan atau tindakan dibandingkan sebelum ada aplikasi", "skala", 20, "Dampak Kinerja & Efisiensi"},

		// 6. Evaluasi Dampak Manfaat (Pilihan Persentase)
		{"Manfaat Business (Peningkatan pendapatan, efisiensi biaya operasional)", "pilihan", 21, "Evaluasi Manfaat Implementasi"},
		{"Manfaat Process (Efisiensi proses bisnis, mempercepat penyelesaian pekerjaan)", "pilihan", 22, "Evaluasi Manfaat Implementasi"},
		{"Manfaat Regulatory (Pemenuhan aturan dan kepatuhan menjadi lebih baik)", "pilihan", 23, "Evaluasi Manfaat Implementasi"},

		// 7. Masukan Kualitatif & Rekomendasi
		{"Apa yang berubah dalam cara Anda bekerja dibandingkan sebelumnya (ditambah, dihilangkan, atau berubah caranya)? Sebutkan perubahan yang paling Anda rasakan (yang lebih baik maupun yang lebih buruk), dan hal dari cara kerja lama yang masih belum tergantikan", "teks", 24, "Evaluasi Kualitatif & Perubahan"},
		{"Tuliskan kebutuhan pengembangan, perbaikan, atau peningkatan modul. Jika tidak ada, tuliskan 'Tidak ada'", "teks", 25, "Rekomendasi Peningkatan Modul"},
	}
}

// EnsureDefaultQuestions memastikan pertanyaan standar lengkap 1-25 ada HANYA untuk metode 'survey' jika belum dibuat
func (r *SurveyRepository) EnsureDefaultQuestions(ctx context.Context, pirFormID string) error {
	// Cek metode PIR Form terlebih dahulu. Hanya metode 'survey' yang memiliki kuesioner pertanyaan!
	var metode string
	err := r.DB.QueryRowContext(ctx, "SELECT ISNULL(metode, 'survey') FROM pir_form WHERE id = @id", sql.Named("id", pirFormID)).Scan(&metode)
	if err != nil {
		return err
	}

	metodeClean := strings.ToLower(strings.TrimSpace(metode))
	if metodeClean != "survey" {
		// Jika bukan metode survey (observasi / wawancara), bersihkan pertanyaan survey yang mungkin terlanjur terbuat tanpa respon
		var ansCount int
		_ = r.DB.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM survey_jawaban sj
			JOIN survey_pertanyaan sp ON sj.pertanyaan_id = sp.id
			WHERE sp.pir_form_id = @id
		`, sql.Named("id", pirFormID)).Scan(&ansCount)
		if ansCount == 0 {
			_, _ = r.DB.ExecContext(ctx, "DELETE FROM survey_pertanyaan WHERE pir_form_id = @id", sql.Named("id", pirFormID))
		}
		return nil
	}

	var count int
	err = r.DB.QueryRowContext(ctx, "SELECT COUNT(*) FROM survey_pertanyaan WHERE pir_form_id = @id", sql.Named("id", pirFormID)).Scan(&count)
	if err != nil {
		return err
	}
	if count >= 20 {
		return nil
	}
	if count > 0 && count < 20 {
		// Cek apakah sudah ada respon
		var ansCount int
		_ = r.DB.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM survey_jawaban sj
			JOIN survey_pertanyaan sp ON sj.pertanyaan_id = sp.id
			WHERE sp.pir_form_id = @id
		`, sql.Named("id", pirFormID)).Scan(&ansCount)
		if ansCount > 0 {
			return nil
		}
		// Hapus pertanyaan placeholder lama yang belum berespon agar digantikan 25 butir lengkap
		_, _ = r.DB.ExecContext(ctx, "DELETE FROM survey_pertanyaan WHERE pir_form_id = @id", sql.Named("id", pirFormID))
	}

	standards := GetStandardQuestions()

	for _, s := range standards {
		_, err := r.DB.ExecContext(ctx, `
			INSERT INTO survey_pertanyaan (pir_form_id, pertanyaan, tipe, urutan, kategori)
			VALUES (@id, @pertanyaan, @tipe, @urutan, @kategori)
		`, sql.Named("id", pirFormID), sql.Named("pertanyaan", s.Pertanyaan), sql.Named("tipe", s.Tipe), sql.Named("urutan", s.Urutan), sql.Named("kategori", s.Kategori))
		if err != nil {
			return err
		}
	}
	return nil
}

// SaveQuestions menyimpan dan menyelaraskan pertanyaan kuesioner yang dirancang PIC (Step 2)
func (r *SurveyRepository) SaveQuestions(ctx context.Context, pirFormID string, questions []model.SurveyItemInput) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Ambil pertanyaan yang sudah ada di DB untuk PIR Form ini
	type ExistingQ struct {
		ID         string
		Urutan     int
		Pertanyaan string
	}
	existingMap := make(map[string]ExistingQ) // key: ID
	existingByUrutan := make(map[int][]ExistingQ)

	rows, err := tx.QueryContext(ctx, `
		SELECT id, urutan, pertanyaan 
		FROM survey_pertanyaan 
		WHERE pir_form_id = @id
		ORDER BY urutan ASC, created_at ASC
	`, sql.Named("id", pirFormID))
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var qID mssql.UniqueIdentifier
		var urutan int
		var teks string
		if err := rows.Scan(&qID, &urutan, &teks); err == nil {
			idStr := qID.String()
			item := ExistingQ{ID: idStr, Urutan: urutan, Pertanyaan: teks}
			existingMap[idStr] = item
			existingByUrutan[urutan] = append(existingByUrutan[urutan], item)
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}

	// Cek metode PIR Form
	var metode string
	_ = tx.QueryRowContext(ctx, "SELECT ISNULL(metode, 'survey') FROM pir_form WHERE id = @id", sql.Named("id", pirFormID)).Scan(&metode)
	metodeClean := strings.ToLower(strings.TrimSpace(metode))

	// 2. Jika input kosong dan metode adalah 'survey', berikan pertanyaan default lengkap
	if len(questions) == 0 && metodeClean == "survey" {
		for _, std := range GetStandardQuestions() {
			questions = append(questions, model.SurveyItemInput{
				Pertanyaan: std.Pertanyaan,
				Tipe:       std.Tipe,
				Urutan:     std.Urutan,
				Kategori:   std.Kategori,
			})
		}
	}

	keptIDs := make(map[string]bool)

	for i, q := range questions {
		urutan := q.Urutan
		if urutan <= 0 {
			urutan = i + 1
		}
		tipe := q.Tipe
		if tipe == "" {
			tipe = "teks"
		}

		// Cari apakah pertanyaan ini sudah ada sebelumnya (berdasarkan ID atau urutan)
		var targetID string
		if q.ID != "" && existingMap[q.ID].ID != "" {
			targetID = q.ID
		} else if list, ok := existingByUrutan[urutan]; ok {
			for _, ex := range list {
				if !keptIDs[ex.ID] {
					targetID = ex.ID
					break
				}
			}
		}

		if targetID != "" {
			// Update pertanyaan yang sudah ada
			keptIDs[targetID] = true
			_, err := tx.ExecContext(ctx, `
				UPDATE survey_pertanyaan 
				SET pertanyaan = @pertanyaan, tipe = @tipe, urutan = @urutan, kategori = @kategori
				WHERE id = @id AND pir_form_id = @pir_id
			`, sql.Named("pertanyaan", q.Pertanyaan), sql.Named("tipe", tipe), sql.Named("urutan", urutan), sql.Named("kategori", q.Kategori), sql.Named("id", targetID), sql.Named("pir_id", pirFormID))
			if err != nil {
				return err
			}
		} else {
			// Insert pertanyaan baru
			var newID mssql.UniqueIdentifier
			err := tx.QueryRowContext(ctx, `
				INSERT INTO survey_pertanyaan (pir_form_id, pertanyaan, tipe, urutan, kategori)
				OUTPUT INSERTED.id
				VALUES (@pir_id, @pertanyaan, @tipe, @urutan, @kategori)
			`, sql.Named("pir_id", pirFormID), sql.Named("pertanyaan", q.Pertanyaan), sql.Named("tipe", tipe), sql.Named("urutan", urutan), sql.Named("kategori", q.Kategori)).Scan(&newID)
			if err != nil {
				return err
			}
			keptIDs[newID.String()] = true
		}
	}

	// 3. Hapus pertanyaan di DB yang tidak dipertahankan (dihapus oleh user di UI atau duplikat lama)
	for id := range existingMap {
		if !keptIDs[id] {
			_, _ = tx.ExecContext(ctx, `
				DELETE FROM survey_pertanyaan WHERE id = @id AND pir_form_id = @pir_id
			`, sql.Named("id", id), sql.Named("pir_id", pirFormID))
		}
	}

	return tx.Commit()
}

// GetPublicSurvey mengambil data proyek dan kuesioner untuk halaman publik responden
func (r *SurveyRepository) GetPublicSurvey(ctx context.Context, pirFormID string) (*model.PublicSurveyResponse, error) {
	_ = r.EnsureDefaultQuestions(ctx, pirFormID)

	query := `
		SELECT 
			pf.id, pf.proyek_id, ISNULL(pf.responden, ''), ISNULL(pf.tujuan_awal, ''),
			p.nama_kegiatan, p.kategori, p.unit_kerja, ISNULL(p.vendor, ''), p.tanggal_go_live
		FROM pir_form pf
		JOIN proyek p ON pf.proyek_id = p.id
		WHERE pf.id = @pir_id
	`
	var resp model.PublicSurveyResponse
	var pfID, pID mssql.UniqueIdentifier
	var tglGoLive sql.NullTime

	err := r.DB.QueryRowContext(ctx, query, sql.Named("pir_id", pirFormID)).Scan(
		&pfID, &pID, &resp.Responden, &resp.TujuanAwal,
		&resp.NamaKegiatan, &resp.Kategori, &resp.UnitKerja, &resp.Vendor, &tglGoLive,
	)
	if err != nil {
		return nil, err
	}

	resp.PIRFormID = pfID.String()
	resp.ProyekID = pID.String()
	if tglGoLive.Valid {
		resp.TanggalGoLive = tglGoLive.Time.Format("2006-01-02")
	}

	// Cek apakah sudah ada respon dan berapa jumlah responden
	var rCount int
	_ = r.DB.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM survey_responden WHERE pir_form_id = @id
	`, sql.Named("id", pirFormID)).Scan(&rCount)
	resp.TotalResponden = rCount
	resp.IsSubmitted = (rCount > 0)

	// Ambil pertanyaan
	rows, err := r.DB.QueryContext(ctx, `
		SELECT id, pertanyaan, tipe, urutan, ISNULL(kategori, '')
		FROM survey_pertanyaan
		WHERE pir_form_id = @id
		ORDER BY urutan ASC
	`, sql.Named("id", pirFormID))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resp.Pertanyaan = make([]model.PublicSurveyQuestion, 0)
	seenPublicUrutan := make(map[int]bool)
	for rows.Next() {
		var q model.PublicSurveyQuestion
		var qID mssql.UniqueIdentifier
		if err := rows.Scan(&qID, &q.Pertanyaan, &q.Tipe, &q.Urutan, &q.Kategori); err == nil {
			if !seenPublicUrutan[q.Urutan] {
				seenPublicUrutan[q.Urutan] = true
				q.ID = qID.String()
				resp.Pertanyaan = append(resp.Pertanyaan, q)
			}
		}
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &resp, nil
}

// PublicSubmit menangani submit kuesioner dari user responden melalui tautan publik
func (r *SurveyRepository) PublicSubmit(ctx context.Context, pirFormID string, req *model.PublicSubmitSurveyRequest) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Simpan responden baru ke tabel survey_responden
	var respID mssql.UniqueIdentifier
	err = tx.QueryRowContext(ctx, `
		INSERT INTO survey_responden (pir_form_id, nama, no_badge, jabatan, unit_kerja, kepuasan, kesesuaian_tujuan, kendala, saran)
		OUTPUT INSERTED.id
		VALUES (@pir_id, @nama, @no_badge, @jabatan, @unit_kerja, @kepuasan, @kesesuaian, @kendala, @saran)
	`, sql.Named("pir_id", pirFormID),
		sql.Named("nama", strings.TrimSpace(req.NamaResponden)),
		sql.Named("no_badge", strings.TrimSpace(req.NoBadge)),
		sql.Named("jabatan", strings.TrimSpace(req.JabatanResponden)),
		sql.Named("unit_kerja", strings.TrimSpace(req.UnitKerjaResponden)),
		sql.Named("kepuasan", req.Kepuasan),
		sql.Named("kesesuaian", req.KesesuaianTujuan),
		sql.Named("kendala", strings.TrimSpace(req.Kendala)),
		sql.Named("saran", strings.TrimSpace(req.Saran)),
	).Scan(&respID)
	if err != nil {
		return err
	}

	// 2. Ambil seluruh pertanyaan untuk PIR form ini
	type QItem struct {
		ID     string
		Urutan int
		Tipe   string
	}
	qMap := make(map[int]QItem)
	qByID := make(map[string]QItem)

	rows, err := tx.QueryContext(ctx, `
		SELECT id, urutan, tipe FROM survey_pertanyaan WHERE pir_form_id = @id ORDER BY urutan ASC
	`, sql.Named("id", pirFormID))
	if err != nil {
		return err
	}
	for rows.Next() {
		var id mssql.UniqueIdentifier
		var urutan int
		var tipe string
		if err := rows.Scan(&id, &urutan, &tipe); err == nil {
			item := QItem{ID: id.String(), Urutan: urutan, Tipe: tipe}
			qMap[urutan] = item
			qByID[id.String()] = item
		}
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	// Helper simpan jawaban per responden
	insertAnswer := func(qID string, teks *string, angka *int) error {
		var jTeks any = nil
		if teks != nil {
			jTeks = *teks
		}
		var jAngka any = nil
		if angka != nil {
			jAngka = *angka
		}
		_, err := tx.ExecContext(ctx, `
			INSERT INTO survey_jawaban (pertanyaan_id, jawaban_teks, jawaban_angka, responden_id)
			VALUES (@q_id, @teks, @angka, @resp_id)
		`, sql.Named("q_id", qID), sql.Named("teks", jTeks), sql.Named("angka", jAngka), sql.Named("resp_id", respID))
		return err
	}

	savedQIDs := make(map[string]bool)

	// Simpan seluruh butir jawaban dari kuesioner
	for _, jl := range req.JawabanLain {
		if jl.PertanyaanID != "" {
			_ = insertAnswer(jl.PertanyaanID, jl.JawabanTeks, jl.JawabanAngka)
			savedQIDs[jl.PertanyaanID] = true
		}
	}

	// Fallback untuk backward compatibility jika ada pertanyaan 1-4 yang belum tersimpan
	if q, ok := qMap[1]; ok && !savedQIDs[q.ID] {
		_ = insertAnswer(q.ID, nil, &req.Kepuasan)
	}
	if q, ok := qMap[2]; ok && !savedQIDs[q.ID] {
		_ = insertAnswer(q.ID, nil, &req.KesesuaianTujuan)
	}
	if q, ok := qMap[3]; ok && !savedQIDs[q.ID] && req.Kendala != "" {
		_ = insertAnswer(q.ID, &req.Kendala, nil)
	}
	if q, ok := qMap[4]; ok && !savedQIDs[q.ID] && req.Saran != "" {
		_ = insertAnswer(q.ID, &req.Saran, nil)
	}

	// 3. Update pir_form status jika draft
	_, err = tx.ExecContext(ctx, `
		UPDATE pir_form
		SET status = CASE WHEN status = 'draft' THEN 'in_progress' ELSE status END,
		    updated_at = SYSUTCDATETIME()
		WHERE id = @id
	`, sql.Named("id", pirFormID))
	if err != nil {
		return err
	}

	return tx.Commit()
}

// GetByPIRFormID mengambil data survey lengkap untuk PIC di Step 2
func (r *SurveyRepository) GetByPIRFormID(ctx context.Context, pirFormID string) (*model.SurveyResponse, error) {
	_ = r.EnsureDefaultQuestions(ctx, pirFormID)

	queryProyek := `
		SELECT 
			pf.id, pf.proyek_id, ISNULL(pf.responden, ''), ISNULL(pf.metode, 'survey'),
			p.nama_kegiatan, p.unit_kerja, ISNULL(p.vendor, '')
		FROM pir_form pf
		JOIN proyek p ON pf.proyek_id = p.id
		WHERE pf.id = @pir_id
	`
	var resp model.SurveyResponse
	var pfID, pID mssql.UniqueIdentifier
	var metode string
	err := r.DB.QueryRowContext(ctx, queryProyek, sql.Named("pir_id", pirFormID)).Scan(
		&pfID, &pID, &resp.Responden, &metode,
		&resp.NamaKegiatan, &resp.UnitKerja, &resp.Vendor,
	)
	if err != nil {
		return nil, err
	}
	resp.PIRFormID = pfID.String()
	resp.ProyekID = pID.String()
	resp.Pertanyaan = make([]model.SurveyPertanyaan, 0)
	resp.DaftarResponden = make([]model.RespondenDetail, 0)

	// 1. Ambil daftar responden dari tabel survey_responden
	rowsResp, err := r.DB.QueryContext(ctx, `
		SELECT id, pir_form_id, nama, ISNULL(no_badge, ''), ISNULL(jabatan, ''), ISNULL(unit_kerja, ''), kepuasan, kesesuaian_tujuan, ISNULL(kendala, ''), ISNULL(saran, ''), created_at
		FROM survey_responden
		WHERE pir_form_id = @pir_id
		ORDER BY created_at DESC
	`, sql.Named("pir_id", pirFormID))
	if err == nil {
		defer rowsResp.Close()
		for rowsResp.Next() {
			var rd model.RespondenDetail
			var rdID, formID mssql.UniqueIdentifier
			if err := rowsResp.Scan(&rdID, &formID, &rd.Nama, &rd.NoBadge, &rd.Jabatan, &rd.UnitKerja, &rd.Kepuasan, &rd.KesesuaianTujuan, &rd.Kendala, &rd.Saran, &rd.CreatedAt); err == nil {
				rd.ID = rdID.String()
				rd.PIRFormID = formID.String()
				resp.DaftarResponden = append(resp.DaftarResponden, rd)
			}
		}
		_ = rowsResp.Err()
	}

	resp.TotalResponden = len(resp.DaftarResponden)
	resp.IsSubmitted = (resp.TotalResponden > 0)

	// Hitung rata-rata dan gabungan kendala/saran TANPA nama responden (dirangkum kolektif)
	if resp.TotalResponden > 0 && metode != "observasi" && metode != "wawancara" {
		var sumKepuasan, sumKesesuaian int
		var listKendala, listSaran []string
		seenKendala := make(map[string]bool)
		seenSaran := make(map[string]bool)

		cleanText := func(s string) string {
			s = strings.TrimSpace(s)
			// Hapus prefix nama jika tersimpan pola "Nama: Isi" sebelumnya
			if idx := strings.Index(s, ": "); idx != -1 && idx < 35 {
				s = strings.TrimSpace(s[idx+2:])
			}
			// Hapus karakter bullet atau nomor di awal
			s = strings.TrimLeft(s, "-*•0123456789. ")
			return strings.TrimSpace(s)
		}

		for _, rd := range resp.DaftarResponden {
			sumKepuasan += rd.Kepuasan
			sumKesesuaian += rd.KesesuaianTujuan

			if k := cleanText(rd.Kendala); k != "" {
				key := strings.ToLower(k)
				if !seenKendala[key] {
					seenKendala[key] = true
					listKendala = append(listKendala, k)
				}
			}
			if s := cleanText(rd.Saran); s != "" {
				key := strings.ToLower(s)
				if !seenSaran[key] {
					seenSaran[key] = true
					listSaran = append(listSaran, s)
				}
			}
		}

		avgKep := float64(sumKepuasan) / float64(resp.TotalResponden)
		resp.KepuasanAvg = math.Round(avgKep*10) / 10
		resp.Kepuasan = int(math.Round(avgKep))
		resp.KesesuaianTujuan = int(math.Round(float64(sumKesesuaian) / float64(resp.TotalResponden)))

		if len(listKendala) > 0 {
			resp.Kendala = strings.Join(listKendala, "\n")
		} else {
			resp.Kendala = "Tidak ada kendala operasional yang dilaporkan oleh responden."
		}

		if len(listSaran) > 0 {
			resp.Saran = strings.Join(listSaran, "\n")
		} else {
			resp.Saran = "Seluruh responden menyatakan sistem telah beroperasi secara optimal."
		}
	}

	// 2. Pertanyaan kuesioner HANYA diambil dan digunakan jika metode adalah 'survey'
	metodeClean := strings.ToLower(strings.TrimSpace(metode))
	if metodeClean != "survey" {
		resp.Pertanyaan = make([]model.SurveyPertanyaan, 0)
		return &resp, nil
	}

	query := `
		SELECT 
			p.id, p.pir_form_id, p.pertanyaan, p.tipe, p.urutan, ISNULL(p.kategori, ''), p.created_at,
			CAST(j.id AS NVARCHAR(36)) AS j_id, j.jawaban_teks, j.jawaban_angka, j.created_at
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

	seenUrutan := make(map[int]bool)
	resp.Pertanyaan = make([]model.SurveyPertanyaan, 0)
	for rows.Next() {
		var p model.SurveyPertanyaan
		var pID, pfID mssql.UniqueIdentifier
		var jID sql.NullString
		var jTeks sql.NullString
		var jAngka sql.NullInt64
		var jCreated sql.NullTime

		err := rows.Scan(
			&pID, &pfID, &p.Pertanyaan, &p.Tipe, &p.Urutan, &p.Kategori, &p.CreatedAt,
			&jID, &jTeks, &jAngka, &jCreated,
		)
		if err != nil {
			return nil, err
		}

		p.ID = pID.String()
		p.PIRFormID = pfID.String()

		if jID.Valid && jID.String != "" && jID.String != "00000000-0000-0000-0000-000000000000" && (jTeks.Valid || jAngka.Valid) {
			var jawaban model.SurveyJawaban
			jawaban.ID = jID.String
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

			// Evaluasi teknis observasi & wawancara atau fallback jika survey belum ada responden
			if resp.TotalResponden == 0 || metode == "observasi" || metode == "wawancara" {
				resp.IsSubmitted = true
				switch p.Urutan {
				case 1:
					if jawaban.JawabanAngka != nil {
						resp.Kepuasan = *jawaban.JawabanAngka
						resp.KepuasanAvg = float64(*jawaban.JawabanAngka)
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
			}
		}

		if !seenUrutan[p.Urutan] {
			seenUrutan[p.Urutan] = true
			resp.Pertanyaan = append(resp.Pertanyaan, p)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &resp, nil
}

// SaveSurvey menyimpan jawaban secara manual oleh PIC (Step 2 fallback jika wawancara/observasi)
func (r *SurveyRepository) SaveSurvey(ctx context.Context, pirFormID string, req *model.SubmitSurveyRequest) error {
	tx, err := r.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 1. Hapus pertanyaan & jawaban lama
	_, err = tx.ExecContext(ctx, `
		DELETE FROM survey_pertanyaan WHERE pir_form_id = @pir_form_id
	`, sql.Named("pir_form_id", pirFormID))
	if err != nil {
		return err
	}

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

	// 4 Pertanyaan Standar
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

	// Update pir_form responden jika ada & status
	if req.Responden != "" {
		_, _ = tx.ExecContext(ctx, `
			UPDATE pir_form
			SET responden = @responden,
			    status = CASE WHEN status = 'draft' THEN 'in_progress' ELSE status END,
			    updated_at = SYSUTCDATETIME()
			WHERE id = @pir_form_id
		`, sql.Named("responden", req.Responden), sql.Named("pir_form_id", pirFormID))
	} else {
		_, _ = tx.ExecContext(ctx, `
			UPDATE pir_form
			SET status = CASE WHEN status = 'draft' THEN 'in_progress' ELSE status END,
			    updated_at = SYSUTCDATETIME()
			WHERE id = @pir_form_id
		`, sql.Named("pir_form_id", pirFormID))
	}

	return tx.Commit()
}
