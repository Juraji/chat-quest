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
  character_id        INTEGER PRIMARY KEY NOT NULL,
  appearance          TEXT                         DEFAULT NULL,
  personality         TEXT                         DEFAULT NULL,
  history             TEXT                         DEFAULT NULL,
  group_talkativeness FLOAT               NOT NULL DEFAULT 0.5,

  CONSTRAINT fk_ct__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);

CREATE TABLE character_tags
(
  character_id INTEGER NOT NULL,
  tag_id       INTEGER NOT NULL,

  CONSTRAINT fk_ct__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE,
  CONSTRAINT fk_ct__tag FOREIGN KEY (tag_id) REFERENCES tags (id)
);

CREATE TABLE character_dialogue_examples
(
  character_id INTEGER NOT NULL,
  text         TEXT    NOT NULL,

  CONSTRAINT fk_cde__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);

CREATE TABLE character_greetings
(
  character_id INTEGER NOT NULL,
  text         TEXT    NOT NULL,

  CONSTRAINT fk_cag__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);

CREATE TABLE character_group_greetings
(
  character_id INTEGER NOT NULL,
  text         TEXT    NOT NULL,

  CONSTRAINT fk_cgg__character FOREIGN KEY (character_id) REFERENCES characters (id) ON DELETE CASCADE
);
