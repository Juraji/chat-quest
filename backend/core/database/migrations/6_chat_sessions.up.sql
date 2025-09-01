CREATE TABLE chat_sessions
(
  id                        INTEGER PRIMARY KEY AUTOINCREMENT,
  world_id                  INTEGER      NOT NULL REFERENCES worlds (id) ON DELETE CASCADE,
  created_at                TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name                      VARCHAR(100) NOT NULL,
  scenario_id               INTEGER REFERENCES scenarios (id) ON DELETE CASCADE,
  generate_memories BIT(1) NOT NULL,
  use_memories      BIT(1) NOT NULL,
  pause_automatic_responses BIT(1)       NOT NULL
);

create table chat_participants
(
  chat_session_id INTEGER                             NOT NULL REFERENCES chat_sessions ON DELETE CASCADE,
  character_id    INTEGER                             NOT NULL REFERENCES characters ON DELETE CASCADE,
  added_on        TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  removed_on      TIMESTAMP DEFAULT NULL,
  constraint chat_participants_pk primary key (chat_session_id, character_id)
);

CREATE TABLE chat_messages
(
  id              INTEGER PRIMARY KEY AUTOINCREMENT,
  chat_session_id INTEGER   NOT NULL REFERENCES chat_sessions (id) ON DELETE CASCADE,
  created_at      TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  is_user         BIT(1)    NOT NULL,
  is_generating BIT(1) NOT NULL,
  is_archived   BIT(1) NOT NULL DEFAULT FALSE,
  character_id    INTEGER   REFERENCES characters (id) ON DELETE SET NULL,
  content         TEXT      NOT NULL
);

CREATE INDEX idx_chat_messages__is_archived ON chat_messages (is_archived);
