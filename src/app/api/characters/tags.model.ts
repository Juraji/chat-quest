import {ChatQuestModel} from '@api/common';
import {SseEvent} from '@api/sse';

export interface Tag extends ChatQuestModel {
  label: string
  readonly lowercase: string
}

export const CharacterTagAdded: SseEvent<[number, number]> = 'CharacterTagAdded'
export const CharacterTagRemoved: SseEvent<[number, number]> = 'CharacterTagRemoved'
export const TagCreated: SseEvent<Tag> = 'TagCreated'
export const TagUpdated: SseEvent<Tag> = 'TagUpdated'
export const TagDeleted: SseEvent<number> = 'TagDeleted'
