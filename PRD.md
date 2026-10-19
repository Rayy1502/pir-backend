# Product Requirements Document (PRD)

## PIR Monitoring — PT PUSRI Palembang

**Versi:** 1.1  
**Tanggal Terakhir Diperbarui:** 19 Okt 2026
**Pemilik Produk:** Divisi Teknologi Informasi PT PUSRI  
**Stakeholder:** Budi Santoso (Staff Divisi TI)  
**Status:** In Progress (Aktif Dikembangkan)

---

## 1. Latar Belakang

PT PUSRI Palembang memiliki banyak proyek pengembangan TI (aplikasi, infrastruktur, end user computing) yang go live setiap tahunnya. Setiap proyek **wajib** menjalani **Post Implementation Review (PIR)** maksimal **2 bulan setelah go live** untuk mengevaluasi keberhasilan implementasi.

Saat ini proses PIR dilakukan secara manual menggunakan dokumen terpisah, sehingga:

- Sulit memantau proyek mana yang belum melakukan PIR
- Tidak ada notifikasi otomatis saat proyek melewati batas waktu
- Pelaporan tidak terstruktur dan sulit diekspor
- Tindak lanjut dari temuan PIR tidak terpantau

**PIR Monitoring** hadir sebagai solusi berbasis web untuk mendigitalkan seluruh siklus PIR: mulai dari pembuatan form, survey, laporan, hingga tindak lanjut.

---

## 2. Tujuan Produk

1. **Memantau** seluruh proyek pengembangan TI dan status PIR-nya secara real-time.
2. **Mengotomatisasi** notifikasi untuk proyek yang terlambat atau mendekati batas waktu PIR.
3. **Menstandarkan** proses PIR melalui form multi-step: Pembuatan → Survey → Laporan → Tindak Lanjut.
4. **Menghasilkan** laporan PIR yang siap cetak/ekspor PDF dengan format resmi PT PUSRI.
5. **Melacak** tindak lanjut temuan PIR hingga selesai dengan PIC dan bukti.

---

## 3. Target Pengguna

| Persona                   | Peran                   | Kebutuhan Utama                              |
| ------------------------- | ----------------------- | -------------------------------------------- |
| **Staff Divisi TI**       | Pembuat & pengelola PIR | Form PIR, laporan, tracking tindak lanjut    |
| **Kepala Divisi TI**      | Approver & reviewer     | Dashboard ringkasan, laporan ekspor          |
| **Tim Pengembang/Vendor** | Narasumber PIR          | Diinput oleh Staff TI (tidak login langsung) |

**Scope v1.0:** hanya untuk Staff Divisi TI (single role).

---

## 4. Fitur Utama

### 4.1 Autentikasi

- Login dengan **Username/NIP** dan **Password**
- Fitur "Ingat saya" dan "Lupa password"
- Redirect ke Dashboard setelah login
- Guard route: halaman dashboard hanya bisa diakses user terautentikasi

### 4.2 Dashboard

- **4 StatCard ringkasan:**
  - Total Projek
  - Projek Selesai
  - PIR Pending (>2 bulan)
  - PIR Selesai
- **Alert banner** merah jika ada proyek terlambat PIR, dengan tombol "Lihat Notifikasi"
- **Tabel Daftar Pengembangan TI:**
  - Kolom: Nama Kegiatan, Kategori, Tgl Go Live, Unit Kerja, Status PIR, Aksi
  - Filter: pencarian nama kegiatan, kategori, status PIR
  - Aksi kondisional: "Buat PIR" (jika terlambat/pending) atau "Detail" (jika selesai/belum jatuh tempo)
- **Donut Chart Statistik Status PIR:** Selesai, Pending, Terlambat, Belum jatuh tempo

### 4.3 Detail Proyek

- Header proyek: nama, kategori, unit kerja, vendor, tanggal go live
- Badge status PIR
- **Timeline tracker** 5 tahap: Go Live → Notifikasi PIR → Survey → Laporan → Tindak Lanjut
- Legend warna: Hijau (selesai), Kuning (pending), Merah (terlambat), Abu (belum dimulai)
- **3 KPI card:** Tingkat Kepuasan, Kesesuaian Tujuan, Progress Tindak Lanjut

### 4.4 Form PIR (Multi-step)

**Step 1 — Pembuatan PIR:**

- Nama Kegiatan (readonly, pre-filled dari proyek)
- Tanggal Go Live (readonly)
- Unit Kerja (readonly)
- Vendor/Tim Pengembang (wajib, validasi)
- Tujuan Awal Pengembangan (textarea, wajib)
- **Metode PIR** (pilih salah satu):
  - **Observasi Langsung** → upload foto + catatan kondisi lapangan
  - **Wawancara** → nama narasumber, jabatan, daftar pertanyaan, kesimpulan
  - **Survey** → field responden
- Tombol: Simpan Draft | Submit & Lanjut

**Step 2 — Form Survey:**

- Progress bar "Tersisa X dari Y pertanyaan"
- Pertanyaan skala kepuasan 1–5
- Pertanyaan kesesuaian hasil dengan tujuan (skala 1–5)
- Textarea kendala
- Textarea saran perbaikan
- Fitur **Tambah Pertanyaan** (custom, dengan tipe: Skala 1–5 / Isian / Pilihan)
- Tombol: Kembali | Kirim Survey

**Step 3 — Laporan PIR:**

- Form input: Ringkasan Hasil, Temuan Utama, Tingkat Kepuasan (angka), Kesesuaian Tujuan (%), Rekomendasi
- **Preview cetak** dengan:
  - Kop surat PT PUSRI — Divisi Teknologi Informasi
  - Nomor dokumen: PIR/TI/[tahun]/[nomor urut]
  - Judul: "Laporan Post Implementation Review"
  - Metadata: Kegiatan, Go Live, Metode
  - Section: Ringkasan, Hasil Utama, Rekomendasi
  - Tanda tangan pembuat
- Tombol: Cetak | Export PDF | Simpan | Lanjut ke Tindak Lanjut

**Step 4 — Tindak Lanjut:**

- Tabel temuan/rekomendasi dengan kolom:
  - Temuan/Rekomendasi
  - PIC (editable)
  - Target Selesai (date picker)
  - Status (dropdown: Open / In Progress / Done)
  - Bukti (upload / lihat file)
- Fitur **Tambah Temuan**
- Dropzone upload bukti (maks 10 MB per file)
- Tombol: Simpan Draft | Simpan & Lihat Detail Projek

### 4.5 Laporan PIR

- Daftar semua laporan PIR yang sudah dibuat
- Filter berdasarkan proyek, tanggal, status
- Aksi: Lihat, Cetak, Export PDF

### 4.6 Tindak Lanjut

- Daftar seluruh tindak lanjut lintas proyek
- Filter berdasarkan status (Open / In Progress / Done)
- Update status dan upload bukti langsung dari list

### 4.7 Notifikasi

- **Filter chip:** Terlambat, Segera
- List notifikasi otomatis untuk proyek yang:
  - **Terlambat** — sudah lewat batas waktu PIR unit kerja (lihat section 4.8)
  - **Segera** — H-7 sebelum batas waktu PIR unit kerja
- Batas waktu dihitung berdasarkan **unit kerja proyek**, bukan global:
  - Proyek unit Pengadaan (2 bulan) → notifikasi H+60
  - Proyek unit Operasi (3 bulan) → notifikasi H+90
- Setiap item: judul proyek, deskripsi, tombol "Buat PIR Sekarang"
- Tombol **Tandai semua dibaca**
- Badge count di sidebar & header

### 4.8 Pengaturan

- Konfigurasi **batas waktu notifikasi PIR per unit kerja** (bukan global)
- Setiap unit kerja bisa punya batas waktu berbeda:
  - Infrastruktur: 2 bulan
  - Layanan: 3 bulan
  - Tata Kelola & Perencanaan: 2 bulan
  - Digitalisasi: 2 bulan
- Form berupa **tabel editable**: setiap baris = 1 unit kerja dengan input angka (1-24 bulan)
- Setiap perubahan disimpan per baris dengan tombol "Simpan"
- Hanya bisa diakses Staff TI
- Perubahan langsung mempengaruhi perhitungan notifikasi otomatis

## 5. Alur Pengguna (User Flow)

Login → Dashboard
├─ Klik "Buat PIR" pada proyek terlambat
│ → Form PIR Step 1 (Pembuatan)
│ → Step 2 (Survey)
│ → Step 3 (Laporan)
│ → Step 4 (Tindak Lanjut)
│ → Kembali ke Detail Proyek (semua step hijau)
│
├─ Klik "Detail" pada proyek selesai
│ → Detail Proyek (timeline + KPI)
│
└─ Klik "Notifikasi" di sidebar
→ List notifikasi → Buat PIR

## 6. Aturan Bisnis

1. **Batas waktu PIR:** 2 bulan setelah tanggal go live (dapat dikonfigurasi di Pengaturan).
2. **Status PIR otomatis:**

- **Belum jatuh tempo** — belum lewat 2 bulan
- **Pending** — dalam proses PIR
- **Terlambat** — lewat 2 bulan & belum PIR
- **Selesai** — PIR lengkap sampai tindak lanjut

3. **Vendor/Tim Pengembang wajib diisi** sebelum bisa lanjut ke step berikutnya.
4. **Metode PIR wajib dipilih** (salah satu dari tiga).
5. **Laporan tidak bisa diekspor** jika step 1–3 belum lengkap.
6. **Tindak lanjut minimal 1 temuan** untuk menyelesaikan PIR.
7. **Notifikasi otomatis** dikirim pada H+60 dan H-7 sebelum batas waktu.
8. **Batas waktu PIR berbeda per unit kerja.** Setiap unit kerja punya konfigurasi
   sendiri di tabel `pengaturan_unit_kerja`. Notifikasi otomatis dihitung
   berdasarkan unit kerja proyek, bukan global. Contoh: proyek di unit Operasi
   punya batas 3 bulan, sementara unit Pengadaan 2 bulan.

---

## 7. Non-Functional Requirements

| Aspek             | Requirement                                          |
| ----------------- | ---------------------------------------------------- |
| **Performance**   | First Contentful Paint < 2s, dashboard load < 3s     |
| **Responsive**    | Desktop-first, mobile-friendly (sidebar collapsible) |
| **Aksesibilitas** | WCAG 2.1 AA, keyboard navigation, aria labels        |
| **Browser**       | Chrome, Edge, Firefox versi terbaru                  |
| **Bahasa**        | Indonesia (UI)                                       |
| **Keamanan**      | Auth wajib, session timeout, HTTPS only              |
| **Data**          | SQL Server via backend Go                            |

---

## 8. Design System

- **Font:** Inter / Plus Jakarta Sans
- **Warna Sidebar:** Gradient `#1e3a8a` → `#059669`
- **Primary Action:** Emerald `#10b981`
- **Background:** `#f3f4f6`
- **Card:** Putih, `rounded-xl`, `shadow-sm`, border `gray-100`
- **Badge Status:**
  - Selesai → `bg-green-100 text-green-700`
  - Pending → `bg-yellow-100 text-yellow-700`
  - Terlambat → `bg-red-100 text-red-700`
  - Belum jatuh tempo → `bg-gray-100 text-gray-700`
- **Sidebar Width:** 260px, fixed
- **Icons:** Lucide React

---

## 9. Tech Stack

### Frontend

- Next.js 15 (App Router) + TypeScript
- Tailwind CSS + shadcn/ui
- Lucide React (icons)
- React Hook Form + Zod (form validation)
- Recharts (donut chart)
- Zustand (state management)

### Backend

- Go 1.22+ dengan Chi router
- SQL Server 2022 (via Docker / local)
- JWT (golang-jwt) + bcrypt
- go-playground/validator

### Infrastructure

- Docker Compose: SQL Server 2022
- GitHub untuk version control
- Deployment: VPS + Nginx (production)

---

## 10. Database Schema

**11 tabel utama:**

1. `users` — data user (NIP, nama, password_hash, role)
2. `proyek` — data proyek pengembangan TI
3. `pir_form` — form PIR step 1
4. `survey_pertanyaan` — pertanyaan survey custom
5. `survey_jawaban` — jawaban survey
6. `laporan_pir` — laporan PIR step 3
7. `laporan_temuan` — temuan laporan
8. `tindak_lanjut` — item tindak lanjut step 4
9. `tindak_lanjut_bukti` — bukti upload
10. `notifikasi` — notifikasi otomatis
11. `pengaturan_unit_kerja` — konfigurasi batas waktu PIR per unit kerja
    - `id` (UNIQUEIDENTIFIER, PK)
    - `unit_kerja` (NVARCHAR(100), UNIQUE)
    - `batas_waktu_bulan` (INT, default 2)
    - `aktif` (BIT, default 1)
    - `updated_at` (DATETIME2)

---

## 11. API Endpoints

### Auth

- `POST /v1/auth/login` — login
- `POST /v1/auth/logout` — logout
- `GET /v1/auth/me` — profil user saat ini

### Proyek

- `GET /v1/proyek` — list proyek (filter: status, kategori, q)
- `GET /v1/proyek/:id` — detail proyek + timeline + KPI
- `POST /v1/proyek` — tambah proyek
- `PUT /v1/proyek/:id` — update proyek
- `DELETE /v1/proyek/:id` — hapus proyek
- `GET /v1/proyek/:id/timeline` — timeline 5 step

### PIR Form

- `POST /v1/pir` — buat PIR form (step 1)
- `GET /v1/pir/:id` — ambil PIR form
- `PUT /v1/pir/:id` — update PIR form

### Survey

- `POST /v1/pir/:id/survey` — submit survey (step 2)
- `GET /v1/pir/:id/survey` — ambil survey

### Laporan

- `POST /v1/laporan` — buat laporan (step 3)
- `GET /v1/laporan` — list laporan
- `GET /v1/laporan/:id` — detail laporan
- `GET /v1/laporan/:id/pdf` — export PDF

### Tindak Lanjut

- `GET /v1/tindak-lanjut` — list tindak lanjut
- `POST /v1/tindak-lanjut` — tambah item
- `PUT /v1/tindak-lanjut/:id` — update item
- `POST /v1/tindak-lanjut/:id/bukti` — upload bukti

### Notifikasi

- `GET /v1/notifikasi` — list notifikasi
- `PUT /v1/notifikasi/:id/read` — tandai dibaca
- `PUT /v1/notifikasi/read-all` — tandai semua dibaca

### Dashboard

- `GET /v1/dashboard/stats` — 4 StatCard
- `GET /v1/dashboard/statistik-pir` — data donut chart

### Pengaturan

- `GET /v1/pengaturan` — list semua unit kerja + batas waktunya
  - Response: array of `{ id, unit_kerja, batas_waktu_bulan, aktif, updated_at }`
- `PUT /v1/pengaturan/{unit}` — update batas waktu 1 unit kerja
  - Body: `{ "batas_waktu_bulan": 3 }`
  - Validasi: 1-24 bulan
  - Response: `{ "message": "Pengaturan berhasil disimpan" }`

---

## 12. Out of Scope (v1.0)

- Multi-role (approver, reviewer)
- Notifikasi email/WhatsApp
- Integrasi langsung ke sistem HRIS
- Mobile app native
- Approval workflow berlapis
- Audit log lengkap

---

## 13. Metrik Keberhasilan

- 100% proyek go live memiliki PIR dalam 2 bulan
- Waktu pembuatan PIR turun dari 2 jam → 20 menit
- 0 proyek terlewat tanpa notifikasi
- Tindak lanjut selesai dalam 30 hari setelah laporan

---

## 14. Roadmap & Status Pengerjaan

### 14.1 Status Ringkasan Milestone

| Milestone | Deskripsi                                                                  |  Backend   |  Frontend  | Status Keseluruhan |
| --------- | -------------------------------------------------------------------------- | :--------: | :--------: | :----------------: |
| **M1**    | Setup Lingkungan (Go, Chi, Docker SQL Server)                |  ✅ Done   |     -      |    **✅ Done**     |
| **M2**    | Autentikasi (Login NIP + Password, JWT, Middleware)                        |  ✅ Done   |  ✅ Done   |    **✅ Done**     |
| **M3**    | Layout Dashboard & Desain Sistem (Sidebar gradien PUSRI, Header responsif) |     -      |  ✅ Done   |    **✅ Done**     |
| **M3.5**  | Konfigurasi Batas Waktu PIR per Unit Kerja                                 |  ✅ Done   |  ✅ Done   |    **✅ Done**     |
| **M4**    | Dashboard Terintegrasi (4 StatCard, Donut Chart, Tabel Proyek)             |  ✅ Done   |  ✅ Done   |    **✅ Done**     |
| **M4.1**  | Navigasi Aksi Dashboard & Detail Proyek (Timeline 5 Tahap + KPI)           | ✅ Done | ✅ Done | **✅ Done** |
| **M5**    | PIR Form Step 1: Pembuatan Form (Metode, Vendor, Tujuan)                   |  ✅ Done   |  ✅ Done   |    **✅ Done**     |
| **M6**    | PIR Form Step 2: Form Survey (Pertanyaan Dinamis & Skala Kepuasan)         | ⏳ Pending | ⏳ Pending |   **⏳ Pending**   |
| **M7**    | PIR Form Step 3: Laporan PIR & Ekspor PDF Format PUSRI                     | ⏳ Pending | ⏳ Pending |   **⏳ Pending**   |
| **M8**    | PIR Form Step 4: Tindak Lanjut & Upload Bukti                              | ⏳ Pending | ⏳ Pending |   **⏳ Pending**   |
| **M9**    | Otomasi Notifikasi & Scheduler Cron (H-7 dan H+batas_waktu)                | ⏳ Pending | ⏳ Pending |   **⏳ Pending**   |
| **M10**   | Testing End-to-End, Security Hardening, & Deployment VPS                   | ⏳ Pending | ⏳ Pending |   **⏳ Pending**   |

---

### 14.2 Rincian Status Fitur (Feature Checklist)

#### 1. Autentikasi & Profil

- [x] Backend: `POST /v1/auth/login` (Bcrypt, JWT token 24 jam)
- [x] Backend: `GET /v1/auth/me` & `POST /v1/auth/logout`
- [x] Frontend: Halaman `/login` dengan validasi form & auth store (Zustand)
- [x] Frontend: Route protection (middleware redirect jika belum login)

#### 2. Dashboard & Master Proyek

- [x] Database: Tabel `proyek` & seed data
- [x] Backend: `GET /v1/proyek` (Filter status, kategori, query pencarian)
- [x] Backend: `GET /v1/dashboard/stats` (4 StatCard)
- [x] Backend: `GET /v1/dashboard/statistik-pir` (Donut chart data)
- [x] Frontend: 4 StatCard (Total, Selesai, Pending >2 bln, PIR Selesai)
- [x] Frontend: Alert Banner merah jika ada proyek lewat batas waktu
- [x] Frontend: Tabel Daftar Proyek dengan badge status warna
- [x] Frontend: Donut Chart status PIR (Recharts)
- [x] Frontend: Routing tombol "Buat PIR" -> Form Step 1
- [x] Frontend: Routing tombol "Detail" -> Halaman Detail Proyek
- [x] Frontend: Routing tombol "Lihat Notifikasi" -> Halaman Notifikasi

#### 3. Pengaturan Batas Waktu per Unit Kerja

- [x] Database: Tabel `pengaturan_unit_kerja` (Infrastruktur, Layanan, Tata Kelola, Digitalisasi)
- [x] Backend: `GET /v1/pengaturan`
- [x] Backend: `PUT /v1/pengaturan/{unit}` (Validasi 1-24 bulan)
- [x] Frontend: Halaman `/pengaturan` (Tabel editable per unit kerja + tombol simpan)

#### 4. Detail Proyek & Visual Timeline

- [x] Backend: `GET /v1/proyek/{id}` (Detail lengkap + status 5 tahap)
- [x] Frontend: Halaman `/proyek/[id]`
- [x] Frontend: Visual Timeline 5 tahap: Go Live → Notifikasi → PIR Review → Laporan → Tindak Lanjut (Hijau, Kuning, Merah, Abu)
- [x] Frontend: 3 KPI Card (Kepuasan, Kesesuaian Tujuan, Progress Tindak Lanjut)

#### 5. PIR Form Step 1: Pembuatan PIR

- [x] Database: Tabel `pir_form`
- [x] Backend: `POST /v1/pir` (Buat draft/submit step 1)
- [x] Backend: `GET /v1/pir/{id}` & `PUT /v1/pir/{id}`
- [x] Frontend: Halaman form wizard Step 1 (Pre-fill data proyek, vendor, tujuan awal, radio metode: Observasi / Wawancara / Survey)

#### 6. PIR Form Step 2: Form Survey

- [x] Database: Tabel `survey_pertanyaan` & `survey_jawaban`
- [ ] Backend: `POST /v1/pir/{id}/survey`
- [ ] Backend: `GET /v1/pir/{id}/survey`
- [ ] Frontend: Form survey kepuasan 1-5, kesesuaian tujuan 1-5, textarea kendala & saran
- [ ] Frontend: Fitur tambah pertanyaan kustom dinamis

#### 7. PIR Form Step 3: Laporan PIR & PDF

- [x] Database: Tabel `laporan_pir` & `laporan_temuan`
- [ ] Backend: `POST /v1/laporan` (Generate nomor urut `PIR/TI/[tahun]/[nomor]`)
- [ ] Backend: `GET /v1/laporan` & `GET /v1/laporan/{id}`
- [ ] Backend: `GET /v1/laporan/{id}/pdf` (Generate PDF resmi)
- [ ] Frontend: Halaman `/laporan-pir` (List laporan)
- [ ] Frontend: Form input temuan utama & ringkasan hasil
- [ ] Frontend: Preview cetak format resmi Kop Surat PT PUSRI Divisi TI

#### 8. PIR Form Step 4: Tindak Lanjut

- [x] Database: Tabel `tindak_lanjut` & `tindak_lanjut_bukti`
- [ ] Backend: `GET /v1/tindak-lanjut` & `POST /v1/tindak-lanjut`
- [ ] Backend: `PUT /v1/tindak-lanjut/{id}` (Update status Open/In Progress/Done & PIC)
- [ ] Backend: `POST /v1/tindak-lanjut/{id}/bukti` (Upload file ke MinIO, max 10MB)
- [ ] Frontend: Halaman `/tindak-lanjut` (Daftar lintas proyek)
- [ ] Frontend: Form tambah temuan, penugasan PIC, datepicker target selesai, dan dropzone upload bukti

#### 9. Notifikasi & Background Scheduler

- [x] Database: Tabel `notifikasi`
- [ ] Backend: Worker scheduler di `internal/scheduler` (Evaluasi harian selisih tanggal go live vs batas waktu unit kerja)
- [ ] Backend: `GET /v1/notifikasi`, `PUT /v1/notifikasi/{id}/read`, `PUT /v1/notifikasi/read-all`
- [x] Frontend: Halaman `/notifikasi` (Filter Terlambat & Segera H-7)
- [x] Frontend: Badge count notifikasi di Header dan Sidebar

---

### 14.3 Log Riwayat Pengerjaan (Changelog)

| Tanggal         | Area                               | Deskripsi Perubahan                                                                                                                                                                         |   Status   |
| --------------- | ---------------------------------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- | :--------: |
| **19 Okt 2026** | Frontend                           | Integrasi `useNotificationStore`: Navigasi klik Bell Header ke `/notifikasi`, badge dinamis Header & Sidebar sesuai notifikasi belum dibaca, tombol Tandai Dibaca & Tandai Semua Dibaca   | ✅ Selesai |
| **19 Okt 2026** | Full Stack (Backend & Frontend) | M4.1: API `GET /v1/proyek/{id}` dengan status 5 tahap & KPI, halaman `/proyek/[id]`, komponen `ProjectTimeline` (5 tahap + legend), `ProjectKPICards`, halaman `/notifikasi`, dan routing aksi dashboard | ✅ Selesai |
| **16 Okt 2026** | Full Stack (DB, Backend, Frontend) | M3.5 UI: Halaman Pengaturan batas waktu unit kerja selesai. M5: API Step 1 PIR (`POST/GET/PUT /v1/pir`), integrasi Redis & MinIO storage, dan UI multi-step wizard Step 1 PIR (`/pir-form`) | ✅ Selesai |
| **16 Okt 2026** | Dokumentasi & Tracking             | Pembaruan PRD v1.1, pembagian milestone detail, dan penyusunan checklist fitur otomatis                                                                                                     | ✅ Selesai |
| **09 Okt 2026** | Frontend                           | Integrasi dashboard dengan backend (StatCard, Donut Chart, Project Table)                                                                                                                   | ✅ Selesai |
| **08 Okt 2026** | Frontend & Backend                 | Pembuatan layout sidebar gradient PUSRI, header responsif, dan endpoint pengaturan unit kerja                                                                                               | ✅ Selesai |
| **07 Okt 2026** | Backend & Database                 | Setup Chi router, koneksi SQL Server Docker, migrasi 11 tabel, dan Auth JWT                                                                                                                 | ✅ Selesai |

---

### 14.4 Instruksi Pembaruan Otomatis untuk Agent

> [!IMPORTANT]
> **Aturan untuk AI Agent**: Setiap kali menyelesaikan task atau pengerjaan fitur baru di proyek ini:
>
> 1. Ubah status item terkait pada tabel **14.1** dan **14.2** dari `⏳ Pending` / `🚧 In Progress` menjadi `✅ Done`.
> 2. Tambahkan baris baru pada tabel **14.3 (Log Riwayat Pengerjaan)** yang mencakup tanggal, area (Frontend/Backend/DB), deskripsi pekerjaan, dan statusnya.
> 3. Perbarui `Tanggal Terakhir Diperbarui` pada header PRD.

---

## 15. Kontak & Repo

- **Pemilik Produk:** (Staff Divisi TI)
- **Repo Backend:** https://github.com/Rayy1502/pir-backend
- **Repo Frontend:** https://github.com/Rayy1502/pir-frontend
- **Database:** SQL Server 2022 (Docker: `pir-sqlserver`)

---

_Dokumen ini adalah sumber kebenaran untuk semua keputusan produk. Setiap perubahan scope harus melalui persetujuan pemilik produk._
