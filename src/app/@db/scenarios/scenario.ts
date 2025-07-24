import {NewRecord, StoreRecord} from '@db/core';

export interface Scenario extends StoreRecord {
  name: string
  sceneDescription: string
}

export const NEW_SCENARIO: NewRecord<Scenario> = {
  name: '',
  sceneDescription: ''
}
