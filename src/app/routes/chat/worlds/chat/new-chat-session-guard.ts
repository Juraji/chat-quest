import {CanActivateFn, RedirectCommand, Router} from '@angular/router';
import {ChatSession, ChatSessions} from '@api/chat-sessions';
import {inject} from '@angular/core';
import {NEW_ID} from '@api/common';
import {map} from 'rxjs';

export const newChatSessionGuard: CanActivateFn = (route) => {
  const p = route.paramMap;
  const sessionId = p.get('chatSessionId')!;

  if (sessionId !== 'new') {
    return true;
  }

  const service = inject(ChatSessions)
  const router = inject(Router);

  const q = route.queryParamMap;
  const worldId = parseInt(p.get('worldId')!);
  const sessionName = q.get('sessionName') ?? 'New Session';
  const characterIds = q.getAll('with')
  const scenarioId = q.has('scenarioId') ? parseInt(q.get('scenarioId')!) : null;
  const enableMemories = q.has('enableMemories') ? q.get('enableMemories') === 'true' : true;

  if (characterIds.length === 0) {
    console.error("No character ids set")
    return false;
  }

  const newChatSession: ChatSession = {
    id: NEW_ID,
    worldId,
    createdAt: null,
    name: sessionName,
    scenarioId,
    enableMemories
  }

  return service
    .save(worldId, newChatSession)
    .pipe(map(res => {
      const urlTree = router.createUrlTree(['chat', 'worlds', worldId, 'chat', res.id])

      return new RedirectCommand(urlTree, {replaceUrl: true});
    }));
};
