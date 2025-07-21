import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Character, Characters, NEW_CHARACTER} from '@db/characters';
import {NewRecord} from '@db/core';

export const editCharacterResolver: ResolveFn<Character | NewRecord<Character>> = (route, state) => {
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
