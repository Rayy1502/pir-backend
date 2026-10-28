package model

import "time"

type SurveyPertanyaan struct {
	ID          string         `json:"id"`
	PIRFormID   string         `json:"pir_form_id"`
	Pertanyaan  string         `json:"pertanyaan"`
	Tipe        string         `json:"tipe"` // 'skala', 'teks', 'pilihan'
	Urutan      int            `json:"urutan"`
	CreatedAt   time.Time      `json:"created_at"`
	Jawaban     *SurveyJawaban `json:"jawaban,omitempty"`
}

type SurveyJawaban struct {
	ID           string    `json:"id"`
	PertanyaanID string    `json:"pertanyaan_id"`
	JawabanTeks  *string   `json:"jawaban_teks,omitempty"`
	JawabanAngka *int      `json:"jawaban_angka,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
}

type SurveyItemInput struct {
	ID           string  `json:"id,omitempty"`
	Pertanyaan   string  `json:"pertanyaan"`
	Tipe         string  `json:"tipe"` // 'skala', 'teks', 'pilihan'
	Urutan       int     `json:"urutan"`
	JawabanTeks  *string `json:"jawaban_teks,omitempty"`
	JawabanAngka *int    `json:"jawaban_angka,omitempty"`
}

type SaveQuestionsRequest struct {
	Pertanyaan []SurveyItemInput `json:"pertanyaan"`
}

type PublicSurveyQuestion struct {
	ID         string `json:"id"`
	Pertanyaan string `json:"pertanyaan"`
	Tipe       string `json:"tipe"` // 'skala', 'teks', 'pilihan'
	Urutan     int    `json:"urutan"`
}

type PublicSurveyResponse struct {
	PIRFormID      string                 `json:"pir_form_id"`
	ProyekID       string                 `json:"proyek_id"`
	NamaKegiatan   string                 `json:"nama_kegiatan"`
	Kategori       string                 `json:"kategori"`
	UnitKerja      string                 `json:"unit_kerja"`
	Vendor         string                 `json:"vendor"`
	TanggalGoLive  string                 `json:"tanggal_go_live"`
	TujuanAwal     string                 `json:"tujuan_awal"`
	Responden      string                 `json:"responden,omitempty"`
	IsSubmitted    bool                   `json:"is_submitted"`
	TotalResponden int                    `json:"total_responden"`
	Pertanyaan     []PublicSurveyQuestion `json:"pertanyaan"`
}

type SurveyJawabanItem struct {
	PertanyaanID string  `json:"pertanyaan_id"`
	JawabanTeks  *string `json:"jawaban_teks,omitempty"`
	JawabanAngka *int    `json:"jawaban_angka,omitempty"`
}

type PublicSubmitSurveyRequest struct {
	NamaResponden      string              `json:"nama_responden"`
	JabatanResponden   string              `json:"jabatan_responden,omitempty"`
	UnitKerjaResponden string              `json:"unit_kerja_responden,omitempty"`
	Kepuasan           int                 `json:"kepuasan"`          // Skala 1-5
	KesesuaianTujuan   int                 `json:"kesesuaian_tujuan"` // Skala 1-5
	Kendala            string              `json:"kendala"`
	Saran              string              `json:"saran"`
	JawabanLain        []SurveyJawabanItem `json:"jawaban_lain,omitempty"`
}

type SubmitSurveyRequest struct {
	Responden         string            `json:"responden,omitempty"`
	Kepuasan          int               `json:"kepuasan"`           // Skala 1-5
	KesesuaianTujuan  int               `json:"kesesuaian_tujuan"`  // Skala 1-5
	Kendala           string            `json:"kendala"`            // Teks kendala
	Saran             string            `json:"saran"`              // Teks saran
	CustomPertanyaan  []SurveyItemInput `json:"custom_pertanyaan"`  // Pertanyaan dinamis tambahan
}

type RespondenDetail struct {
	ID               string    `json:"id"`
	PIRFormID        string    `json:"pir_form_id"`
	Nama             string    `json:"nama"`
	Jabatan          string    `json:"jabatan"`
	UnitKerja        string    `json:"unit_kerja"`
	Kepuasan         int       `json:"kepuasan"`
	KesesuaianTujuan int       `json:"kesesuaian_tujuan"`
	Kendala          string    `json:"kendala"`
	Saran            string    `json:"saran"`
	CreatedAt        time.Time `json:"created_at"`
}

type SurveyResponse struct {
	PIRFormID        string             `json:"pir_form_id"`
	ProyekID         string             `json:"proyek_id"`
	NamaKegiatan     string             `json:"nama_kegiatan"`
	UnitKerja        string             `json:"unit_kerja"`
	Vendor           string             `json:"vendor"`
	Responden        string             `json:"responden"`
	IsSubmitted      bool               `json:"is_submitted"`
	TotalResponden   int                `json:"total_responden"`
	Kepuasan         int                `json:"kepuasan"`
	KepuasanAvg      float64            `json:"kepuasan_avg"`
	KesesuaianTujuan int                `json:"kesesuaian_tujuan"`
	Kendala          string             `json:"kendala"`
	Saran            string             `json:"saran"`
	DaftarResponden  []RespondenDetail  `json:"daftar_responden,omitempty"`
	Pertanyaan       []SurveyPertanyaan `json:"pertanyaan"`
}

