import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Instructions} from './instructions.service';
import {Instruction} from './instructions.model';
import {NEW_ID} from '@api/common';
import {resolveNewOrExisting} from "@util/ng";

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
        temperature: 1.3,
        maxTokens: 300,
        topP: 0.95,
        presencePenalty: 1.1,
        frequencyPenalty: 1.1,
        stream: true,
        stopSequences: null,
        includeReasoning: false,
        reasoningPrefix: '<think>',
        reasoningSuffix: '</think>',
        characterIdPrefix: '<character>',
        characterIdSuffix: '</character>',
        systemPrompt: '',
        worldSetup: '',
        instruction: ''
      }),
      id => service.get(id)
    )
  }
}
