import {Component, inject, Signal} from '@angular/core';
import {ActivatedRoute, RouterLink} from "@angular/router";
import {EmptyPipe} from '@components/empty.pipe';
import {Instruction} from '@api/instructions';
import {routeDataSignal} from '@util/ng';

@Component({
  selector: 'instruction-overview',
  imports: [
    RouterLink,
    EmptyPipe
  ],
  templateUrl: './instructions-overview.html'
})
export class InstructionOverview {
  private readonly activatedRoute = inject(ActivatedRoute)

  readonly instructionList: Signal<Instruction[]> = routeDataSignal(this.activatedRoute, 'instructions')
}
