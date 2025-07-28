import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Characters} from '@api/clients';
import {resolveNewOrExisting} from '@util/resolvers';
import {Character, CharacterDetails, NEW_ID, Tag} from '@api/model';
import {map} from 'rxjs';

export const editCharacterResolver: ResolveFn<Character> = (route, state) => {
  const service = inject(Characters)
  return resolveNewOrExisting(
    route.params['characterId'],
    () => ({
      id: NEW_ID,
      createdAt: null,
      name: '',
      favorite: false,
      avatarUrl: null
    }),
    id => service.get(id)
  )
};


export const editCharacterDetailsResolver: ResolveFn<CharacterDetails> = route => {
  const service = inject(Characters)
  return resolveNewOrExisting(
    route.params['characterId'],
    () => ({
      characterId: NEW_ID,
      appearance: null,
      personality: null,
      history: null,
      scenario: null,
      groupTalkativeness: 0.5
    }),
    id => service.getDetails(id)
  )
}

export const editCharacterTagsResolver: ResolveFn<Tag[]> = route => {
  const service = inject(Characters)
  return resolveNewOrExisting(
    route.params['characterId'],
    () => [],
    id => service.getTags(id)
  )
}

export const editCharacterDialogueExamplesResolver: ResolveFn<string[]> = route => {
  const service = inject(Characters)
  return resolveNewOrExisting(
    route.params['characterId'],
    () => [],
    id => service
      .getDialogueExamples(id)
      .pipe(map(blocks => blocks.map(b => b.text)))
  )
}

export const editCharacterGreetingsResolver: ResolveFn<string[]> = route => {
  const service = inject(Characters)
  return resolveNewOrExisting(
    route.params['characterId'],
    () => [],
    id => service
      .getGreetings(id)
      .pipe(map(blocks => blocks.map(b => b.text)))
  )
}

export const editCharacterGroupGreetingsResolver: ResolveFn<string[]> = route => {
  const service = inject(Characters)
  return resolveNewOrExisting(
    route.params['characterId'],
    () => [],
    id => service
      .getGroupGreetings(id)
      .pipe(map(blocks => blocks.map(b => b.text)))
  )
}
