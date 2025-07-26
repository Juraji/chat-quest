CREATE TABLE system_prompts
(
  id     INTEGER PRIMARY KEY AUTOINCREMENT,
  name   VARCHAR(100) NOT NULL,
  prompt TEXT         NOT NULL
);
