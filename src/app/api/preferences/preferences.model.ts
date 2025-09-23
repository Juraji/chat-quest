import {SseEvent} from '@api/sse';

export interface CQPreferences {
  chatModelId: Nullable<number>
  chatInstructionId: Nullable<number>
  maxMessagesInContext: number
  embeddingModelId: Nullable<number>
  memoriesModelId: Nullable<number>
  memoriesInstructionId: Nullable<number>
  memoryMinP: number
  memoryTriggerAfter: number
  memoryWindowSize: number
  memoryIncludeChatSize: number
  memoryIncludeChatNotes: boolean
}

export const PreferencesUpdated: SseEvent<CQPreferences> = 'PreferencesUpdated'
