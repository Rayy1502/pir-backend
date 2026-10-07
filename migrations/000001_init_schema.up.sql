CREATE TABLE users (
    id            UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    nip           NVARCHAR(20) UNIQUE NOT NULL,
    nama          NVARCHAR(100) NOT NULL,
    password_hash NVARCHAR(255) NOT NULL,
    jabatan       NVARCHAR(100),
    divisi        NVARCHAR(100),
    unit          NVARCHAR(100),
    role          NVARCHAR(20) DEFAULT 'staff_ti',
    is_active     BIT DEFAULT 1,
    created_at    DATETIME2 DEFAULT SYSUTCDATETIME(),
    updated_at    DATETIME2 DEFAULT SYSUTCDATETIME()
);

CREATE TABLE proyek (
    id              UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    nama_kegiatan   NVARCHAR(255) NOT NULL,
    kategori        NVARCHAR(50) NOT NULL,
    unit_kerja      NVARCHAR(100) NOT NULL,
    vendor          NVARCHAR(255),
    tanggal_go_live DATE NOT NULL,
    status_pir      NVARCHAR(30) NOT NULL DEFAULT 'belum-jatuh-tempo',
    created_at      DATETIME2 DEFAULT SYSUTCDATETIME(),
    updated_at      DATETIME2 DEFAULT SYSUTCDATETIME()
);
CREATE INDEX idx_proyek_status ON proyek(status_pir);
CREATE INDEX idx_proyek_go_live ON proyek(tanggal_go_live);

CREATE TABLE pir_form (
    id                   UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    proyek_id            UNIQUEIDENTIFIER REFERENCES proyek(id) ON DELETE CASCADE,
    tujuan_awal          NVARCHAR(MAX),
    metode               NVARCHAR(20),
    catatan_kondisi      NVARCHAR(MAX),
    nama_narasumber      NVARCHAR(100),
    jabatan_narasumber   NVARCHAR(100),
    daftar_pertanyaan    NVARCHAR(MAX),
    kesimpulan_wawancara NVARCHAR(MAX),
    responden            NVARCHAR(255),
    status               NVARCHAR(20) DEFAULT 'draft',
    created_by           UNIQUEIDENTIFIER REFERENCES users(id),
    created_at           DATETIME2 DEFAULT SYSUTCDATETIME(),
    updated_at           DATETIME2 DEFAULT SYSUTCDATETIME()
);

CREATE TABLE survey_pertanyaan (
    id            UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    pir_form_id   UNIQUEIDENTIFIER REFERENCES pir_form(id) ON DELETE CASCADE,
    pertanyaan    NVARCHAR(MAX) NOT NULL,
    tipe          NVARCHAR(20) NOT NULL,
    urutan        INT NOT NULL,
    created_at    DATETIME2 DEFAULT SYSUTCDATETIME()
);

CREATE TABLE survey_jawaban (
    id            UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    pertanyaan_id UNIQUEIDENTIFIER REFERENCES survey_pertanyaan(id) ON DELETE CASCADE,
    jawaban_teks  NVARCHAR(MAX),
    jawaban_angka INT,
    created_at    DATETIME2 DEFAULT SYSUTCDATETIME()
);

CREATE TABLE laporan_pir (
    id                  UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    nomor_pir           NVARCHAR(50) UNIQUE NOT NULL,
    proyek_id           UNIQUEIDENTIFIER REFERENCES proyek(id),
    pir_form_id         UNIQUEIDENTIFIER REFERENCES pir_form(id),
    ringkasan           NVARCHAR(MAX),
    kepuasan            DECIMAL(3,1),
    kesesuaian_tujuan   INT,
    rekomendasi         NVARCHAR(MAX),
    tanggal_laporan     DATE NOT NULL,
    pembuat_id          UNIQUEIDENTIFIER REFERENCES users(id),
    created_at          DATETIME2 DEFAULT SYSUTCDATETIME()
);

CREATE TABLE laporan_temuan (
    id            UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    laporan_id    UNIQUEIDENTIFIER REFERENCES laporan_pir(id) ON DELETE CASCADE,
    temuan        NVARCHAR(MAX) NOT NULL,
    urutan        INT
);

CREATE TABLE tindak_lanjut (
    id                UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    laporan_id        UNIQUEIDENTIFIER REFERENCES laporan_pir(id) ON DELETE CASCADE,
    temuan            NVARCHAR(MAX) NOT NULL,
    pic               NVARCHAR(100),
    target_selesai    DATE,
    status            NVARCHAR(20) DEFAULT 'open',
    created_at        DATETIME2 DEFAULT SYSUTCDATETIME(),
    updated_at        DATETIME2 DEFAULT SYSUTCDATETIME()
);

CREATE TABLE tindak_lanjut_bukti (
    id                UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    tindak_lanjut_id  UNIQUEIDENTIFIER REFERENCES tindak_lanjut(id) ON DELETE CASCADE,
    nama_file         NVARCHAR(255),
    url               NVARCHAR(MAX),
    ukuran_bytes      BIGINT,
    uploaded_at       DATETIME2 DEFAULT SYSUTCDATETIME()
);

CREATE TABLE notifikasi (
    id              UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    proyek_id       UNIQUEIDENTIFIER REFERENCES proyek(id) ON DELETE CASCADE,
    tipe            NVARCHAR(20) NOT NULL,
    judul           NVARCHAR(255),
    deskripsi       NVARCHAR(MAX),
    dibaca          BIT DEFAULT 0,
    created_at      DATETIME2 DEFAULT SYSUTCDATETIME()
);
CREATE INDEX idx_notifikasi_dibaca ON notifikasi(dibaca);

CREATE TABLE pengaturan (
    [key] NVARCHAR(100) PRIMARY KEY,
    value NVARCHAR(MAX) NOT NULL
);