import {ChatQuestModel} from './model';

export interface Memory extends ChatQuestModel {
  worldId: number
  chatSessionId: number
  characterId: number
  createdAt: Nullable<number>
  content: string
}
