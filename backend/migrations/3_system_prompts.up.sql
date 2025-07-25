create table system_prompts
(
  id     integer primary key autoincrement,
  name   varchar(100) not null,
  prompt text         not null
);
