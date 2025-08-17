CREATE TABLE chat_sessions
(
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  world_id        INTEGER      NOT NULL REFERENCES worlds (id) ON DELETE CASCADE,
  created_at      TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name            VARCHAR(100) NOT NULL,
  scenario_id     INTEGER REFERENCES scenarios (id) ON DELETE CASCADE,
  enable_memories BIT(1)       NOT NULL
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
  is_system BIT(1) NOT NULL,
  character_id    INTEGER   REFERENCES characters (id) ON DELETE SET NULL,
  content         TEXT      NOT NULL
);
