import {ChatQuestModel} from '@api/common';

export interface Tag extends ChatQuestModel {
  label: string
  readonly lowercase: string
}
