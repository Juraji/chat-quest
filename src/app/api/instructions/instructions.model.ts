import {ChatQuestModel} from '@api/common';

export type InstructionType = 'CHAT' | 'MEMORIES'

export interface Instruction extends ChatQuestModel {
  name: string
  type: InstructionType,
  temperature: number
  maxTokens: number
  topP: number
  presencePenalty: number
  frequencyPenalty: number
  stream: boolean
  stopSequences: Nullable<string>
  systemPrompt: string
  worldSetup: string
  instruction: string
}
