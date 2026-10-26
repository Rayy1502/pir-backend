package model

import "time"

type PIRForm struct {
	ID                  string     `json:"id"`
	ProyekID            string     `json:"proyek_id"`
	TujuanAwal          string     `json:"tujuan_awal"`
	Metode              string     `json:"metode"` // 'observasi', 'wawancara', 'survey'
	CatatanKondisi      string     `json:"catatan_kondisi"`
	NamaNarasumber      string     `json:"nama_narasumber"`
	JabatanNarasumber   string     `json:"jabatan_narasumber"`
	DaftarPertanyaan    string     `json:"daftar_pertanyaan"`
	KesimpulanWawancara string     `json:"kesimpulan_wawancara"`
	Responden           string     `json:"responden"`
	Status              string     `json:"status"` // 'draft', 'submitted', 'approved', 'in_progress', 'selesai'
	ApprovedBy          *string    `json:"approved_by,omitempty"`
	ApprovedAt          *time.Time `json:"approved_at,omitempty"`
	CreatedBy           *string    `json:"created_by,omitempty"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
	
	// Joined Proyek detail
	Proyek              *Proyek    `json:"proyek,omitempty"`
}

type ApprovePIRRequest struct {
	ApprovedBy string `json:"approved_by"`
}

type CreatePIRRequest struct {
	ProyekID            string `json:"proyek_id"`
	Vendor              string `json:"vendor"`
	TujuanAwal          string `json:"tujuan_awal"`
	Metode              string `json:"metode"`
	CatatanKondisi      string `json:"catatan_kondisi"`
	NamaNarasumber      string `json:"nama_narasumber"`
	JabatanNarasumber   string `json:"jabatan_narasumber"`
	DaftarPertanyaan    string `json:"daftar_pertanyaan"`
	KesimpulanWawancara string `json:"kesimpulan_wawancara"`
	Responden           string `json:"responden"`
	Status              string `json:"status"` // 'draft' or 'submitted'
}

type UpdatePIRRequest struct {
	Vendor              string `json:"vendor"`
	TujuanAwal          string `json:"tujuan_awal"`
	Metode              string `json:"metode"`
	CatatanKondisi      string `json:"catatan_kondisi"`
	NamaNarasumber      string `json:"nama_narasumber"`
	JabatanNarasumber   string `json:"jabatan_narasumber"`
	DaftarPertanyaan    string `json:"daftar_pertanyaan"`
	KesimpulanWawancara string `json:"kesimpulan_wawancara"`
	Responden           string `json:"responden"`
	Status              string `json:"status"`
}
