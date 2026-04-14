import {Provider} from '@angular/core';

export class ChatQuestUIConfig {

  constructor(private readonly backingStorage: Storage) {
  }

  get apiBaseUrl(): string {
    return this.getItem('apiBaseUrl', 'http://localhost:8080/api')
  }

  set apiBaseUrl(value: string) {
    this.setItem('apiBaseUrl', value);
  }

  get sseMaxReconnectAttempts(): number {
    return this.getItem('sseMaxReconnectAttempts', 10);
  }

  set sseMaxReconnectAttempts(value: number) {
    this.setItem('sseMinReconnectDelayMillis', value);
  }

  get sseMinReconnectDelayMillis(): number {
    return this.getItem('sseMinReconnectDelayMillis', 100);
  }

  set sseMinReconnectDelayMillis(value: number) {
    this.setItem('sseMinReconnectDelayMillis', value);
  }

  get maxMessagesInChatView(): number {
    return this.getItem('maxMessagesInChatView', 10);
  }

  set maxMessagesInChatView(value: number) {
    this.setItem('maxMessagesInChatView', value);
  }

  private getItem<T>(key: string, defaultValue: T): T {
    const storedValue = this.backingStorage.getItem(key);
    if (storedValue != null) {
      return JSON.parse(storedValue);
    } else {
      return defaultValue;
    }
  }

  private setItem<T>(key: string, value: T): void {
    const serializedValue = JSON.stringify(value);
    this.backingStorage.setItem(key, serializedValue);
  }
}

export function provideChatQuestUIConfig(backingStorage: Storage): Provider[] {
  return [
    {
      provide: ChatQuestUIConfig,
      useValue: new ChatQuestUIConfig(backingStorage)
    }
  ]
}
