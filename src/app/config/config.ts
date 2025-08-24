import {InjectionToken, Provider} from '@angular/core';
import {DATE_PIPE_DEFAULT_OPTIONS} from '@angular/common';

export interface ChatQuestConfig {
  apiBaseUrl: string,
  sse: {
    maxReconnectAttempts: number,
    reconnectionDelay: number
  }
}

const DEFAULT_CONFIG: ChatQuestConfig = {
  apiBaseUrl: 'http://localhost:8080/api',
  sse: {
    maxReconnectAttempts: 15,
    reconnectionDelay: 100,
  }
}

export const ChatQuestConfig = new InjectionToken<ChatQuestConfig>('chat_quest_config')

export function provideChatQuestConfig(config: Nullable<ChatQuestConfig>): Provider[] {
  return [
    {
      provide: ChatQuestConfig,
      useValue: config ?? DEFAULT_CONFIG
    },
    {
      provide: DATE_PIPE_DEFAULT_OPTIONS,
      useValue: {
        dateFormat: 'medium',
      }
    }
  ]
}
