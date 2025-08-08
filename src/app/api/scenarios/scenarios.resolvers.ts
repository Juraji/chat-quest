import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {resolveNewOrExisting} from '@util/resolvers';
import {Scenarios} from './scenarios.service';
import {Scenario} from './scenarios.model';
import {NEW_ID} from '@api/common';

export const scenariosResolver: ResolveFn<Scenario[]> = () => {
  const service = inject(Scenarios)
  return service.getAll();
};

export function scenarioResolverFactory(idParam: string): ResolveFn<Scenario> {
  return route => {
    const service = inject(Scenarios)
    return resolveNewOrExisting(
      route,
      idParam,
      () => ({
        id: NEW_ID,
        name: '',
        description: '',
        avatarUrl: null,
        linkedCharacterId: null
      }),
      id => service.get(id)
    )
  }
}
