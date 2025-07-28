import {ResolveFn} from '@angular/router';
import {ConnectionProfile, LlmModel} from '@api/model';
import {inject} from '@angular/core';
import {ConnectionProfiles} from '@api/clients';
import {resolveNewOrExisting} from '@util/resolvers';

export const editConnectionProfileResolver: ResolveFn<ConnectionProfile> = route => {
  const service = inject(ConnectionProfiles)
  const profileId = route.params['profileId']

  return resolveNewOrExisting(
    profileId,
    () => service.getDefaults('OPEN_AI'),
    id => service.get(id)
  )
}

export const editConnectionProfileLlmModelsResolver: ResolveFn<LlmModel[]> = route => {
  const service = inject(ConnectionProfiles)
  const profileId = route.params['profileId']

  return resolveNewOrExisting(
    profileId,
    () => [],
    id => service.getModels(id)
  )
}
