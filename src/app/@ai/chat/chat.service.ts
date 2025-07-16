import {effect, inject, Injectable, Signal, signal, WritableSignal} from '@angular/core';
import {SettingsStore} from '@db/core';
import {HttpClient, HttpHeaders} from '@angular/common/http';
import {Observable, throwError} from 'rxjs';
import {OpenAiChatCompletionRequest, OpenAiChatCompletionResponse, OpenAiModel, OpenAiSettings} from './interface';

@Injectable({providedIn: 'root'})
export class ChatService {
  private readonly SETTING_NAME = 'open-ai';

  private readonly http: HttpClient = inject(HttpClient)

  readonly settings: WritableSignal<OpenAiSettings | null> = signal(null)

  constructor() {
    const settingsStore = inject(SettingsStore)
    this.settings.set(settingsStore.get(this.SETTING_NAME));

    effect(() => {
      const upd = this.settings()
      settingsStore.set(this.SETTING_NAME, upd)
    });
  }

  getModels(): Observable<OpenAiModel[]> {
    const settings = this.settings()
    if (!settings) return throwError(() => new Error('Settings not initialized'));

    const endpoint = `${settings.baseUri}/models`

    return this.http.get<OpenAiModel[]>(endpoint)
  }

  chatCompletions(request: OpenAiChatCompletionRequest): Observable<OpenAiChatCompletionResponse> {
    const settings = this.settings()
    if (!settings) return throwError(() => new Error('Settings not initialized'));

    const headers = new HttpHeaders({
      Authorization: 'Bearer ' + settings.apiKey,
      'Content-Type': 'application/json'
    });
    const endpoint = `${settings.baseUri}/chat/completions`;

    return this.http.post<OpenAiChatCompletionResponse>(endpoint, request, {headers})
  }
}
