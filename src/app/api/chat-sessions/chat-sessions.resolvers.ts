import {ResolveFn} from '@angular/router';
import {ChatSessions} from '@api/chat-sessions/chat-sessions.service';
import {ChatMessage, ChatParticipant, ChatSession} from '@api/chat-sessions/chat-sessions.model';
import {inject} from '@angular/core';
import {paramAsId, resolveNewOrExisting} from '@util/ng';

export function chatSessionsResolverFactory(worldIdParam: string): ResolveFn<ChatSession[]> {
  return route => {
    const service = inject(ChatSessions)
    return resolveNewOrExisting(
      route, worldIdParam,
      () => [],
      worldId => service.getAll(worldId)
    )
  }
}

export function chatSessionResolverFactory(worldIdParam: string, sessionIdParam: string): ResolveFn<ChatSession> {
  return route => {
    const service = inject(ChatSessions)
    const worldId = paramAsId(route, worldIdParam)
    const sessionId = paramAsId(route, sessionIdParam)
    return service.get(worldId, sessionId)
  }
}

export function chatParticipantsResolverFactory(worldIdParam: string, sessionIdParam: string): ResolveFn<ChatParticipant[]> {
  return route => {
    const service = inject(ChatSessions)
    const worldId = paramAsId(route, worldIdParam)
    const sessionId = paramAsId(route, sessionIdParam)
    return service.getParticipants(worldId, sessionId)
  }
}

export function chatMessagesResolverFactory(worldIdParam: string, sessionIdParam: string): ResolveFn<ChatMessage[]> {
  return route => {
    const service = inject(ChatSessions)
    const worldId = paramAsId(route, worldIdParam)
    const sessionId = paramAsId(route, sessionIdParam)
    return service.getMessages(worldId, sessionId)
  }
}
