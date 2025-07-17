import {Component, computed, inject, Signal} from '@angular/core';
import {ActivatedRoute} from '@angular/router';
import {routeDataSignal} from '@util/ng';
import {Character} from '@db/characters';
import {ReactiveFormsModule} from '@angular/forms';

@Component({
  selector: 'app-character-edit-page',
  imports: [
    ReactiveFormsModule
  ],
  templateUrl: './character-edit-page.html'
})
export class CharacterEditPage {
  private readonly activatedRoute = inject(ActivatedRoute)

  readonly character: Signal<Character> = routeDataSignal(this.activatedRoute, 'character')
  readonly isNew = computed(() => !this.character().id)
  readonly name: Signal<string> = computed(() => this.character().name)
}
