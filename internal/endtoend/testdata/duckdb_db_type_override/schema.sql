CREATE TABLE records (
  id INTEGER PRIMARY KEY,
  data_json JSON NOT NULL,
  preferences_json JSON NOT NULL
);

CREATE TABLE reviews (
  id INTEGER PRIMARY KEY,
  required_json JSON NOT NULL,
  optional_json JSON
);

CREATE VIEW direct_json_view AS
SELECT id, required_json, optional_json FROM reviews;

CREATE VIEW cast_json_view AS
SELECT id,
       CAST(required_json AS JSON) AS required_json,
       CAST(optional_json AS JSON) AS optional_json,
       CAST(coalesce(optional_json, CAST('{}' AS JSON)) AS JSON) AS review_json,
       CAST(required_json AS VARCHAR) AS plain_text
FROM reviews;

CREATE VIEW nested_json_view AS
SELECT * FROM cast_json_view;
