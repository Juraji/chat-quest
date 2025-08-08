import {Component, input, InputSignal} from '@angular/core';
import {RouterLink} from "@angular/router";
import {EmptyPipe} from '@components/empty-pipe';
import {Instruction} from '@api/instructions';

@Component({
  selector: 'instruction-overview',
  imports: [
    RouterLink,
    EmptyPipe
  ],
  templateUrl: './instructions-overview.html'
})
export class InstructionOverview {
  readonly templates: InputSignal<Instruction[]> = input.required()
}
