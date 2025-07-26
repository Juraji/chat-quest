import {StoreRecord} from '@db/core';

export type SystemPromptType = 'CHAT' | 'MEMORIES' | 'SUMMARIES'

export interface SystemPrompt extends StoreRecord {
  name: string
  prompt: string
  type: SystemPromptType
}
