import {Component, inject, Signal} from '@angular/core';
import {PageHeader} from "@components/page-header";
import {ActivatedRoute, RouterLink} from "@angular/router";
import {InstructionTemplate} from '@api/model';
import {routeDataSignal} from '@util/ng';
import {EmptyPipe} from '@components/empty-pipe';

@Component({
  selector: 'app-instruction-templates-overview-page',
  imports: [
    PageHeader,
    RouterLink,
    EmptyPipe
  ],
  templateUrl: './instruction-templates-overview-page.html'
})
export class InstructionTemplatesOverviewPage {
private readonly activatedRoute = inject(ActivatedRoute)

  readonly templates: Signal<InstructionTemplate[]> = routeDataSignal(this.activatedRoute, 'templates')
}
