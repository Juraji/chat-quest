import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Worlds} from '@api/worlds/worlds.service';
import {World} from '@api/worlds/index';
import {NEW_ID} from '@api/common';
import {resolveNewOrExisting} from '@util/ng';
import {map} from 'rxjs';

export const worldsResolver: ResolveFn<World[]> = () => {
  const service = inject(Worlds)
  return service
    .getAll()
    .pipe(map(a => a.sort((a, b) => a.name.localeCompare(b.name))))
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
