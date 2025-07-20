import {Injectable} from '@angular/core';
import {Store} from '@db/core';
import {Scenario} from '@db/scenarios/scenario';

@Injectable({
  providedIn: 'root'
})
export class Scenarios extends Store<Scenario> {
  constructor() {
    super('scenarios');
  }
}
