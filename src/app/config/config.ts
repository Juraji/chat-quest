import {InjectionToken, Provider} from '@angular/core';

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
    maxReconnectAttempts: 5,
    reconnectionDelay: 500,
  }
}

export const ChatQuestConfig = new InjectionToken<ChatQuestConfig>('chat_quest_config')

export function provideChatQuestConfig(config: Nullable<ChatQuestConfig>): Provider {
  return {
    provide: ChatQuestConfig,
    useValue: config ?? DEFAULT_CONFIG
  }
}
