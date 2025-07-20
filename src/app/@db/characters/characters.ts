import {Injectable} from '@angular/core';
import {Character} from './character';
import {Store} from '@db/core';

@Injectable({
  providedIn: 'root'
})
export class Characters extends Store<Character> {
  constructor() {
    super('characters');
  }
}
