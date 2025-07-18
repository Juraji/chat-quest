import {Injectable} from '@angular/core';
import {Tag} from '@db/tags/tag';
import {Store} from '@db/store';

@Injectable({
  providedIn: 'root'
})
export class Tags extends Store<Tag> {
  static readonly STORE_NAME: string = 'tags'

  constructor() {
    super(Tags.STORE_NAME)
  }
}
