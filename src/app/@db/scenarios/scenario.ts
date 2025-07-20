import {StoreRecord} from '@db/core';

export interface Scenario extends StoreRecord {
  name: string
  sceneDescription: string
}
