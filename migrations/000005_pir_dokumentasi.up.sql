-- Migrasi 000005: Tabel pir_dokumentasi untuk upload bukti observasi & wawancara
IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'pir_dokumentasi')
BEGIN
    CREATE TABLE pir_dokumentasi (
        id           UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
        pir_form_id  UNIQUEIDENTIFIER REFERENCES pir_form(id) ON DELETE CASCADE,
        kategori     NVARCHAR(50) NOT NULL DEFAULT 'observasi', -- 'observasi' atau 'wawancara'
        nama_file    NVARCHAR(255) NOT NULL,
        url          NVARCHAR(MAX) NOT NULL,
        ukuran_bytes BIGINT DEFAULT 0,
        keterangan   NVARCHAR(MAX) NULL,
        uploaded_at  DATETIME2 DEFAULT SYSUTCDATETIME()
    );

    CREATE INDEX idx_pir_dok_form ON pir_dokumentasi(pir_form_id);
    CREATE INDEX idx_pir_dok_kategori ON pir_dokumentasi(kategori);
END;
