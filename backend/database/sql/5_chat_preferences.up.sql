CREATE TABLE chat_preferences
(
  id                   INTEGER NOT NULL PRIMARY KEY,
  chat_model_id           INTEGER REFERENCES llm_models (id) ON DELETE SET NULL,
  chat_instruction_id     INTEGER REFERENCES instruction_templates (id) ON DELETE SET NULL,
  memories_model_id       INTEGER REFERENCES llm_models (id) ON DELETE SET NULL,
  memories_instruction_id INTEGER REFERENCES instruction_templates (id) ON DELETE SET NULL,
  embedding_model_id   INTEGER REFERENCES llm_models (id) ON DELETE SET NULL,
  memory_top_p         FLOAT   NOT NULL,
  memory_trigger_after INTEGER NOT NULL,
  memory_window_size   INTEGER NOT NULL
);

-- Just making sure our default record exists
INSERT INTO chat_preferences (id, chat_model_id, chat_instruction_id, memories_model_id,
                              memories_instruction_id, embedding_model_id,
                              memory_top_p, memory_trigger_after, memory_window_size)
VALUES (0,
        NULL,
        NULL,
        NULL,
        NULL,
        NULL,
        0.95,
        15,
        4);
