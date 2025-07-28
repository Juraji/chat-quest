import {ChatQuestModel} from './model';

export interface Tag extends ChatQuestModel {
  label: string
  readonly lowercase: string
}
