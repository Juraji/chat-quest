import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Memories} from './memories.service';
import {paramAsId} from '@util/ng';

export function memoryBookmarkResolverFactory(worldIdParam: string, sessionIdParam: string): ResolveFn<Nullable<number>> {
  return route => {
    const service = inject(Memories)
    const worldId = paramAsId(route, worldIdParam)
    const sessionId = paramAsId(route, sessionIdParam)
    return service.getBookmarkMessageId(worldId, sessionId)
  }
}
