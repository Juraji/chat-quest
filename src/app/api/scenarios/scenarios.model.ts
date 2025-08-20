import {ChatQuestModel} from '@api/common';
import {SseEvent} from '@api/sse';

export interface Scenario extends ChatQuestModel {
  name: string
  description: string
  avatarUrl: Nullable<string>
  linkedCharacterId: Nullable<number>
}

export const ScenarioCreated: SseEvent<Scenario> = 'ScenarioCreated'
export const ScenarioUpdated: SseEvent<Scenario> = 'ScenarioUpdated'
export const ScenarioDeleted: SseEvent<number> = 'ScenarioDeleted'
