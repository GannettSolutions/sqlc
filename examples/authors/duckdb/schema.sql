CREATE SEQUENCE authors_id_seq;

CREATE TABLE authors (
  id               BIGINT PRIMARY KEY DEFAULT nextval('authors_id_seq'),
  name             TEXT NOT NULL,
  bio              TEXT,
  data_json        JSON NOT NULL DEFAULT '{}',
  preferences_json JSON NOT NULL DEFAULT '[]'
);
