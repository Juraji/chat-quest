CREATE TABLE worlds
(
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        VARCHAR(100) NOT NULL,
  description TEXT DEFAULT NULL,
  avatar_url  TEXT DEFAULT NULL
);
