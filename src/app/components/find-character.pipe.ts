import {Pipe, PipeTransform} from '@angular/core';
import {Character} from '@api/characters';

@Pipe({name: 'findCharacter'})
export class FindCharacterPipe implements PipeTransform {

  transform(characters: Character[], characterId: Nullable<number>): Nullable<Character> {
    if (!characterId) return null
    return characters.find(c => c.id === characterId);
  }
}
