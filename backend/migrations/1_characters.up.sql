CREATE TABLE tags
(
  id        INTEGER PRIMARY KEY AUTOINCREMENT,
  label     VARCHAR(50) NOT NULL,
  lowercase VARCHAR(50) NOT NULL,

  CONSTRAINT uk_t__label UNIQUE (label)
);

CREATE TABLE characters
(
  id                  INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at          TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
  name                VARCHAR(100) NOT NULL,
  favorite            BIT(1)       NOT NULL DEFAULT 0,

  -- base
  appearance          TEXT                  DEFAULT NULL,
  personality         TEXT                  DEFAULT NULL,
  history             TEXT                  DEFAULT NULL,

  -- chat
  scenario            TEXT                  DEFAULT NULL,
  first_message       TEXT                  DEFAULT NULL,
  group_talkativeness FLOAT        NOT NULL DEFAULT 0.5
);

CREATE TABLE character_tags
(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  character_id INTEGER NOT NULL,
  tag_id       INTEGER NOT NULL,

  CONSTRAINT fk_ct__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE,
  CONSTRAINT fk_ct__tag FOREIGN KEY (tag_id) REFERENCES tags (id)
);

CREATE TABLE character_likely_actions
(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  character_id INTEGER NOT NULL,
  action       TEXT    NOT NULL,

  CONSTRAINT fk_cla__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);

CREATE TABLE character_unlikely_actions
(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  character_id INTEGER NOT NULL,
  action       TEXT    NOT NULL,

  CONSTRAINT fk_cua__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);

CREATE TABLE character_dialogue_examples
(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  character_id INTEGER NOT NULL,
  example      TEXT    NOT NULL,

  CONSTRAINT fk_cde__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);

CREATE TABLE character_alternate_greetings
(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  character_id INTEGER NOT NULL,
  greeting     TEXT    NOT NULL,

  CONSTRAINT fk_cag__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);

CREATE TABLE character_group_greetings
(
  id           INTEGER PRIMARY KEY AUTOINCREMENT,
  character_id INTEGER NOT NULL,
  greeting     TEXT    NOT NULL,

  CONSTRAINT fk_cgg__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);
