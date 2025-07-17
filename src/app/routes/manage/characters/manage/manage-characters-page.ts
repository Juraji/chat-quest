import {Component, inject, Signal} from '@angular/core';
import {routeDataSignal} from '@util/ng';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {CharacterCard} from '@components/character-card/character-card';
import {Character} from '@db/characters';

@Component({
  selector: 'app-manage-characters-page',
  imports: [
    CharacterCard,
    RouterLink
  ],
  templateUrl: './manage-characters-page.html',
  styleUrls: ['./manage-characters-page.scss']
})
export class ManageCharactersPage {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly characters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'characters');
}
