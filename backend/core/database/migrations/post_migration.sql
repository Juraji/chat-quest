-- Default Chat Instruction
UPDATE preferences
SET chat_instruction_id = (SELECT id FROM instructions WHERE type = 'CHAT' LIMIT 1)
WHERE id = 0
  AND chat_instruction_id is null;
-- Default Memory Instruction
UPDATE preferences
SET memories_instruction_id = (SELECT id FROM instructions WHERE type = 'MEMORIES' LIMIT 1)
WHERE id = 0
  AND memories_instruction_id is null;
-- Default Title Generation Instruction
UPDATE preferences
SET title_generation_instruction_id = (SELECT id FROM instructions WHERE type = 'TITLE_GENERATION' LIMIT 1)
WHERE id = 0
  AND title_generation_instruction_id is null;
