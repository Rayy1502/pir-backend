package pdf

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/go-pdf/fpdf"
)

func GenerateLaporanPDF(l *model.LaporanPIR) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 15, 15)
	pdf.SetAutoPageBreak(true, 15)
	pdf.AddPage()

	// 1. KOP SURAT
	pdf.SetFont("Arial", "B", 13)
	pdf.SetTextColor(0, 51, 102) // PUSRI Navy
	pdf.CellFormat(0, 6, "PT PUPUK SRIWIDJAJA PALEMBANG", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "B", 10)
	pdf.SetTextColor(0, 102, 51) // PUSRI Green
	pdf.CellFormat(0, 5, "DIVISI TEKNOLOGI INFORMASI", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	pdf.SetTextColor(100, 100, 100)
	pdf.CellFormat(0, 4, "Jl. Mayor Zen, Palembang 30118, Sumatera Selatan | Telp: (0711) 712222 | ti@pusri.co.id", "", 1, "C", false, 0, "")

	// Double Line Divider
	y := pdf.GetY() + 2
	pdf.SetDrawColor(0, 51, 102)
	pdf.SetLineWidth(0.8)
	pdf.Line(15, y, 195, y)
	pdf.SetLineWidth(0.3)
	pdf.Line(15, y+1.2, 195, y+1.2)
	pdf.SetY(y + 6)

	// 2. JUDUL DOKUMEN & NOMOR
	pdf.SetFont("Arial", "B", 12)
	pdf.SetTextColor(20, 20, 20)
	pdf.CellFormat(0, 6, "LAPORAN POST IMPLEMENTATION REVIEW (PIR)", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 9)
	pdf.SetTextColor(70, 70, 70)
	pdf.CellFormat(0, 5, fmt.Sprintf("Nomor: %s", l.NomorPIR), "", 1, "C", false, 0, "")
	pdf.Ln(3)

	// 3. TABEL METADATA KEGIATAN
	pdf.SetFillColor(240, 244, 248)
	pdf.SetTextColor(0, 51, 102)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 6, "  A. INFORMASI PROYEK PENGEMBANGAN TI", "1", 1, "L", true, 0, "")

	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Arial", "", 8.5)

	namaKegiatan := "-"
	kategori := "-"
	unitKerja := "-"
	vendor := "-"
	goLive := "-"
	metode := "-"
	tujuan := "-"

	if l.Proyek != nil {
		namaKegiatan = l.Proyek.NamaKegiatan
		kategori = strings.ToUpper(l.Proyek.Kategori)
		unitKerja = l.Proyek.UnitKerja
		if l.Proyek.Vendor != "" {
			vendor = l.Proyek.Vendor
		}
		goLive = l.Proyek.TanggalGoLive.Format("02 January 2006")
	}
	if l.PIRForm != nil {
		if l.PIRForm.Metode != "" {
			metode = strings.ToUpper(l.PIRForm.Metode)
		}
		if l.PIRForm.TujuanAwal != "" {
			tujuan = l.PIRForm.TujuanAwal
		}
	}

	renderMetaRow := func(label, val string) {
		pdf.SetFont("Arial", "B", 8)
		pdf.CellFormat(45, 5.5, "  "+label, "L", 0, "L", false, 0, "")
		pdf.CellFormat(5, 5.5, ":", "", 0, "C", false, 0, "")
		pdf.SetFont("Arial", "", 8)
		pdf.CellFormat(130, 5.5, val, "R", 1, "L", false, 0, "")
	}

	renderMetaRow("Nama Kegiatan", namaKegiatan)
	renderMetaRow("Kategori Pengembangan", kategori)
	renderMetaRow("Unit Kerja Pemilik", unitKerja)
	renderMetaRow("Vendor / Pengembang", vendor)
	renderMetaRow("Tanggal Go Live", goLive)
	renderMetaRow("Metode PIR", metode)
	renderMetaRow("Tanggal Laporan", l.TanggalLaporan)
	pdf.CellFormat(0, 1, "", "T", 1, "L", false, 0, "")
	pdf.Ln(3)

	// 4. TUJUAN AWAL & RINGKASAN HASIL
	if tujuan != "-" && tujuan != "" {
		pdf.SetFillColor(240, 244, 248)
		pdf.SetTextColor(0, 51, 102)
		pdf.SetFont("Arial", "B", 9)
		pdf.CellFormat(0, 6, "  B. TUJUAN AWAL PENGEMBANGAN", "1", 1, "L", true, 0, "")
		pdf.SetTextColor(40, 40, 40)
		pdf.SetFont("Arial", "", 8.5)
		pdf.MultiCell(0, 5, tujuan, "LRB", "L", false)
		pdf.Ln(3)
	}

	pdf.SetFillColor(240, 244, 248)
	pdf.SetTextColor(0, 51, 102)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 6, "  C. RINGKASAN HASIL EVALUASI", "1", 1, "L", true, 0, "")
	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Arial", "", 8.5)
	ringkasan := l.Ringkasan
	if ringkasan == "" {
		ringkasan = "Implementasi proyek telah melalui evaluasi menyeluruh bersama stakeholder terkait."
	}
	pdf.MultiCell(0, 5, ringkasan, "LRB", "L", false)
	pdf.Ln(3)

	// 5. INDIKATOR SKOR KPI
	pdf.SetFillColor(240, 244, 248)
	pdf.SetTextColor(0, 51, 102)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 6, "  D. INDIKATOR KEBERHASILAN IMPLEMENTASI", "1", 1, "L", true, 0, "")

	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Arial", "", 8.5)
	kepuasanStr := fmt.Sprintf("%.1f / 5.0", l.Kepuasan)
	kesesuaianStr := fmt.Sprintf("%d%%", l.KesesuaianTujuan)
	renderMetaRow("Tingkat Kepuasan Pengguna", kepuasanStr)
	renderMetaRow("Kesesuaian dengan Tujuan", kesesuaianStr)
	pdf.CellFormat(0, 1, "", "T", 1, "L", false, 0, "")
	pdf.Ln(3)

	// 6. TEMUAN UTAMA
	pdf.SetFillColor(240, 244, 248)
	pdf.SetTextColor(0, 51, 102)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 6, "  E. TEMUAN UTAMA EVALUASI", "1", 1, "L", true, 0, "")
	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Arial", "", 8.5)

	if len(l.Temuan) == 0 {
		pdf.CellFormat(0, 6, "  Tidak ada temuan kritis yang dicatat.", "LRB", 1, "L", false, 0, "")
	} else {
		for i, t := range l.Temuan {
			isLast := (i == len(l.Temuan)-1)
			border := "LR"
			if isLast {
				border = "LRB"
			}
			pdf.MultiCell(0, 5, fmt.Sprintf("  %d. %s", i+1, t.Temuan), border, "L", false)
		}
	}
	pdf.Ln(3)

	// 7. REKOMENDASI
	pdf.SetFillColor(240, 244, 248)
	pdf.SetTextColor(0, 51, 102)
	pdf.SetFont("Arial", "B", 9)
	pdf.CellFormat(0, 6, "  F. REKOMENDASI TINDAK LANJUT", "1", 1, "L", true, 0, "")
	pdf.SetTextColor(40, 40, 40)
	pdf.SetFont("Arial", "", 8.5)
	rekom := l.Rekomendasi
	if rekom == "" {
		rekom = "Melaksanakan monitoring berkala dan menindaklanjuti seluruh temuan PIR yang telah ditetapkan."
	}
	pdf.MultiCell(0, 5, rekom, "LRB", "L", false)
	pdf.Ln(6)

	// 8. TANDA TANGAN (2 Kolom: Disusun Oleh & Disetujui Oleh VP TI)
	pdf.SetFont("Arial", "", 8.5)
	pdf.SetTextColor(40, 40, 40)
	pdf.CellFormat(95, 4, "Palembang, "+l.TanggalLaporan, "", 0, "C", false, 0, "")
	pdf.CellFormat(95, 4, "Menyetujui,", "", 1, "C", false, 0, "")

	pdf.CellFormat(95, 4, "Disusun oleh,", "", 0, "C", false, 0, "")
	pdf.CellFormat(95, 4, "VP Teknologi Informasi,", "", 1, "C", false, 0, "")
	pdf.Ln(18)

	pembuat := l.PembuatNama
	if pembuat == "" {
		pembuat = "Staff Divisi TI"
	}
	pdf.SetFont("Arial", "B", 8.5)
	pdf.CellFormat(95, 4, pembuat, "", 0, "C", false, 0, "")
	pdf.CellFormat(95, 4, "Budi Santoso", "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 8)
	pdf.CellFormat(95, 4, "Divisi Teknologi Informasi", "", 0, "C", false, 0, "")
	pdf.CellFormat(95, 4, "PT Pupuk Sriwidjaja Palembang", "", 1, "C", false, 0, "")

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
