CREATE TABLE connection_profiles
(
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    provider_type VARCHAR(50)   NOT NULL,
    base_url      VARCHAR(255)  NOT NULL,
    api_key       VARCHAR(1024) NOT NULL
);

CREATE TABLE llm_models
(
    id                    INTEGER PRIMARY KEY AUTOINCREMENT,
    connection_profile_id INTEGER       NOT NULL,
    model_id              VARCHAR(1024) NOT NULL,
    temperature           FLOAT         NOT NULL DEFAULT 0.7,
    max_tokens            INTEGER       NOT NULL DEFAULT 256,
    top_p                 FLOAT         NOT NULL DEFAULT 0.95,
    stream                BIT(1)        NOT NULL DEFAULT 0,
    stop                  VARCHAR(2048) NOT NULL DEFAULT '',

    constraint fk_lm__connection_profile FOREIGN KEY (connection_profile_id) REFERENCES connection_profiles (id) ON DELETE CASCADE
)
