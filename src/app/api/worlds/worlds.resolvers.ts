import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Worlds} from '@api/worlds/worlds.service';
import {ChatPreferences, World} from '@api/worlds/index';
import {resolveNewOrExisting} from '@util/resolvers';
import {NEW_ID} from '@api/common';

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
        description: ''
      }),
      id => service.get(id)
    )
  }
}

export const chatSettingsResolver: ResolveFn<ChatPreferences> = () => {
  const service = inject(Worlds)
  return service.getChatPreferences();
};
