import {ChatQuestModel} from '@api/common';

export interface World extends ChatQuestModel {
  name: string
  description: Nullable<string>
  avatarUrl: Nullable<string>
}
