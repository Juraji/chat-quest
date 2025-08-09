import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {AiProviders, ConnectionProfile, LlmModel, LlmModelView} from './providers.model';
import {Providers} from './providers.service';
import {NEW_ID} from '@api/common';
import {resolveNewOrExisting} from '@util/ng';

export const connectionProfilesResolver: ResolveFn<ConnectionProfile[]> = () => {
  const service = inject(Providers)
  return service.getAll()
}

export function connectionProfileResolverFactory(idParam: string): ResolveFn<ConnectionProfile> {
  return route => {
    const service = inject(Providers)
    return resolveNewOrExisting(
      route,
      idParam,
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
}

export const connectionProfileTemplatesResolver: ResolveFn<AiProviders> = () => {
  const service = inject(Providers)
  return service.getTemplates()
}

export function llmModelsResolverFactory(profileIdParam: string): ResolveFn<LlmModel[]> {
  return route => {
    const service = inject(Providers)
    return resolveNewOrExisting(
      route,
      profileIdParam,
      () => [],
      id => service.getModels(id)
    )
  }
}

export const llmModelViewsResolver: ResolveFn<LlmModelView[]> = () => {
  const service = inject(Providers)
  return service.getAllLlmModelViews()
}
