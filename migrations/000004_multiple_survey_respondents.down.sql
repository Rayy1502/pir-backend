IF EXISTS (SELECT * FROM sys.tables WHERE name = 'survey_responden')
BEGIN
    DROP TABLE survey_responden;
END;
