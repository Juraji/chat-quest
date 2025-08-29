CREATE TABLE memories
(
  id                 INTEGER PRIMARY KEY AUTOINCREMENT,
  world_id           INTEGER   NOT NULL REFERENCES worlds (id) ON DELETE CASCADE,
  character_id INTEGER REFERENCES characters (id) ON DELETE CASCADE,
  created_at         TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  content            TEXT      NOT NULL,
  embedding    BLOB,
  embedding_model_id INTEGER   REFERENCES llm_models (id) ON DELETE SET NULL
);
