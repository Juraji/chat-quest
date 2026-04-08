import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Scenarios} from './scenarios.service';
import {Scenario} from './scenarios.model';
import {NEW_ID} from '@api/common';
import {resolveNewOrExisting} from '@util/ng';
import {map} from 'rxjs';

export const scenariosResolver: ResolveFn<Scenario[]> = () => {
  const service = inject(Scenarios)
  return service
    .getAll()
    .pipe(map(a => a.sort((a, b) => a.name.localeCompare(b.name))));
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
