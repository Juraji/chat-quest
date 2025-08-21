import {ChatQuestModel} from '@api/common';
import {SseEvent} from '@api/sse';

export interface World extends ChatQuestModel {
  name: string
  description: Nullable<string>
  avatarUrl: Nullable<string>
}

export interface ChatPreferences {
  chatModelId: Nullable<number>
  chatInstructionId: Nullable<number>
}

export const ChatPreferencesUpdated: SseEvent<ChatPreferences> = 'ChatPreferencesUpdated'
