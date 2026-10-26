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

type SubmitSurveyRequest struct {
	Kepuasan          int               `json:"kepuasan"`           // Skala 1-5
	KesesuaianTujuan  int               `json:"kesesuaian_tujuan"`  // Skala 1-5
	Kendala           string            `json:"kendala"`            // Teks kendala
	Saran             string            `json:"saran"`              // Teks saran
	CustomPertanyaan  []SurveyItemInput `json:"custom_pertanyaan"`  // Pertanyaan dinamis tambahan
}

type SurveyResponse struct {
	PIRFormID        string             `json:"pir_form_id"`
	Kepuasan         int                `json:"kepuasan"`
	KesesuaianTujuan int                `json:"kesesuaian_tujuan"`
	Kendala          string             `json:"kendala"`
	Saran            string             `json:"saran"`
	Pertanyaan       []SurveyPertanyaan `json:"pertanyaan"`
}
