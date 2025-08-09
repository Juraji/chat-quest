import {ResolveFn} from '@angular/router';
import {ChatSessions} from '@api/chat-sessions/chat-sessions.service';
import {ChatMessage, ChatSession} from '@api/chat-sessions/chat-sessions.model';
import {inject} from '@angular/core';
import {paramAsNumber} from '@util/ng';
import {Character} from '@api/characters';

export function chatSessionsResolverFactory(worldIdParam: string): ResolveFn<ChatSession[]> {
  return route => {
    const service = inject(ChatSessions)
    const worldId = paramAsNumber(route, worldIdParam)
    return service.getAll(worldId)
  }
}

export function chatParticipantsResolverFactory(worldIdParam: string, sessionIdParam: string): ResolveFn<Character[]> {
  return route => {
    const service = inject(ChatSessions)
    const worldId = paramAsNumber(route, worldIdParam)
    const sessionId = paramAsNumber(route, sessionIdParam)
    return service.getParticipants(worldId, sessionId)
  }
}

export function chatMessagesResolverFactory(worldIdParam: string, sessionIdParam: string): ResolveFn<ChatMessage[]> {
  return route => {
    const service = inject(ChatSessions)
    const worldId = paramAsNumber(route, worldIdParam)
    const sessionId = paramAsNumber(route, sessionIdParam)
    return service.getMessages(worldId, sessionId)
  }
}
