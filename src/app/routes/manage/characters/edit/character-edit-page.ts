import {Component, computed, inject, Signal} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {routeDataSignal} from '@util/ng';
import {Character} from '@db/characters';
import {ReactiveFormsModule} from '@angular/forms';
import {
  CharacterEditBaseInformation
} from './components/character-edit-base-information/character-edit-base-information';

@Component({
  selector: 'app-character-edit-page',
  imports: [
    ReactiveFormsModule,
    CharacterEditBaseInformation,
  ],
  templateUrl: './character-edit-page.html'
})
export class CharacterEditPage {
  private readonly activatedRoute = inject(ActivatedRoute)

  readonly character: Signal<Character> = routeDataSignal(this.activatedRoute, 'character')
  readonly isNew = computed(() => !this.character().id)
  readonly name: Signal<string> = computed(() => this.character().name)
}
