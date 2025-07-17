import {Injectable} from '@angular/core';
import {deleteRecord, getAllRecords, getRecord, saveRecord} from '@db/actions';
import {Observable, toArray} from 'rxjs';
import {Character} from './character';
import {NewRecord} from '@db/storeRecord';

@Injectable({
  providedIn: 'root'
})
export class Characters {
  static readonly STORE_NAME = 'characters';

  getAllCharacters(): Observable<Character[]> {
    return getAllRecords<Character>(Characters.STORE_NAME).pipe(toArray())
  }

  getCharacter(id: number): Observable<Character> {
    return getRecord(Characters.STORE_NAME, id)
  }

  saveCharacter(character: Character | NewRecord<Character>): Observable<Character> {
    return saveRecord(Characters.STORE_NAME, character)
  }

  deleteCharacter(id: number): Observable<void> {
    return deleteRecord(Characters.STORE_NAME, id)
  }
}
