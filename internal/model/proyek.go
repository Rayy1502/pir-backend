package model

import "time"

type Proyek struct {
	ID                string     `json:"id"`
	NamaKegiatan      string     `json:"nama_kegiatan"`
	Kategori          string     `json:"kategori"`
	UnitKerja         string     `json:"unit_kerja"`
	Vendor            string     `json:"vendor"`
	TanggalGoLive     time.Time  `json:"tanggal_go_live"`
	StatusPIR         string     `json:"status_pir"`
	Kepuasan          *float64   `json:"kepuasan,omitempty"`
	KesesuaianTujuan  *int       `json:"kesesuaian_tujuan,omitempty"`
	TindakLanjutTotal *int       `json:"tindak_lanjut_total,omitempty"`
	TindakLanjutDone  *int       `json:"tindak_lanjut_done,omitempty"`
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}

type TimelineStage struct {
	Key       string `json:"key"`
	Nama      string `json:"nama"`
	Deskripsi string `json:"deskripsi"`
	Status    string `json:"status"` // 'selesai', 'pending', 'terlambat', 'belum_dimulai'
	Tanggal   string `json:"tanggal,omitempty"`
}

type ProyekDetail struct {
	Proyek
	BatasWaktuBulan   int             `json:"batas_waktu_bulan"`
	TanggalJatuhTempo string          `json:"tanggal_jatuh_tempo"`
	Timeline          []TimelineStage `json:"timeline"`
	PIRFormID         *string         `json:"pir_form_id,omitempty"`
	PIRFormStatus     *string         `json:"pir_form_status,omitempty"`
	ApprovedBy        *string         `json:"approved_by,omitempty"`
	ApprovedAt        *string         `json:"approved_at,omitempty"`
	LaporanID         *string         `json:"laporan_id,omitempty"`
	StatusLaporan     *string         `json:"status_laporan,omitempty"`
	FileTTDURL        *string         `json:"file_ttd_url,omitempty"`
}

type StatistikPIR struct {
	Selesai         int `json:"selesai"`
	Pending         int `json:"pending"`
	Terlambat       int `json:"terlambat"`
	BelumJatuhTempo int `json:"belum_jatuh_tempo"`
	Total           int `json:"total"`
}

type DashboardStats struct {
	TotalProyek   int `json:"total_proyek"`
	ProyekSelesai int `json:"proyek_selesai"`
	PIRPending    int `json:"pir_pending"`
	PIRSelesai    int `json:"pir_selesai"`
}