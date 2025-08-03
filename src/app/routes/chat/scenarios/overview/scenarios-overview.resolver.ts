import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Scenarios} from '@api/clients';
import {Scenario} from '@api/model';

export const scenariosOverviewResolver: ResolveFn<Scenario[]> = () => {
  const service = inject(Scenarios)
  return service.getAll();
};
