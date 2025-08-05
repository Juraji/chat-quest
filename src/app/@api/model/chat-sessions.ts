import {ChatQuestModel} from './model';

export interface ChatSession extends ChatQuestModel {
  worldId: number
  createdAt: Nullable<number>
  name: string
  scenarioId: Nullable<number>
  enableMemories: boolean
}

export interface ChatMessage extends ChatQuestModel {
  chatSessionId: number
  createdAt: Nullable<number>
  isUser: boolean
  characterId: number
  content: string
  readonly memoryId: Nullable<number>
}
