import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Characters} from '@api/clients';
import {resolveNewOrExisting} from '@util/resolvers';
import {CharacterTextBlock, NEW_ID} from '@api/model';
import {forkJoin, map, Observable, OperatorFunction} from 'rxjs';
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
  const flattenTextBlocks: OperatorFunction<CharacterTextBlock[], string[]> =
    map(blocks => blocks.map(b => b.text))

  return forkJoin<ForkJoinSource<CharacterFormData>>({
    character: service.get(id),
    characterDetails: service.getDetails(id),
    tags: service.getTags(id),
    dialogueExamples: service
      .getDialogueExamples(id)
      .pipe(flattenTextBlocks),
    greetings: service
      .getGreetings(id)
      .pipe(flattenTextBlocks),
    groupGreetings: service
      .getGroupGreetings(id)
      .pipe(flattenTextBlocks),
  })
}
