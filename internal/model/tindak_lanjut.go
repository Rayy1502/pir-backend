package model

import "time"

type TindakLanjut struct {
	ID            string              `json:"id"`
	LaporanID     string              `json:"laporan_id"`
	NomorPIR      string              `json:"nomor_pir,omitempty"`
	ProyekID      string              `json:"proyek_id,omitempty"`
	NamaKegiatan  string              `json:"nama_kegiatan,omitempty"`
	UnitKerja     string              `json:"unit_kerja,omitempty"`
	Temuan        string              `json:"temuan"`
	PIC           string              `json:"pic"`
	TargetSelesai string              `json:"target_selesai"`
	Status        string              `json:"status"` // 'open', 'in_progress', 'done'
	CreatedAt     time.Time           `json:"created_at"`
	UpdatedAt     time.Time           `json:"updated_at"`
	Bukti         []TindakLanjutBukti `json:"bukti,omitempty"`
}

type TindakLanjutBukti struct {
	ID             string    `json:"id"`
	TindakLanjutID string    `json:"tindak_lanjut_id"`
	NamaFile       string    `json:"nama_file"`
	URL            string    `json:"url"`
	UkuranBytes    int64     `json:"ukuran_bytes"`
	UploadedAt     time.Time `json:"uploaded_at"`
}

type CreateTindakLanjutRequest struct {
	LaporanID     string `json:"laporan_id"`
	Temuan        string `json:"temuan"`
	PIC           string `json:"pic"`
	TargetSelesai string `json:"target_selesai"`
	Status        string `json:"status"` // default 'open'
}

type UpdateTindakLanjutRequest struct {
	Temuan        string `json:"temuan"`
	PIC           string `json:"pic"`
	TargetSelesai string `json:"target_selesai"`
	Status        string `json:"status"`
}
