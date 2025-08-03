import {ChatQuestModel} from './model';

export type InstructionType = 'CHAT' | 'MEMORIES' | 'SUMMARIES'

export interface InstructionTemplate extends ChatQuestModel {
  name: string
  type: InstructionType,
  temperature: Nullable<number>
  systemPrompt: string
  instruction: string
}
