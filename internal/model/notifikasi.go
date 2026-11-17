package model

import "time"

type Notifikasi struct {
	ID           string    `json:"id"`
	ProyekID     string    `json:"proyek_id"`
	Tipe         string    `json:"tipe"` // "waktu_pir", "segera", "terlambat"
	Judul        string    `json:"judul"`
	Deskripsi    string    `json:"deskripsi"`
	Dibaca       bool      `json:"dibaca"`
	CreatedAt    time.Time `json:"created_at"`
	NamaKegiatan string    `json:"nama_kegiatan,omitempty"`
	UnitKerja    string    `json:"unit_kerja,omitempty"`
	TanggalRilis string    `json:"tanggal_rilis,omitempty"`
}

type NotifikasiFilter struct {
	Dibaca *bool
	Limit  int
	Offset int
}
