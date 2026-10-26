package model

import "time"

type LaporanPIR struct {
	ID               string          `json:"id"`
	NomorPIR         string          `json:"nomor_pir"`
	ProyekID         string          `json:"proyek_id"`
	PIRFormID        string          `json:"pir_form_id"`
	Ringkasan        string          `json:"ringkasan"`
	Kepuasan         float64         `json:"kepuasan"`
	KesesuaianTujuan int             `json:"kesesuaian_tujuan"`
	Rekomendasi      string          `json:"rekomendasi"`
	TanggalLaporan   string          `json:"tanggal_laporan"`
	PembuatID        *string         `json:"pembuat_id,omitempty"`
	PembuatNama      string          `json:"pembuat_nama,omitempty"`
	FileTTDURL       *string         `json:"file_ttd_url,omitempty"`
	FileTTDNama      *string         `json:"file_ttd_nama,omitempty"`
	StatusLaporan    string          `json:"status_laporan"` // 'draft', 'signed', 'approved'
	SignedAt         *time.Time      `json:"signed_at,omitempty"`
	CreatedAt        time.Time       `json:"created_at"`
	Proyek           *Proyek         `json:"proyek,omitempty"`
	PIRForm          *PIRForm        `json:"pir_form,omitempty"`
	Temuan           []LaporanTemuan `json:"temuan,omitempty"`
}

type LaporanTemuan struct {
	ID        string `json:"id"`
	LaporanID string `json:"laporan_id"`
	Temuan    string `json:"temuan"`
	Urutan    int    `json:"urutan"`
}

type CreateLaporanRequest struct {
	ProyekID         string   `json:"proyek_id"`
	PIRFormID        string   `json:"pir_form_id"`
	Ringkasan        string   `json:"ringkasan"`
	Kepuasan         float64  `json:"kepuasan"`
	KesesuaianTujuan int      `json:"kesesuaian_tujuan"`
	Rekomendasi      string   `json:"rekomendasi"`
	TanggalLaporan   string   `json:"tanggal_laporan"`
	Temuan           []string `json:"temuan"`
}
