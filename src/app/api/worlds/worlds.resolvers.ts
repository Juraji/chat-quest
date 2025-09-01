import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Worlds} from '@api/worlds/worlds.service';
import {World} from '@api/worlds/index';
import {NEW_ID} from '@api/common';
import {resolveNewOrExisting} from '@util/ng';

export const worldsResolver: ResolveFn<World[]> = () => {
  const service = inject(Worlds)
  return service.getAll()
}

export function worldResolverFactory(idParam: string): ResolveFn<World> {
  return route => {
    const service = inject(Worlds)
    return resolveNewOrExisting(
      route,
      idParam,
      () => ({
        id: NEW_ID,
        name: '',
        description: null,
        avatarUrl: null,
        personaId: null,
      }),
      id => service.get(id)
    )
  }
}
