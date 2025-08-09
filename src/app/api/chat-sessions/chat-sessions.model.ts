import {ChatQuestModel} from "@api/common"

export interface ChatSession extends ChatQuestModel {
  worldId: number
  createdAt: Nullable<string>
  name: string
  scenarioId: Nullable<number>
  enableMemories: boolean
}

export interface ChatMessage extends ChatQuestModel {
  chatSessionId: number
  createdAt: Nullable<string>
  isUser: boolean
  characterId: number
  content: string
  readonly memoryId: Nullable<number>
}
