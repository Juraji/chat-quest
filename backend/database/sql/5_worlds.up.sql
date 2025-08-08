CREATE TABLE worlds
(
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        VARCHAR(100) NOT NULL,
  description TEXT NOT NULL,
  avatar_url  TEXT DEFAULT NULL
);

CREATE TABLE chat_preferences
(
  id                  INTEGER NOT NULL PRIMARY KEY,
  chat_model_id       INTEGER REFERENCES llm_models (id) ON DELETE SET NULL,
  chat_instruction_id INTEGER REFERENCES instruction_templates (id) ON DELETE SET NULL
);

-- Just making sure our default record exists
INSERT INTO chat_preferences (id, chat_model_id, chat_instruction_id)
VALUES (0,
        NULL,
        1);
