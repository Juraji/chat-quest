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
