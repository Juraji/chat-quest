import {ChatQuestModel} from '@api/common';

export type InstructionType = 'CHAT' | 'MEMORIES'

export interface Instruction extends ChatQuestModel {
  name: string
  type: InstructionType,

  // Model Settings
  temperature: number
  maxTokens: number
  topP: number
  presencePenalty: number
  frequencyPenalty: number
  stream: boolean
  stopSequences: Nullable<string>
  includeReasoning: boolean

  // Parsing
  reasoningPrefix: string
  reasoningSuffix: string
  characterIdPrefix: string
  characterIdSuffix: string

  // Prompt Templates
  systemPrompt: Nullable<string>
  worldSetup: Nullable<string>
  instruction: string
}
