import {ChatQuestModel} from '@api/common';

export type InstructionType = 'CHAT' | 'MEMORIES' | 'TITLE_GENERATION' | 'CHARACTER_EXPORT';

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
  allowMultiCharacterResponses: boolean,
  enableReasoningParsing: boolean,
  reasoningPrefix: string,
  reasoningSuffix: string,
  enableCharacterMarkers: boolean,
  characterIdPrefix: string,
  characterIdSuffix: string,

  // Prompt Templates
  systemPrompt: Nullable<string>
  worldSetup: Nullable<string>
  instruction: string
}
