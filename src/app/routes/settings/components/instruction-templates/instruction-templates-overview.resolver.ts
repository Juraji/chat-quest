import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {InstructionTemplates} from '@api/clients';
import {InstructionTemplate} from '@api/model';

export const instructionTemplatesOverviewResolver: ResolveFn<InstructionTemplate[]> = () => {
  const service = inject(InstructionTemplates)
  return service.getAll();
};
