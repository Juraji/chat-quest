// noinspection JSUnusedGlobalSymbols

import {ChatQuestModel} from "@api/common"
import {SseEvent} from '@api/sse';

export interface ChatSession extends ChatQuestModel {
  worldId: number
  createdAt: Nullable<string>
  name: string
  scenarioId: Nullable<number>
  enableMemories: boolean
}

export interface ChatMessage extends ChatQuestModel {
  chatSessionId: number
  createdAt: Nullable<string>
  isUser: boolean
  characterId: Nullable<number>
  content: string
  readonly memoryId: Nullable<number>
}

export interface ChatParticipant {
  chatSessionId: number
  characterId: number
}

export const ChatSessionCreated: SseEvent<ChatSession> = 'ChatSessionCreated'
export const ChatSessionUpdated: SseEvent<ChatSession> = 'ChatSessionUpdated'
export const ChatSessionDeleted: SseEvent<number> = 'ChatSessionDeleted'
export const ChatMessageCreated: SseEvent<ChatMessage> = 'ChatMessageCreated'
export const ChatMessageUpdated: SseEvent<ChatMessage> = 'ChatMessageUpdated'
export const ChatMessageDeleted: SseEvent<number> = 'ChatMessageDeleted'
export const ChatParticipantAdded: SseEvent<ChatParticipant> = 'ChatParticipantAdded'
export const ChatParticipantRemoved: SseEvent<ChatParticipant> = 'ChatParticipantRemoved'
