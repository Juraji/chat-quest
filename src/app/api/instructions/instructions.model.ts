import {ChatQuestModel} from '@api/common';

export type InstructionType = 'CHAT' | 'MEMORIES'

export interface Instruction extends ChatQuestModel {
  name: string
  type: InstructionType,
  temperature: Nullable<number>
  systemPrompt: string
  worldSetup: string
  instruction: string
}
