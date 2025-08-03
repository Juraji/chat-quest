import {ResolveFn} from '@angular/router';
import {InstructionTemplate, NEW_ID} from '@api/model';
import {inject} from '@angular/core';
import {InstructionTemplates} from '@api/clients';
import {resolveNewOrExisting} from '@util/resolvers';

export const editInstructionTemplateResolver: ResolveFn<InstructionTemplate> = (route) => {
  const service = inject(InstructionTemplates)
  const profileId = route.params['templateId']

  return resolveNewOrExisting(
    profileId,
    () => ({
      id: NEW_ID,
      name: '',
      type: 'CHAT',
      temperature: null,
      systemPrompt: '',
      instruction: ''
    }),
    id => service.get(id)
  )
};
