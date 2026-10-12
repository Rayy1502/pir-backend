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