import {Injectable} from '@angular/core';
import {Tag} from './tag';
import {Store} from '@db/store';

@Injectable({
  providedIn: 'root'
})
export class Tags extends Store<Tag> {
  constructor() {
    super('tags')
  }
}
