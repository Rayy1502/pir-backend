-- Migrasi 000006: Tambah No. Badge di survey_responden dan Kategori di survey_pertanyaan
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('survey_responden') AND name = 'no_badge')
BEGIN
    ALTER TABLE survey_responden ADD no_badge NVARCHAR(100) NULL;
END;

IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('survey_pertanyaan') AND name = 'kategori')
BEGIN
    ALTER TABLE survey_pertanyaan ADD kategori NVARCHAR(100) NULL;
END;
