import {ResolveFn} from '@angular/router';
import {Character, Characters} from '@db/characters';
import {inject} from '@angular/core';

export const manageCharactersResolver: ResolveFn<Character[]> = () => {
  const characters = inject(Characters)
  return characters.getAll();
};
