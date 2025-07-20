import {Component, inject, input, InputSignal} from '@angular/core';
import {ActivatedRoute, RouterLink, UrlTree} from '@angular/router';

@Component({
  selector: 'app-page-header',
  imports: [
    RouterLink
  ],
  templateUrl: './page-header.html',
})
export class PageHeader {
  readonly activatedRoute = inject(ActivatedRoute)
  readonly backUrl: InputSignal<readonly any[] | string | UrlTree | null | undefined> = input()
}
