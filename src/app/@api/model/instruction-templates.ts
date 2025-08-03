import {ChatQuestModel} from './model';

export type TemplateType = 'CHAT' | 'MEMORIES' | 'SUMMARIES'

export interface InstructionTemplate extends ChatQuestModel {
  name: string
  type: TemplateType,
  temperature: Nullable<number>
  systemPrompt: string
  instruction: string
}
