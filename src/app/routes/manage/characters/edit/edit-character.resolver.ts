import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Character, Characters} from '@db/characters';
import {NewRecord} from '@db/core';

const NEW_CHARACTER: NewRecord<Character> = {
  name: 'string',
  appearance: 'string',
  personality: 'string',
  avatar: null,
  favorite: false,
  tagIds: [],

  // Extended
  history: '',
  likelyActions: [],
  unlikelyActions: [],
  dialogueExamples: [],

  // Chat Defaults
  scenario: '',
  firstMessage: '',
  alternateGreetings: [],
  groupGreetings: [],
  groupTalkativeness: 0.5
}

export const editCharacterResolver: ResolveFn<Character | NewRecord<Character>> = (route) => {
  const service = inject(Characters)
  const characterId = route.paramMap.get('characterId')!!
  const iCharacterId = Number(characterId)

  if (characterId === 'new') {
    return {...NEW_CHARACTER}
  } else if (!isNaN(iCharacterId)) {
    return service.get(iCharacterId)
  } else {
    throw new Error(`Character with id "${characterId}" can not be loaded.`)
  }
};
