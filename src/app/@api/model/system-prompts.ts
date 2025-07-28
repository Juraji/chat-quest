import {ChatQuestModel} from './model';

export type SystemPromptType = 'CHAT' | 'MEMORIES' | 'SUMMARIES'

export interface SystemPrompt extends ChatQuestModel {
  name: string
  type: SystemPromptType,
  prompt: string
}
