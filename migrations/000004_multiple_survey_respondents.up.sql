-- Migrasi 000004: Dukungan Multiple Responden Survey Kuesioner
IF NOT EXISTS (SELECT * FROM sys.tables WHERE name = 'survey_responden')
BEGIN
    CREATE TABLE survey_responden (
        id                 UNIQUEIDENTIFIER PRIMARY KEY DEFAULT NEWSEQUENTIALID(),
        pir_form_id        UNIQUEIDENTIFIER NOT NULL REFERENCES pir_form(id) ON DELETE CASCADE,
        nama               NVARCHAR(255) NOT NULL,
        jabatan            NVARCHAR(255) NULL,
        unit_kerja         NVARCHAR(255) NULL,
        kepuasan           INT NOT NULL,          -- Skala 1-5
        kesesuaian_tujuan  INT NOT NULL,          -- Skala 1-5
        kendala            NVARCHAR(MAX) NULL,
        saran              NVARCHAR(MAX) NULL,
        created_at         DATETIME2 DEFAULT SYSUTCDATETIME()
    );
    CREATE INDEX idx_survey_responden_pir ON survey_responden(pir_form_id);
END;

-- Tambahkan kolom responden_id di survey_jawaban jika belum ada
IF NOT EXISTS (SELECT * FROM sys.columns WHERE object_id = OBJECT_ID('survey_jawaban') AND name = 'responden_id')
BEGIN
    ALTER TABLE survey_jawaban ADD responden_id UNIQUEIDENTIFIER NULL REFERENCES survey_responden(id) ON DELETE NO ACTION;
END;
