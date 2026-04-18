// noinspection JSUnusedGlobalSymbols

import {ChatQuestModel} from "@api/common"
import {SseEvent} from '@api/sse';

export type TimeOfDay =
  "MIDNIGHT" |
  "NIGHT" |
  "EARLY_MORNING" |
  "MORNING" |
  "NOON" |
  "AFTERNOON" |
  "EVENING" |
  "LATE_NIGHT" |
  "REAL_TIME"

export interface ChatSession extends ChatQuestModel {
  worldId: number
  createdAt: Nullable<string>
  name: string
  scenarioId: Nullable<number>
  generateMemories: boolean
  useMemories: boolean
  pauseAutomaticResponses: boolean
  currentTimeOfDay: Nullable<TimeOfDay>
  chatNotes: Nullable<string>
  personaId: Nullable<number>
  chatModelId: Nullable<number>
  chatInstructionId: Nullable<number>
  readonly lastTotalTokens: number
  readonly lastCompletionTokens: number
}

export interface ChatMessage extends ChatQuestModel {
  chatSessionId: number
  createdAt: Nullable<string>
  isUser: boolean
  isGenerating: boolean
  characterId: Nullable<number>
  content: string
  reasoning: string
}

export interface ChatParticipant {
  chatSessionId: number
  characterId: number
  addedOn: string
  removedOn: Nullable<string>
  muted: boolean
  newlyAdded: boolean
}

export const ChatSessionCreated: SseEvent<ChatSession> = 'ChatSessionCreated'
export const ChatSessionUpdated: SseEvent<ChatSession> = 'ChatSessionUpdated'
export const ChatSessionDeleted: SseEvent<number> = 'ChatSessionDeleted'
export const ChatMessageCreated: SseEvent<ChatMessage> = 'ChatMessageCreated'
export const ChatMessageUpdated: SseEvent<ChatMessage> = 'ChatMessageUpdated'
export const ChatMessageDeleted: SseEvent<number> = 'ChatMessageDeleted'
export const ChatParticipantAdded: SseEvent<ChatParticipant> = 'ChatParticipantAdded'
export const ChatParticipantRemoved: SseEvent<ChatParticipant> = 'ChatParticipantRemoved'
