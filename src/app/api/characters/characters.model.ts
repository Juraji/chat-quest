import {ChatQuestModel} from '@api/common';
import {SseEvent} from '@api/sse';

export interface Character extends ChatQuestModel {
  createdAt: Nullable<string>
  name: string
  favorite: boolean
  avatarUrl: Nullable<string>
  appearance: Nullable<string>
  personality: Nullable<string>
  history: Nullable<string>
  groupTalkativeness: number
}

export const CharacterCreated: SseEvent<Character> = 'CharacterCreated'
export const CharacterUpdated: SseEvent<Character> = 'CharacterUpdated'
export const CharacterDeleted: SseEvent<number> = 'CharacterDeleted'
