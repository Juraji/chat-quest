CREATE TABLE tags
(
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  label     VARCHAR(50) NOT NULL,
  lowercase VARCHAR(50) NOT NULL,

  CONSTRAINT uk_t__label UNIQUE (label)
);

CREATE TABLE characters
(
  id         INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name       VARCHAR(100) NOT NULL,
  favorite   BIT(1)       NOT NULL DEFAULT 0,
  avatar_url TEXT                  DEFAULT NULL
);

CREATE TABLE character_details
(
  character_id        INTEGER PRIMARY KEY NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
  appearance          TEXT                         DEFAULT NULL,
  personality         TEXT                         DEFAULT NULL,
  history             TEXT                         DEFAULT NULL,
  group_talkativeness FLOAT               NOT NULL DEFAULT 0.5

);

CREATE TABLE character_tags
(
  character_id INTEGER NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
  tag_id       INTEGER NOT NULL REFERENCES tags (id)
);

CREATE TABLE character_dialogue_examples
(
  character_id INTEGER NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
  text         TEXT    NOT NULL
);

CREATE TABLE character_greetings
(
  character_id INTEGER NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
  text         TEXT    NOT NULL
);

CREATE TABLE character_group_greetings
(
  character_id INTEGER NOT NULL REFERENCES characters (id) ON DELETE CASCADE,
  text         TEXT    NOT NULL
);
