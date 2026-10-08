-- Hapus tabel pengaturan lama
DROP TABLE IF EXISTS pengaturan;

-- Buat tabel pengaturan per unit kerja
CREATE TABLE pengaturan_unit_kerja (
    id                  UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
    unit_kerja          NVARCHAR(100) NOT NULL UNIQUE,
    batas_waktu_bulan   INT NOT NULL DEFAULT 2,
    aktif               BIT DEFAULT 1,
    updated_at          DATETIME2 DEFAULT SYSUTCDATETIME()
);

-- Seed data unit kerja
INSERT INTO pengaturan_unit_kerja (unit_kerja, batas_waktu_bulan) VALUES
('Infrastruktur', 2),
('Layanan', 3),
('Tata Kelola & Perencanaan', 2),
('Digitalisasi', 2);