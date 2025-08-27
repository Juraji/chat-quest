CREATE TABLE preferences
(
  id                      INTEGER PRIMARY KEY,
  -- Chat
  chat_model_id           INTEGER REFERENCES llm_models (id) ON DELETE SET NULL,
  chat_instruction_id     INTEGER REFERENCES instruction_templates (id) ON DELETE SET NULL,
  -- Embedding
  embedding_model_id      INTEGER REFERENCES llm_models (id) ON DELETE SET NULL,
  -- Memories
  memories_model_id       INTEGER REFERENCES llm_models (id) ON DELETE SET NULL,
  memories_instruction_id INTEGER REFERENCES instruction_templates (id) ON DELETE SET NULL,
  memory_min_p            FLOAT   NOT NULL,
  memory_trigger_after    INTEGER NOT NULL,
  memory_window_size      INTEGER NOT NULL
);

-- Insert default record
INSERT INTO preferences (id, chat_model_id, memories_model_id, embedding_model_id, chat_instruction_id,
                         memories_instruction_id, memory_min_p, memory_trigger_after, memory_window_size)
VALUES (0,
        NULL,
        (SELECT id FROM instruction_templates WHERE type = 'CHAT' LIMIT 1),
        NULL,
        NULL,
        (SELECT id FROM instruction_templates WHERE type = 'MEMORIES' LIMIT 1),
        0.95,
        15,
        10);
