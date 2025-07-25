create table tags
(
  id        integer primary key autoincrement,
  label     varchar(50) not null,
  lowercase varchar(50) not null,

  constraint uk_t__label unique (label)
);

create table characters
(
  id                  integer primary key autoincrement,
  name                varchar(100) not null,
  favorite            bit(1)       not null default 0,
  created_at          timestamp    not null default current_timestamp,

  -- base
  appearance          text                  default null,
  personality         text                  default null,
  history             text                  default null,

  -- chat
  scenario            text                  default null,
  first_message       text                  default null,
  group_talkativeness float        not null default 0.5
);

create table character_tags
(
  id           integer primary key autoincrement,
  character_id integer not null,
  tag_id       integer not null,

  constraint fk_ct__character foreign key (character_id) references characters (id) on delete cascade,
  constraint fk_ct__tag foreign key (tag_id) references tags (id)
);

create table character_likely_actions
(
  id           integer primary key autoincrement,
  character_id integer not null,
  action       text    not null,

  constraint fk_cla__character foreign key (character_id) references characters (id) on delete cascade
);

create table character_unlikely_actions
(
  id           integer primary key autoincrement,
  character_id integer not null,
  action       text    not null,

  constraint fk_cua__character foreign key (character_id) references characters (id) on delete cascade
);

create table character_dialogue_examples
(
  id           integer primary key autoincrement,
  character_id integer not null,
  example      text    not null,

  constraint fk_cde__character foreign key (character_id) references characters (id) on delete cascade
);

create table character_alternate_greetings
(
  id           integer primary key autoincrement,
  character_id integer not null,
  greeting     text    not null,

  constraint fk_cag__character foreign key (character_id) references characters (id) on delete cascade
);

create table character_group_greetings
(
  id           integer primary key autoincrement,
  character_id integer not null,
  greeting     text    not null,

  constraint fk_cgg__character foreign key (character_id) references characters (id) on delete cascade
);
