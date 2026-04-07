ALTER TABLE preferences
  ADD COLUMN title_generation_model_id INTEGER REFERENCES llm_models (id) ON DELETE SET NULL;
ALTER TABLE preferences
  ADD COLUMN title_generation_instruction_id INTEGER REFERENCES instructions (id) ON DELETE SET NULL;
ALTER TABLE preferences
  ADD COLUMN title_generation_message_window INTEGER NOT NULL DEFAULT 10;
