CREATE TABLE connection_profiles
(
  id            INTEGER PRIMARY KEY AUTOINCREMENT,
  name          VARCHAR(100)  NOT NULL,
  provider_type VARCHAR(50)   NOT NULL,
  base_url      VARCHAR(255)  NOT NULL,
  api_key       VARCHAR(1024) NOT NULL
);

CREATE TABLE llm_models
(
  id                    INTEGER PRIMARY KEY AUTOINCREMENT,
  connection_profile_id INTEGER      NOT NULL REFERENCES connection_profiles (id) ON DELETE CASCADE,
  model_id              VARCHAR(255) NOT NULL,
  temperature           FLOAT        NOT NULL,
  max_tokens            INTEGER      NOT NULL,
  top_p                 FLOAT        NOT NULL,
  stream                BIT(1)       NOT NULL,
  stop_sequences        VARCHAR(1024),
  disabled              BIT(1)       NOT NULL
)
