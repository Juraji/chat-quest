import {ChatQuestModel} from './model';

export interface Scenario extends ChatQuestModel {
  name: string
  scene: string
}
