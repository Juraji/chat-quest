import {ChatQuestModel} from '@api/common';

export interface Memory extends ChatQuestModel {
  worldId: number
  chatSessionId: number
  characterId: number
  createdAt: Nullable<number>
  content: string
}

export interface MemoryPreferences {
  memoriesModelId: Nullable<number>
  memoriesInstructionId: Nullable<number>
  embeddingModelId: Nullable<number>
  memoryMinP: number
  memoryTriggerAfter: number
  memoryWindowSize: number
}
