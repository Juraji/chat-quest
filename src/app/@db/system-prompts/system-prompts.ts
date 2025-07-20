import {Injectable} from '@angular/core';
import {Store} from '@db/core';
import {SystemPrompt} from './system-prompt';

@Injectable({
  providedIn: 'root'
})
export class SystemPrompts extends Store<SystemPrompt> {
  constructor() {
    super('system-prompts');
  }
}
