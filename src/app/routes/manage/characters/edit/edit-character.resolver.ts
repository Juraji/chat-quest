import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Character, Characters} from '@db/characters';
import {NewRecord} from '@db/model';

const NEW_CHARACTER: NewRecord<Character> = {
  name: '',
  appearance: '',
  personality: '',
  history: '',
  likelyActions: [],
  unlikelyActions: [],
  dialogueExamples: [],
  avatar: null,
  favorite: false,
  tagIds: [],
}

export const editCharacterResolver: ResolveFn<Character | NewRecord<Character>> = (route) => {
  const service = inject(Characters)
  const characterId = route.paramMap.get('characterId')!!
  const iCharacterId = Number(characterId)

  if (characterId === 'new') {
    return {...NEW_CHARACTER}
  } else if(!isNaN(iCharacterId)) {
    return service.get(iCharacterId)
  } else{
    throw new Error(`Character with id "${characterId}" can not be loaded.`)
  }
};
