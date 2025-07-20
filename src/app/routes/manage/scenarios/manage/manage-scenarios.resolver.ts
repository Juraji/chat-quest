import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Scenario, Scenarios} from '@db/scenarios';

export const manageScenariosResolver: ResolveFn<Scenario[]> = () => {
  const service = inject(Scenarios)
  return service.getAll();
};
