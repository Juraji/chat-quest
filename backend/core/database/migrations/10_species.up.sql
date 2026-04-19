CREATE TABLE species
(
  id          INTEGER PRIMARY KEY AUTOINCREMENT,
  name        VARCHAR(100) NOT NULL,
  description TEXT         NOT NULL,
  avatar_url  TEXT DEFAULT NULL
);

create table characters_dg_tmp
(
  id                  INTEGER PRIMARY KEY AUTOINCREMENT,
  created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
  name                VARCHAR(100)                        NOT NULL,
  favorite            BIT(1)    DEFAULT 0                 NOT NULL,
  avatar_url          TEXT      DEFAULT NULL,
  appearance          TEXT      DEFAULT NULL,
  personality         TEXT      DEFAULT NULL,
  history             TEXT      DEFAULT NULL,
  group_talkativeness FLOAT     DEFAULT 0.5               NOT NULL,
  age                 INTEGER,
  pronouns            TEXT,
  species_id          INT                                 REFERENCES species (id) ON DELETE SET NULL
);

insert into characters_dg_tmp(id, created_at, name, favorite, avatar_url, appearance, personality, history,
                              group_talkativeness, age, pronouns, species_id)
select id,
       created_at,
       name,
       favorite,
       avatar_url,
       appearance,
       personality,
       history,
       group_talkativeness,
       age,
       pronouns,
       null
from characters;

drop table characters;

alter table characters_dg_tmp
  rename to characters;

