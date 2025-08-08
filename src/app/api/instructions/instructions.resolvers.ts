import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {resolveNewOrExisting} from '@util/resolvers';
import {Instructions} from './instructions.service';
import {Instruction} from './instructions.model';
import {NEW_ID} from '@api/common';

export const instructionsResolver: ResolveFn<Instruction[]> = () => {
  const service = inject(Instructions)
  return service.getAll();
};

export function instructionResolverFactory(idParam: string): ResolveFn<Instruction> {
  return route => {
    const service = inject(Instructions)
    return resolveNewOrExisting(
      route,
      idParam,
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
  }
}
