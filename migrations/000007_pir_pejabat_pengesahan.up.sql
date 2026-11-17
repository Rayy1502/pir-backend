-- Tambah kolom pengesahan berita acara pada tabel pir_form
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'disiapkan_oleh')
BEGIN
    ALTER TABLE pir_form ADD disiapkan_oleh NVARCHAR(150) NULL DEFAULT 'Budi Pratama';
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'disiapkan_jabatan')
BEGIN
    ALTER TABLE pir_form ADD disiapkan_jabatan NVARCHAR(150) NULL DEFAULT 'TKNO Perencanaan & Tata Kelola TI';
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'direview_oleh')
BEGIN
    ALTER TABLE pir_form ADD direview_oleh NVARCHAR(150) NULL DEFAULT 'Rahmat Hidayat';
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'direview_jabatan')
BEGIN
    ALTER TABLE pir_form ADD direview_jabatan NVARCHAR(150) NULL DEFAULT 'Senior Engineer Perencanaan & Tata Kelola TI';
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'disetujui_oleh')
BEGIN
    ALTER TABLE pir_form ADD disetujui_oleh NVARCHAR(150) NULL DEFAULT 'Ahmad Fauzi';
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'disetujui_jabatan')
BEGIN
    ALTER TABLE pir_form ADD disetujui_jabatan NVARCHAR(150) NULL DEFAULT 'VP Departemen Teknologi Informasi PSP';
END;
