import {ChatQuestModel} from '@api/common';

export interface Memory extends ChatQuestModel {
  worldId: number
  chatSessionId: number
  characterId: number
  createdAt: Nullable<string>
  content: string
}
