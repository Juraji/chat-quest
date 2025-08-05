import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {ChatMessage, ChatSession, isNew} from '@api/model';
import {Observable} from 'rxjs';

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
}
