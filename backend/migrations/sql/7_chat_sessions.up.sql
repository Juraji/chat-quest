CREATE TABLE chat_sessions
(
  id                      INTEGER PRIMARY KEY AUTOINCREMENT,
  world_id                INTEGER      NOT NULL REFERENCES worlds (id) ON DELETE CASCADE,
  created_at              TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name                    VARCHAR(100) NOT NULL,
  scenario_id             INTEGER REFERENCES scenarios (id) ON DELETE CASCADE,
  chat_model_id           INTEGER      REFERENCES llm_models (id) ON DELETE SET NULL,
  chat_instruction_id     INTEGER      REFERENCES instruction_templates (id) ON DELETE SET NULL,
  enable_memories         BIT(1)       NOT NULL,
  memories_model_id       INTEGER      REFERENCES llm_models (id) ON DELETE SET NULL,
  memories_instruction_id INTEGER      REFERENCES instruction_templates (id) ON DELETE SET NULL
);

CREATE TABLE chat_participants
(
  chat_session_id INTEGER NOT NULL REFERENCES chat_sessions (id) ON DELETE CASCADE,
  character_id    INTEGER NOT NULL REFERENCES characters (id) ON DELETE CASCADE
);

CREATE TABLE chat_messages
(
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  chat_session_id INTEGER   NOT NULL REFERENCES chat_sessions (id) ON DELETE CASCADE,
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_user         BIT(1)    NOT NULL,
  character_id    INTEGER   REFERENCES characters (id) ON DELETE SET NULL,
  content         TEXT      NOT NULL
);
