import {inject, Injectable, Signal} from '@angular/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {
  OpenAiChatCompletionRequest,
  OpenAiChatCompletionResponse,
  OpenAIListResponse,
  OpenAiModel,
  OpenAiSettings
} from './interface';
import {Settings} from '@db/settings/settings';
import {toSignal} from '@angular/core/rxjs-interop';

export const CHAT_SETTINGS_NAME = 'open-ai'

@Injectable({providedIn: 'root'})
export class ChatService {
  private readonly settings = inject(Settings)
  private readonly http: HttpClient = inject(HttpClient)

  private readonly chatSettings: Signal<OpenAiSettings | null> = toSignal(this.settings
    .get<OpenAiSettings>(CHAT_SETTINGS_NAME, true), {initialValue: null});

  getModels(): Observable<OpenAIListResponse<OpenAiModel>> {
    const settings = this.chatSettings()
    if (!settings) return throwError(() => new Error('Settings not initialized'));

    const endpoint = `${settings.baseUri}/models`

    return this.http.get<OpenAIListResponse<OpenAiModel>>(endpoint)
  }

  chatCompletions(request: OpenAiChatCompletionRequest): Observable<OpenAiChatCompletionResponse> {
    const settings = this.chatSettings()
    if (!settings) return throwError(() => new Error('Settings not initialized'));

    const headers = new HttpHeaders({
      Authorization: 'Bearer ' + settings.apiKey,
      'Content-Type': 'application/json'
    });
    const endpoint = `${settings.baseUri}/chat/completions`;

    return this.http.post<OpenAiChatCompletionResponse>(endpoint, request, {headers})
  }
}
