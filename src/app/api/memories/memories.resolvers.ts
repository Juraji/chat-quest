import {ResolveFn} from '@angular/router';
import {Memory} from '@api/memories/memories.model';
import {inject} from '@angular/core';
import {Memories} from '@api/memories/memories.service';
import {paramAsId, resolveNewOrExisting} from '@util/ng';

export function memoriesResolverFactory(worldIdParam: string): ResolveFn<Memory[]> {
  return route => {
    const service = inject(Memories)
    return resolveNewOrExisting(
      route, worldIdParam,
      () => [],
      worldId => service.getAll(worldId)
    )
  }
}

export function characterMemoriesResolverFactory(worldIdParam: string, characterIdParam: string): ResolveFn<Memory[]> {
  return route => {
    const service = inject(Memories)
    const worldId = paramAsId(route, worldIdParam)
    const characterId = paramAsId(route, characterIdParam)
    return service.getAllByCharacter(worldId, characterId)
  }
}
