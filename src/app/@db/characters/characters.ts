import {Injectable} from '@angular/core';
import {Character} from './character';
import {Store} from '@db/store';

@Injectable({
  providedIn: 'root'
})
export class Characters extends Store<Character> {
  static readonly STORE_NAME = 'characters';

  constructor() {
    super(Characters.STORE_NAME);
  }
}
