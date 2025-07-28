import {Component, inject, Signal} from '@angular/core';
import {ActivatedRoute, RouterLink} from '@angular/router';
import {routeDataSignal} from '@util/ng';
import {Character} from '@api/model';
import {PageHeader} from '@components/page-header';

@Component({
  selector: 'app-manage-characters',
  imports: [
    PageHeader,
    RouterLink
  ],
  templateUrl: './manage-characters-page.html',
  styleUrls: ['./manage-characters-page.scss'],
})
export class ManageCharactersPage {
  private readonly activatedRoute = inject(ActivatedRoute);

  readonly characters: Signal<Character[]> = routeDataSignal(this.activatedRoute, 'characters');
}
