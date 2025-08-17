DROP TABLE memory_preferences;
DROP TABLE memories;

ALTER TABLE chat_messages
  DROP COLUMN memory_id;
