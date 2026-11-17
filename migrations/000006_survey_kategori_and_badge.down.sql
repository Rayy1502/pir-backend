-- Rollback Migrasi 000006
IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('survey_responden') AND name = 'no_badge')
BEGIN
    ALTER TABLE survey_responden DROP COLUMN no_badge;
END;

IF EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('survey_pertanyaan') AND name = 'kategori')
BEGIN
    ALTER TABLE survey_pertanyaan DROP COLUMN kategori;
END;
