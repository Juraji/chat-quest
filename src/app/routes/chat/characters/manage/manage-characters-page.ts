import {Component, inject, Signal} from '@angular/core';
import {RouterLink} from '@angular/router';
import {PageHeader} from '@components/page-header';
import {CharacterCard} from '@components/cards/character-card';
import {NewItemCard} from '@components/cards/new-item-card';
import {Character, Characters} from '@api/characters';
import {Scalable} from '@components/scalable/scalable';

@Component({
  selector: 'app-manage-characters',
  imports: [
    PageHeader,
    RouterLink,
    CharacterCard,
    NewItemCard,
    Scalable
  ],
  templateUrl: './manage-characters-page.html',
  styleUrls: ['./manage-characters-page.scss'],
})
export class ManageCharactersPage {
  private readonly charactersService = inject(Characters)

  readonly characters: Signal<Character[]> = this.charactersService.all
}
