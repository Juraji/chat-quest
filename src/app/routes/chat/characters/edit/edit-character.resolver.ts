import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Characters} from '@api/clients';
import {resolveNewOrExisting} from '@util/resolvers';
import {NEW_ID} from '@api/model';
import {forkJoin, Observable} from 'rxjs';
import {CharacterFormData} from './character-form-data';
import {ForkJoinSource} from '@util/rx';

export const editCharacterResolver: ResolveFn<CharacterFormData> = route => {
  const service = inject(Characters)
  return resolveNewOrExisting(
    route.params['characterId'],
    newCharacter,
    id => existingCharacter(service, id)
  )
}

function newCharacter(): CharacterFormData {
  return {
    character: {
      id: NEW_ID,
      createdAt: null,
      name: '',
      favorite: false,
      avatarUrl: null
    },
    characterDetails: {
      characterId: NEW_ID,
      appearance: null,
      personality: null,
      history: null,
      groupTalkativeness: 0.5
    },
    tags: [],
    dialogueExamples: [],
    greetings: [],
    groupGreetings: [],
  }
}

function existingCharacter(service: Characters, id: number): Observable<CharacterFormData> {
  return forkJoin<ForkJoinSource<CharacterFormData>>({
    character: service.get(id),
    characterDetails: service.getDetails(id),
    tags: service.getTags(id),
    dialogueExamples: service.getDialogueExamples(id),
    greetings: service.getGreetings(id),
    groupGreetings: service.getGroupGreetings(id),
  })
}
