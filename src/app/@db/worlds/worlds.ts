import {Injectable} from '@angular/core';
import {Store} from '@db/core';
import {World} from './world';

@Injectable({
  providedIn: 'root'
})
export class Worlds extends Store<World> {
  constructor() {
    super('worlds');
  }
}
