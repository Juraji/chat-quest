import {Tag} from '@api/tags';
import {ChatQuestModel} from '@api/common';
import {SseEvent} from '@api/sse';

export interface BaseCharacter extends ChatQuestModel {
  createdAt: Nullable<string>
  name: string
  favorite: boolean
  avatarUrl: Nullable<string>
}

export interface Character extends BaseCharacter {
  appearance: Nullable<string>
  personality: Nullable<string>
  history: Nullable<string>
  groupTalkativeness: number
}

export interface CharacterListView extends BaseCharacter {
  tags: Tag[]
}

export const CharacterCreated: SseEvent<Character> = 'CharacterCreated'
export const CharacterUpdated: SseEvent<Character> = 'CharacterUpdated'
export const CharacterDeleted: SseEvent<number> = 'CharacterDeleted'
