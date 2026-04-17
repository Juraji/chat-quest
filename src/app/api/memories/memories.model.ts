import {ChatQuestModel} from '@api/common';
import {SseEvent} from '@api/sse';

export interface Memory extends ChatQuestModel {
  worldId: number
  characterId: Nullable<number>
  createdAt: Nullable<string>
  content: string
  alwaysInclude: boolean
}

export interface MemoryBookmarkEvent {
  chatSessionId: number
  messageId: number
}

export const MemoryCreated: SseEvent<Memory> = 'MemoryCreated'
export const MemoryUpdated: SseEvent<Memory> = 'MemoryUpdated'
export const MemoryDeleted: SseEvent<number> = 'MemoryDeleted'
export const MemoryBookmarkUpdated: SseEvent<MemoryBookmarkEvent> = 'MemoryBookmarkUpdated'
