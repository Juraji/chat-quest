CREATE TABLE instructions
(
  id                  INTEGER PRIMARY KEY AUTOINCREMENT,
  name                VARCHAR(100) NOT NULL,
  type                VARCHAR(50)  NOT NULL,

  -- Model Settings
  temperature         FLOAT        NOT NULL,
  max_tokens          INTEGER      NOT NULL,
  top_p               FLOAT        NOT NULL,
  presence_penalty    FLOAT        NOT NULL,
  frequency_penalty   FLOAT        NOT NULL,
  stream              BIT(1)       NOT NULL,
  stop_sequences      VARCHAR(1024),
  include_reasoning BIT(1) NOT NULL,

  -- Parsing
  reasoning_prefix    VARCHAR(50)  NOT NULL,
  reasoning_suffix    VARCHAR(50)  NOT NULL,
  character_id_prefix VARCHAR(50)  NOT NULL,
  character_id_suffix VARCHAR(50)  NOT NULL,

  -- Prompt Templates
  system_prompt       TEXT         NOT NULL,
  world_setup         TEXT         NOT NULL,
  instruction         TEXT         NOT NULL
);
