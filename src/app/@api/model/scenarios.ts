import {ChatQuestModel} from './model';

export interface Scenario extends ChatQuestModel {
  name: string
  description: string
  avatarUrl: Nullable<string>
  linkedCharacterId: Nullable<number>
}
