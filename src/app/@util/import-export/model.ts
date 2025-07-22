import {NewRecord} from '@db/core';
import {Character} from '@db/characters';

export interface CharacterImportResult {
  character: NewRecord<Character>;
  tags: string[];
}

export type SerializableCharacter = Omit<Character, 'id' | 'avatar' | 'tagIds'>

export interface ChatQuestCharacterExport {
  spec: 'chat_quest_character_v1'
  character: SerializableCharacter;
  tagStrings: string[];
  avatarData: string | null;
}
