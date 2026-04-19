import {ChatQuestModel} from '@api/common';
import {SseEvent} from '@api/sse';

export interface Species extends ChatQuestModel {
  name: string
  description: string
  avatarUrl: Nullable<string>
}

export const SpeciesCreated: SseEvent<Species> = 'SpeciesCreated'
export const SpeciesUpdated: SseEvent<Species> = 'SpeciesUpdated'
export const SpeciesDeleted: SseEvent<number> = 'SpeciesDeleted'
