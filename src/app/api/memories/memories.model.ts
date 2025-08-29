import {ChatQuestModel} from '@api/common';

export interface Memory extends ChatQuestModel {
  worldId: number
  characterId: Nullable<number>
  createdAt: Nullable<string>
  content: string
}
