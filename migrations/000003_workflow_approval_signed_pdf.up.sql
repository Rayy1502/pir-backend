-- Tambah approved_by & approved_at di pir_form
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'approved_by')
BEGIN
    ALTER TABLE pir_form ADD approved_by NVARCHAR(100) NULL;
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'approved_at')
BEGIN
    ALTER TABLE pir_form ADD approved_at DATETIME2 NULL;
END;

-- Tambah file_ttd_url, file_ttd_nama, status_laporan, signed_at di laporan_pir
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('laporan_pir') AND name = 'file_ttd_url')
BEGIN
    ALTER TABLE laporan_pir ADD file_ttd_url NVARCHAR(MAX) NULL;
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('laporan_pir') AND name = 'file_ttd_nama')
BEGIN
    ALTER TABLE laporan_pir ADD file_ttd_nama NVARCHAR(255) NULL;
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('laporan_pir') AND name = 'status_laporan')
BEGIN
    ALTER TABLE laporan_pir ADD status_laporan NVARCHAR(20) DEFAULT 'draft';
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('laporan_pir') AND name = 'signed_at')
BEGIN
    ALTER TABLE laporan_pir ADD signed_at DATETIME2 NULL;
END;
