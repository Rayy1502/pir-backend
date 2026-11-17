IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'disiapkan_oleh')
BEGIN
    ALTER TABLE pir_form DROP COLUMN disiapkan_oleh;
END;
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'disiapkan_jabatan')
BEGIN
    ALTER TABLE pir_form DROP COLUMN disiapkan_jabatan;
END;
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'direview_oleh')
BEGIN
    ALTER TABLE pir_form DROP COLUMN direview_oleh;
END;
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'direview_jabatan')
BEGIN
    ALTER TABLE pir_form DROP COLUMN direview_jabatan;
END;
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'disetujui_oleh')
BEGIN
    ALTER TABLE pir_form DROP COLUMN disetujui_oleh;
END;
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('pir_form') AND name = 'disetujui_jabatan')
BEGIN
    ALTER TABLE pir_form DROP COLUMN disetujui_jabatan;
END;
