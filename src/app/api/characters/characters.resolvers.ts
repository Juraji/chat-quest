import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Characters} from './characters.service';
import {NEW_ID} from '@api/common';
import {Tag} from './tags.model';
import {resolveNewOrExisting} from '@util/ng';
import {Character} from '@api/characters/characters.model';

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
        avatarUrl: null,
        appearance: null,
        personality: null,
        history: null,
        groupTalkativeness: 0.5
      }),
      id => service.get(id)
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
