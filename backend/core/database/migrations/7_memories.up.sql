CREATE TABLE memories
(
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  world_id           INTEGER   NOT NULL REFERENCES worlds (id) ON DELETE CASCADE,
  chat_session_id    INTEGER REFERENCES chat_sessions (id) ON DELETE CASCADE,
  character_id       INTEGER   NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
  created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  content            TEXT      NOT NULL,
  embedding          BLOB      NOT NULL,
  embedding_model_id INTEGER   REFERENCES llm_models (id) ON DELETE SET NULL
);

ALTER TABLE chat_messages
  ADD COLUMN memory_id INTEGER REFERENCES memories (id) ON DELETE SET NULL;

CREATE TABLE memory_preferences
(
  id                      INTEGER NOT NULL PRIMARY KEY,
  memories_model_id       INTEGER REFERENCES llm_models (id) ON DELETE SET NULL,
  memories_instruction_id INTEGER REFERENCES instruction_templates (id) ON DELETE SET NULL,
  embedding_model_id      INTEGER REFERENCES llm_models (id) ON DELETE SET NULL,
  memory_min_p FLOAT NOT NULL,
  memory_trigger_after    INTEGER NOT NULL,
  memory_window_size      INTEGER NOT NULL
);

-- Insert default record
INSERT INTO memory_preferences (id, memories_model_id, memories_instruction_id, embedding_model_id, memory_min_p,
                                memory_trigger_after, memory_window_size)
VALUES (0,
        NULL,
        NULL,
        NULL,
        0.95,
        15,
        3);
