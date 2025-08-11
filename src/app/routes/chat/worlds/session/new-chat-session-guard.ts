import {CanActivateFn, RedirectCommand, Router} from '@angular/router';
import {ChatSession, ChatSessions} from '@api/chat-sessions';
import {booleanAttribute, inject} from '@angular/core';
import {NEW_ID} from '@api/common';
import {map} from 'rxjs';
import {paramAsId} from '@util/ng';

export const newChatSessionGuard: CanActivateFn = (route) => {
  const params = route.paramMap

  const sessionId = params.get('chatSessionId') || 'new';

  if (sessionId !== 'new') {
    return true;
  }

  const service = inject(ChatSessions)
  const router = inject(Router);
  const query = route.queryParamMap

  const worldId = paramAsId(route, 'worldId');
  const sessionName = query.get('sessionName') || 'New Session'
  const characterIds = query.getAll('with').map(paramAsId)
  const scenarioId = query.has('scenarioId') ? paramAsId(query.get('scenarioId')!) : null
  const enableMemories = booleanAttribute(query.get('enableMemories'))

  const newChatSession: ChatSession = {
    id: NEW_ID,
    worldId,
    createdAt: null,
    name: sessionName,
    scenarioId,
    enableMemories
  }

  return service
    .create(worldId, newChatSession, characterIds)
    .pipe(map(res => {
      const urlTree = router.createUrlTree(['chat', 'worlds', worldId, 'session', res.id])

      return new RedirectCommand(urlTree, {replaceUrl: true});
    }));
};
