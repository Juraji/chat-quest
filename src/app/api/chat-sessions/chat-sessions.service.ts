import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {ChatMessage, ChatParticipant, ChatSession} from './chat-sessions.model';
import {isNew} from '@api/common';

@Injectable({
  providedIn: 'root'
})
export class ChatSessions {
  private http: HttpClient = inject(HttpClient)

  getAll(worldId: number): Observable<ChatSession[]> {
    return this.http.get<ChatSession[]>(`/worlds/${worldId}/chat-sessions`)
  }

  get(worldId: number, sessionId: number): Observable<ChatSession> {
    return this.http.get<ChatSession>(`/worlds/${worldId}/chat-sessions/${sessionId}`)
  }

  save(worldId: number, session: ChatSession): Observable<ChatSession> {
    if (isNew(session)) {
      return this.http.post<ChatSession>(`/worlds/${worldId}/chat-sessions`, session)
    } else {
      return this.http.put<ChatSession>(`/worlds/${worldId}/chat-sessions/${session.id}`, session)
    }
  }

  delete(worldId: number, sessionId: number): Observable<void> {
    return this.http.delete<void>(`/worlds/${worldId}/chat-sessions/${sessionId}`)
  }

  getParticipants(worldId: number, sessionId: number): Observable<ChatParticipant[]> {
    return this.http.get<ChatParticipant[]>(`/worlds/${worldId}/chat-sessions/${sessionId}/participants`);
  }

  addParticipant(worldId: number, sessionId: number, characterId: number, muted: boolean): Observable<void> {
    return this.http.post<void>(
      `/worlds/${worldId}/chat-sessions/${sessionId}/participants/${characterId}`,
      null, {params: {muted}}
    );
  }

  removeParticipant(worldId: number, sessionId: number, characterId: number): Observable<void> {
    return this.http.delete<void>(`/worlds/${worldId}/chat-sessions/${sessionId}/participants/${characterId}`);
  }

  triggerParticipantResponse(worldId: number, sessionId: number, characterId: number): Observable<void> {
    return this.http.post<void>(`/worlds/${worldId}/chat-sessions/${sessionId}/participants/${characterId}/trigger-response`, null);
  }

  getMessages(worldId: number, sessionId: number): Observable<ChatMessage[]> {
    return this.http.get<ChatMessage[]>(`/worlds/${worldId}/chat-sessions/${sessionId}/chat-messages`)
  }

  saveMessage(worldId: number, sessionId: number, message: ChatMessage): Observable<ChatMessage> {
    if (isNew(message)) {
      return this.http.post<ChatMessage>(`/worlds/${worldId}/chat-sessions/${sessionId}/chat-messages`, message)
    } else {
      return this.http.put<ChatMessage>(`/worlds/${worldId}/chat-sessions/${sessionId}/chat-messages/${message.id}`, message)
    }
  }

  deleteMessage(worldId: number, sessionId: number, messageId: number): Observable<void> {
    return this.http.delete<void>(`/worlds/${worldId}/chat-sessions/${sessionId}/chat-messages/${messageId}`)
  }

  forkChatSession(worldId: number, sessionId: number, messageId: number): Observable<ChatSession> {
    return this.http.post<ChatSession>(`/worlds/${worldId}/chat-sessions/${sessionId}/chat-messages/${messageId}/fork`, null)
  }
}
