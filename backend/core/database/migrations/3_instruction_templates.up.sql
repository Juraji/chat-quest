CREATE TABLE instruction_templates
(
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          VARCHAR(100) NOT NULL,
  type          VARCHAR(50)  NOT NULL,
  temperature   FLOAT DEFAULT NULL,
  system_prompt TEXT         NOT NULL,
  world_setup   TEXT         NOT NULL,
  instruction   TEXT         NOT NULL
);
