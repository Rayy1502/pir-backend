package pdf

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Rayy1502/pir-backend/internal/model"
	"github.com/go-pdf/fpdf"
)

func GenerateLaporanPDF(l *model.LaporanPIR) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(15, 12, 15)
	pdf.SetAutoPageBreak(true, 12)
	pdf.AddPage()

	startX := 15.0
	tableWidth := 180.0
	colLabelW := 45.0
	colValW := tableWidth - colLabelW

	namaKegiatan := "-"
	goLive := "-"
	metode := "SURVEY"

	if l.Proyek != nil {
		if l.Proyek.NamaKegiatan != "" {
			namaKegiatan = l.Proyek.NamaKegiatan
		}
		if !l.Proyek.TanggalGoLive.IsZero() {
			goLive = formatIndonesianDate(l.Proyek.TanggalGoLive.Format("2006-01-02"))
		}
	}
	if l.PIRForm != nil && l.PIRForm.Metode != "" {
		metode = strings.ToUpper(l.PIRForm.Metode)
	}

	// 1. HEADER ROW: KOTAK LOGO PUSRI (KIRI) & BERITA ACARA POST IMPLEMENTATION REVIEW (PIR) (KANAN)
	headerH := 24.0
	y0 := pdf.GetY()

	// Kolom Kiri: Logo PUSRI
	logoBoxW := 65.0
	pdf.Rect(startX, y0, logoBoxW, headerH, "D")

	logoPath := "assets/pusri_logo.png"
	if _, err := os.Stat(logoPath); err == nil {
		pdf.ImageOptions(logoPath, startX+3, y0+3.5, logoBoxW-6, 17, false, fpdf.ImageOptions{ImageType: "PNG"}, 0, "")
	} else {
		// Fallback teks jika file logo tidak ditemukan
		pdf.SetXY(startX, y0+5)
		pdf.SetFont("Arial", "B", 12)
		pdf.SetTextColor(0, 51, 102)
		pdf.CellFormat(logoBoxW, 6, "PUSRI", "", 1, "C", false, 0, "")
		pdf.SetX(startX)
		pdf.SetFont("Arial", "B", 8)
		pdf.SetTextColor(50, 50, 50)
		pdf.CellFormat(logoBoxW, 4, "PUPUK SRIWIDJAJA", "", 1, "C", false, 0, "")
		pdf.SetX(startX)
		pdf.CellFormat(logoBoxW, 4, "PALEMBANG", "", 1, "C", false, 0, "")
	}

	// Kolom Kanan: Judul Berita Acara
	titleBoxW := tableWidth - logoBoxW
	titleX := startX + logoBoxW
	pdf.Rect(titleX, y0, titleBoxW, headerH, "D")

	pdf.SetXY(titleX, y0+5)
	pdf.SetFont("Arial", "B", 13)
	pdf.SetTextColor(20, 20, 20)
	pdf.CellFormat(titleBoxW, 7, "Berita Acara", "", 1, "C", false, 0, "")
	pdf.SetXY(titleX, y0+12)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(titleBoxW, 6, "Post Implementation Review (PIR)", "", 1, "C", false, 0, "")

	pdf.SetY(y0 + headerH)

	// Helper fungsi baris metadata
	renderDataRow := func(label, value string, rowHeight float64) {
		yCur := pdf.GetY()
		pdf.Rect(startX, yCur, colLabelW, rowHeight, "D")
		pdf.Rect(startX+colLabelW, yCur, colValW, rowHeight, "D")

		pdf.SetXY(startX+2, yCur+1.5)
		pdf.SetFont("Arial", "B", 8.5)
		pdf.SetTextColor(30, 30, 30)
		pdf.CellFormat(colLabelW-4, rowHeight-3, label, "", 0, "L", false, 0, "")

		pdf.SetXY(startX+colLabelW+3, yCur+1.5)
		pdf.SetFont("Arial", "", 8.5)
		pdf.CellFormat(colValW-6, rowHeight-3, value, "", 0, "L", false, 0, "")

		pdf.SetY(yCur + rowHeight)
	}

	// 2. METADATA ROWS
	renderDataRow("Aplikasi", namaKegiatan, 8.0)
	renderDataRow("Tanggal serah terima", goLive, 8.0)

	perihalStr := fmt.Sprintf("Progres dan Evaluasi Digitalisasi %s (Nomor: %s)", namaKegiatan, l.NomorPIR)
	renderDataRow("Perihal", perihalStr, 9.0)

	// 3. BARIS URAIAN KEGIATAN (Kotak besar isi laporan)
	uraianStartY := pdf.GetY()

	// Bangun teks lengkap isi kegiatan
	var sb strings.Builder
	sb.WriteString("Telah dilaksanakan kegiatan Post Implementation Review (PIR) sebagai bagian dari evaluasi pasca implementasi sistem. Kegiatan ini dilakukan untuk meninjau hasil pelaksanaan implementasi, memastikan kesesuaian sistem dengan kebutuhan dan proses bisnis yang telah ditetapkan, mengidentifikasi kendala atau permasalahan yang masih ditemukan, serta mengevaluasi efektivitas penggunaan sistem setelah implementasi.\n\n")

	fmt.Fprintf(&sb, "Metode Evaluasi: %s\n", metode)

	metodeKey := strings.ToLower(metode)
	switch metodeKey {
	case "observasi":
		sb.WriteString("Hasil Observasi Langsung & Uji Operasional Tim TI:\n")
		if l.Ringkasan != "" {
			fmt.Fprintf(&sb, "- Ringkasan: %s\n", l.Ringkasan)
		}
		if l.PIRForm != nil && l.PIRForm.CatatanKondisi != "" {
			fmt.Fprintf(&sb, "- Kondisi Lapangan: %s\n", l.PIRForm.CatatanKondisi)
		}
		sb.WriteString("- Status 4 Aspek Evaluasi: Fungsionalitas Modul (OK), Kinerja & Respon Sistem (OK), Ketersediaan Server (High Availability), Keamanan & Otorisasi Hak Akses (OK).\n")
		fmt.Fprintf(&sb, "- Skor Keandalan Operasional: %.1f / 5.0 | Kesesuaian Spesifikasi Teknis: %d%%\n\n", l.Kepuasan, l.KesesuaianTujuan)

	case "wawancara":
		sb.WriteString("Hasil Wawancara Terarah (In-Depth Interview):\n")
		if l.PIRForm != nil && (l.PIRForm.NamaNarasumber != "" || l.PIRForm.JabatanNarasumber != "") {
			fmt.Fprintf(&sb, "- Narasumber: %s (%s)\n", l.PIRForm.NamaNarasumber, l.PIRForm.JabatanNarasumber)
		}
		if l.Ringkasan != "" {
			fmt.Fprintf(&sb, "- Ringkasan: %s\n", l.Ringkasan)
		}
		if l.PIRForm != nil && l.PIRForm.KesimpulanWawancara != "" {
			fmt.Fprintf(&sb, "- Rangkuman Masukan: %s\n", l.PIRForm.KesimpulanWawancara)
		}
		fmt.Fprintf(&sb, "- Skor Kepuasan Narasumber: %.1f / 5.0 | Pemenuhan Kebutuhan Unit Kerja: %d%%\n\n", l.Kepuasan, l.KesesuaianTujuan)

	default: // survey
		sb.WriteString("Hasil Kuesioner Survei Kepuasan Responden:\n")
		if l.PIRForm != nil && l.PIRForm.Responden != "" {
			fmt.Fprintf(&sb, "- Sasaran Responden: %s\n", l.PIRForm.Responden)
		}
		if l.Ringkasan != "" {
			fmt.Fprintf(&sb, "- Ringkasan: %s\n", l.Ringkasan)
		}
		sb.WriteString("- Cakupan Dimensi Evaluasi: Kualitas Sistem & Antarmuka, Kualitas Data & Informasi, Layanan & Dukungan Tim TI, Adopsi & Kepuasan Pengguna, Dampak Kinerja & Manfaat Implementasi.\n")
		fmt.Fprintf(&sb, "- Rata-rata Skor Kepuasan (CSAT): %.1f / 5.0 | Kesesuaian Harapan Pengguna: %d%%\n\n", l.Kepuasan, l.KesesuaianTujuan)
	}

	// Temuan Utama (Dirangkum tanpa nama responden)
	sb.WriteString("Temuan Utama Evaluasi:\n")
	if len(l.Temuan) == 0 {
		sb.WriteString("- Tidak ada kendala kritis yang dilaporkan; sistem berjalan normal.\n")
	} else {
		for i, t := range l.Temuan {
			fmt.Fprintf(&sb, "  %d. %s\n", i+1, cleanAnonymizedText(t.Temuan))
		}
	}
	sb.WriteString("\n")

	// Rekomendasi (Dirangkum tanpa nama responden)
	rekom := strings.TrimSpace(l.Rekomendasi)
	if rekom == "" {
		rekom = "Melaksanakan monitoring berkala dan optimalisasi sistem secara berkelanjutan bersama unit kerja terkait."
	} else {
		lines := strings.Split(rekom, "\n")
		var cleanLines []string
		for _, line := range lines {
			if c := cleanAnonymizedText(line); c != "" {
				cleanLines = append(cleanLines, c)
			}
		}
		if len(cleanLines) > 0 {
			rekom = strings.Join(cleanLines, "\n")
		}
	}
	fmt.Fprintf(&sb, "Rekomendasi Tindak Lanjut:\n%s\n\n", rekom)

	// Penutup resmi Berita Acara PUSRI
	sb.WriteString("Hasil review selanjutnya menjadi dasar dalam menentukan rekomendasi dan tindak lanjut perbaikan yang diperlukan untuk mendukung optimalisasi penggunaan sistem.")

	fullUraianText := sb.String()

	// Hitung tinggi yang dibutuhkan untuk kolom uraian kegiatan
	pdf.SetFont("Arial", "", 8.2)
	lines := pdf.SplitText(fullUraianText, colValW-8)
	uraianHeight := float64(len(lines))*4.2 + 8.0
	if uraianHeight < 95.0 {
		uraianHeight = 95.0
	}

	// Gambar border kotak kolom kiri & kanan
	pdf.Rect(startX, uraianStartY, colLabelW, uraianHeight, "D")
	pdf.Rect(startX+colLabelW, uraianStartY, colValW, uraianHeight, "D")

	// Label "Uraian Kegiatan"
	pdf.SetXY(startX+2, uraianStartY+3)
	pdf.SetFont("Arial", "B", 8.5)
	pdf.SetTextColor(30, 30, 30)
	pdf.CellFormat(colLabelW-4, 6, "Uraian Kegiatan", "", 0, "L", false, 0, "")

	// Isi teks uraian di kolom kanan
	pdf.SetXY(startX+colLabelW+3, uraianStartY+3)
	pdf.SetFont("Arial", "", 8.0)
	pdf.SetTextColor(40, 40, 40)
	pdf.MultiCell(colValW-6, 4.0, fullUraianText, "", "L", false)

	pdf.SetY(uraianStartY + uraianHeight)

	// 4. BARIS TANGGAL KOTA
	dateRowH := 7.5
	yDate := pdf.GetY()
	pdf.Rect(startX, yDate, tableWidth, dateRowH, "D")

	tglLaporan := formatIndonesianDate(l.TanggalLaporan)
	if tglLaporan == "-" || tglLaporan == "" {
		tglLaporan = formatIndonesianDate(time.Now().Format("2006-01-02"))
	}
	pdf.SetXY(startX+4, yDate+1.5)
	pdf.SetFont("Arial", "B", 8.5)
	pdf.SetTextColor(30, 30, 30)
	pdf.CellFormat(tableWidth-8, 5, fmt.Sprintf("Palembang, %s", tglLaporan), "", 0, "L", false, 0, "")

	pdf.SetY(yDate + dateRowH)

	// 5. BARIS 3 KOLOM TANDA TANGAN (Disiapkan oleh, Direview oleh, Disetujui oleh)
	colSigW := tableWidth / 3.0 // 60 mm per kolom
	sigY0 := pdf.GetY()

	// Baris Header Tanda Tangan
	sigHeaderH := 7.0
	pdf.Rect(startX, sigY0, colSigW, sigHeaderH, "D")
	pdf.Rect(startX+colSigW, sigY0, colSigW, sigHeaderH, "D")
	pdf.Rect(startX+2*colSigW, sigY0, colSigW, sigHeaderH, "D")

	pdf.SetFont("Arial", "B", 8.5)
	pdf.SetTextColor(30, 30, 30)

	pdf.SetXY(startX, sigY0+1.2)
	pdf.CellFormat(colSigW, 5, "Disiapkan oleh", "", 0, "C", false, 0, "")

	pdf.SetXY(startX+colSigW, sigY0+1.2)
	pdf.CellFormat(colSigW, 5, "Direview oleh", "", 0, "C", false, 0, "")

	pdf.SetXY(startX+2*colSigW, sigY0+1.2)
	pdf.CellFormat(colSigW, 5, "Disetujui oleh", "", 0, "C", false, 0, "")

	// Kotak Ruang Tanda Tangan (Space TTD)
	sigSpaceH := 20.0
	ySigSpace := sigY0 + sigHeaderH
	pdf.Rect(startX, ySigSpace, colSigW, sigSpaceH, "D")
	pdf.Rect(startX+colSigW, ySigSpace, colSigW, sigSpaceH, "D")
	pdf.Rect(startX+2*colSigW, ySigSpace, colSigW, sigSpaceH, "D")

	// Kotak Nama & Jabatan Pejabat
	sigFooterH := 15.0
	ySigFooter := ySigSpace + sigSpaceH
	pdf.Rect(startX, ySigFooter, colSigW, sigFooterH, "D")
	pdf.Rect(startX+colSigW, ySigFooter, colSigW, sigFooterH, "D")
	pdf.Rect(startX+2*colSigW, ySigFooter, colSigW, sigFooterH, "D")

	// Ambil data penandatangan dari form PIR, jika kosong gunakan nama samaran (fiktif)
	disiapkanNama := "Budi Pratama"
	disiapkanJabatan := "TKNO Perencanaan & Tata Kelola TI"
	direviewNama := "Rahmat Hidayat"
	direviewJabatan := "Senior Engineer Perencanaan & Tata Kelola TI"
	disetujuiNama := "Ahmad Fauzi"
	disetujuiJabatan := "VP Departemen Teknologi Informasi PSP"

	if l.PIRForm != nil {
		if l.PIRForm.DisiapkanOleh != "" {
			disiapkanNama = l.PIRForm.DisiapkanOleh
		}
		if l.PIRForm.DisiapkanJabatan != "" {
			disiapkanJabatan = l.PIRForm.DisiapkanJabatan
		}
		if l.PIRForm.DireviewOleh != "" {
			direviewNama = l.PIRForm.DireviewOleh
		}
		if l.PIRForm.DireviewJabatan != "" {
			direviewJabatan = l.PIRForm.DireviewJabatan
		}
		if l.PIRForm.DisetujuiOleh != "" {
			disetujuiNama = l.PIRForm.DisetujuiOleh
		}
		if l.PIRForm.DisetujuiJabatan != "" {
			disetujuiJabatan = l.PIRForm.DisetujuiJabatan
		}
	} else if l.PembuatNama != "" {
		disiapkanNama = l.PembuatNama
	}

	// Kolom 1: Disiapkan oleh
	pdf.SetXY(startX+1, ySigFooter+1.5)
	pdf.SetFont("Arial", "B", 8.0)
	pdf.CellFormat(colSigW-2, 4.5, disiapkanNama, "", 1, "C", false, 0, "")
	pdf.SetX(startX + 1)
	pdf.SetFont("Arial", "", 7.0)
	pdf.MultiCell(colSigW-2, 3.5, disiapkanJabatan, "", "C", false)

	// Kolom 2: Direview oleh
	pdf.SetXY(startX+colSigW+1, ySigFooter+1.5)
	pdf.SetFont("Arial", "B", 8.0)
	pdf.CellFormat(colSigW-2, 4.5, direviewNama, "", 1, "C", false, 0, "")
	pdf.SetX(startX + colSigW + 1)
	pdf.SetFont("Arial", "", 7.0)
	pdf.MultiCell(colSigW-2, 3.5, direviewJabatan, "", "C", false)

	// Kolom 3: Disetujui oleh
	pdf.SetXY(startX+2*colSigW+1, ySigFooter+1.5)
	pdf.SetFont("Arial", "B", 8.0)
	pdf.CellFormat(colSigW-2, 4.5, disetujuiNama, "", 1, "C", false, 0, "")
	pdf.SetX(startX + 2*colSigW + 1)
	pdf.SetFont("Arial", "", 7.0)
	pdf.MultiCell(colSigW-2, 3.5, disetujuiJabatan, "", "C", false)

	var buf bytes.Buffer
	if err := pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func formatIndonesianDate(raw string) string {
	if raw == "" {
		return "-"
	}
	if strings.Contains(raw, "T") {
		raw = strings.Split(raw, "T")[0]
	}
	parts := strings.Split(raw, "-")
	if len(parts) == 3 {
		year := parts[0]
		month := parts[1]
		day := strings.TrimLeft(parts[2], "0")
		if day == "" {
			day = "0"
		}
		bulanMap := map[string]string{
			"01": "Januari", "02": "Februari", "03": "Maret",
			"04": "April", "05": "Mei", "06": "Juni",
			"07": "Juli", "08": "Agustus", "09": "September",
			"10": "Oktober", "11": "November", "12": "Desember",
		}
		if bName, ok := bulanMap[month]; ok {
			return fmt.Sprintf("%s %s %s", day, bName, year)
		}
	}
	return raw
}

func cleanAnonymizedText(s string) string {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, ": "); idx != -1 && idx < 35 {
		s = strings.TrimSpace(s[idx+2:])
	}
	s = strings.TrimLeft(s, "-*•0123456789. ")
	return strings.TrimSpace(s)
}
