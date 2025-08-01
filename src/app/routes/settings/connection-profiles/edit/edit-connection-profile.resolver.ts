import {ResolveFn} from '@angular/router';
import {ConnectionProfile, LlmModel, NEW_ID} from '@api/model';
import {inject} from '@angular/core';
import {ConnectionProfiles} from '@api/clients';
import {resolveNewOrExisting} from '@util/resolvers';

export const editConnectionProfileResolver: ResolveFn<ConnectionProfile> = route => {
  const service = inject(ConnectionProfiles)
  const profileId = route.params['profileId']

  return resolveNewOrExisting(
    profileId,
    () => ({
      id: NEW_ID,
      name: '',
      providerType: "OPEN_AI",
      baseUrl: '',
      apiKey: ''
    }),
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

export const editConnectionProfileTemplatesResolver: ResolveFn<ConnectionProfile[]> = () => {
  const service = inject(ConnectionProfiles)
  return service.getTemplates()
}
