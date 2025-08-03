import {ResolveFn} from '@angular/router';
import {NEW_ID, Scenario} from '@api/model';
import {inject} from '@angular/core';
import {Scenarios} from '@api/clients';
import {resolveNewOrExisting} from '@util/resolvers';

export const editScenarioResolver: ResolveFn<Scenario> = (route) => {
  const service = inject(Scenarios)
  return resolveNewOrExisting(
    route.params['scenarioId'],
    () => ({
      id: NEW_ID,
      name: '',
      description: '',
      avatarUrl: null,
      linkedCharacterId: null
    }),
    id => service.get(id)
  )
};
