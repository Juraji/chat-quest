create table characters_dg_tmp
(
  id                  INTEGER
    primary key autoincrement,
  created_at          TIMESTAMP default CURRENT_TIMESTAMP not null,
  name                VARCHAR(100)                        not null,
  favorite            BIT(1)    default 0                 not null,
  avatar_url          TEXT      default NULL,
  appearance          TEXT      default NULL,
  personality         TEXT      default NULL,
  history             TEXT      default NULL,
  group_talkativeness FLOAT     default 0.5               not null
);

insert into characters_dg_tmp(id, created_at, name, favorite, avatar_url, appearance, personality, history,
                              group_talkativeness)
select id,
       created_at,
       name,
       favorite,
       avatar_url,
       appearance,
       personality,
       history,
       group_talkativeness
from characters;

drop table characters;

alter table characters_dg_tmp
  rename to characters;

