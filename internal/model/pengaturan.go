package model

import "time"

type PengaturanUnitKerja struct {
	ID              string    `json:"id"`
	UnitKerja       string    `json:"unit_kerja"`
	BatasWaktuBulan int       `json:"batas_waktu_bulan"`
	Aktif           bool      `json:"aktif"`
	UpdatedAt       time.Time `json:"updated_at"`
}
