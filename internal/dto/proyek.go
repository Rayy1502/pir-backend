package dto

type ProyekFilter struct {
	Query    string
	Kategori string
	Status   string
	Page     int
	PerPage  int
}

type CreateProyekRequest struct {
	NamaKegiatan  string `json:"nama_kegiatan"`
	Kategori      string `json:"kategori"`
	UnitKerja     string `json:"unit_kerja"`
	Vendor        string `json:"vendor"`
	TanggalGoLive string `json:"tanggal_go_live"`
	StatusPIR     string `json:"status_pir,omitempty"`
}

