import {StoreRecord} from '@db/core';

export interface SystemPrompt extends StoreRecord {
  name: string
  prompt: string
}
