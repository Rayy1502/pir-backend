DROP TABLE IF EXISTS pengaturan_unit_kerja;

-- Kembalikan tabel pengaturan lama
CREATE TABLE pengaturan (
    [key] NVARCHAR(100) PRIMARY KEY,
    value NVARCHAR(MAX) NOT NULL
);