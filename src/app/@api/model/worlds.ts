import {ChatQuestModel} from './model';

export interface World extends ChatQuestModel {
  name: string
  description: string
}

export interface ChatPreferences {
  chatModelId: Nullable<number>
  chatInstructionId: Nullable<number>
}
