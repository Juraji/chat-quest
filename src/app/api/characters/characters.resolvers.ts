import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {resolveNewOrExisting} from '@util/resolvers';
import {Character, CharacterDetails, CharacterWithTags} from './characters.model';
import {Characters} from './characters.service';
import {NEW_ID} from '@api/common';
import {Tag} from '@api/tags';

export const charactersResolver: ResolveFn<CharacterWithTags[]> = () => {
  const service = inject(Characters)
  return service.getAllWithTags();
}

export function characterResolverFactory(idParam: string): ResolveFn<Character> {
  return route => {
    const service = inject(Characters)
    return resolveNewOrExisting(
      route,
      idParam,
      () => ({
        id: NEW_ID,
        createdAt: null,
        name: '',
        favorite: false,
        avatarUrl: null
      }),
      id => service.get(id)
    );
  }
}

export function characterDetailsResolverFactory(idParam: string): ResolveFn<CharacterDetails> {
  return route => {
    const service = inject(Characters)
    return resolveNewOrExisting(
      route,
      idParam,
      () => ({
        characterId: NEW_ID,
        appearance: null,
        personality: null,
        history: null,
        groupTalkativeness: 0.5
      }),
      id => service.getDetails(id)
    );
  }
}

export function characterDialogExamplesResolverFactory(idParam: string): ResolveFn<string[]> {
  return route => {
    const service = inject(Characters)
    return resolveNewOrExisting(
      route,
      idParam,
      () => [],
      id => service.getDialogueExamples(id)
    );
  }
}

export function characterGreetingsResolverFactory(idParam: string): ResolveFn<string[]> {
  return route => {
    const service = inject(Characters)
    return resolveNewOrExisting(
      route,
      idParam,
      () => [],
      id => service.getGreetings(id)
    );
  }
}

export function characterGroupGreetingsResolverFactory(idParam: string): ResolveFn<string[]> {
  return route => {
    const service = inject(Characters)
    return resolveNewOrExisting(
      route,
      idParam,
      () => [],
      id => service.getGroupGreetings(id)
    );
  }
}

export function characterTagsResolverFactory(idParam: string): ResolveFn<Tag[]> {
  return route => {
    const service = inject(Characters)
    return resolveNewOrExisting(
      route,
      idParam,
      () => [],
      id => service.getTags(id)
    );
  }
}
