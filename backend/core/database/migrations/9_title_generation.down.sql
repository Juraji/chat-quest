create table preferences_dg_tmp
(
  id                        INTEGER
    primary key,
  chat_model_id             INTEGER
                                    references llm_models
                                      on delete set null,
  chat_instruction_id       INTEGER
                                    references instructions
                                      on delete set null,
  max_messages_in_context   integer not null,
  embedding_model_id        INTEGER
                                    references llm_models
                                      on delete set null,
  memories_model_id         INTEGER
                                    references llm_models
                                      on delete set null,
  memories_instruction_id   INTEGER
                                    references instructions
                                      on delete set null,
  memory_min_p              FLOAT   not null,
  memory_trigger_after      INTEGER not null,
  memory_window_size        INTEGER not null,
  memory_include_chat_size  integer not null,
  memory_include_chat_notes BIT(1)
);

insert into preferences_dg_tmp(id, chat_model_id, chat_instruction_id, max_messages_in_context, embedding_model_id,
                               memories_model_id, memories_instruction_id, memory_min_p, memory_trigger_after,
                               memory_window_size, memory_include_chat_size, memory_include_chat_notes)
select id,
       chat_model_id,
       chat_instruction_id,
       max_messages_in_context,
       embedding_model_id,
       memories_model_id,
       memories_instruction_id,
       memory_min_p,
       memory_trigger_after,
       memory_window_size,
       memory_include_chat_size,
       memory_include_chat_notes
from preferences;

drop table preferences;

alter table preferences_dg_tmp
  rename to preferences;

