import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Instructions} from './instructions.service';
import {Instruction} from './instructions.model';
import {NEW_ID} from '@api/common';
import {paramAsId} from "@util/ng";

export const instructionsResolver: ResolveFn<Instruction[]> = () => {
  const service = inject(Instructions)
  return service.getAll();
};

export const defaultInstructionTemplates: ResolveFn<Record<string, string>> = () => {
  const service = inject(Instructions)
  return service.defaultTemplates()
}

export function instructionResolverFactory(idParam: string): ResolveFn<Instruction> {
  return route => {
    const service = inject(Instructions)
    const idStr = route.paramMap.get(idParam);

    if (idStr === "new") {
      const templateKey = route.queryParamMap.get("templateKey")

      if (!!templateKey) {
        return service.newOfDefaultTemplate(templateKey)
      } else {
        return {
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
        }
      }
    } else {
      return service.get(paramAsId(idStr))
    }
  }
}
