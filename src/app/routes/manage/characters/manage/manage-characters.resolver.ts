import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Characters} from '@api/clients';
import {CharacterWithTags} from '@api/model';

export const manageCharactersResolver: ResolveFn<CharacterWithTags[]> = () => {
  const service = inject(Characters)
  return service.getAllWithTags();
};
