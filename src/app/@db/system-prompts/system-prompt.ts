import {StoreRecord} from '@db/model';

export interface SystemPrompt extends StoreRecord {
  name: string
  prompt: string
}
