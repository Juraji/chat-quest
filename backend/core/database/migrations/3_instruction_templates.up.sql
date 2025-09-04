CREATE TABLE instructions
(
  id                INTEGER PRIMARY KEY AUTOINCREMENT,
  name              VARCHAR(100) NOT NULL,
  type              VARCHAR(50)  NOT NULL,

  -- Model Settings
  temperature       FLOAT        NOT NULL,
  max_tokens        INTEGER      NOT NULL,
  top_p             FLOAT        NOT NULL,
  presence_penalty  FLOAT        NOT NULL,
  frequency_penalty FLOAT        NOT NULL,
  stream            BIT(1)       NOT NULL,
  stop_sequences    VARCHAR(1024),

  -- Templates
  system_prompt     TEXT         NOT NULL,
  world_setup       TEXT         NOT NULL,
  instruction       TEXT         NOT NULL
);
