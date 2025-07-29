import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Characters} from '@api/clients';
import {Character} from '@api/model';

export const manageCharactersResolver: ResolveFn<Character[]> = () => {
  const service = inject(Characters)
  return service.getAll();
};
