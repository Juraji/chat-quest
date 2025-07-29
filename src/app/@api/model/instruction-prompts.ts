import {ChatQuestModel} from './model';

export type InstructionPromptType = 'CHAT' | 'MEMORIES' | 'SUMMARIES'

export interface InstructionPrompt extends ChatQuestModel {
  name: string
  type: InstructionPromptType,
  temperature: Nullable<number>
  systemPrompt: string
  instruction: string
}
