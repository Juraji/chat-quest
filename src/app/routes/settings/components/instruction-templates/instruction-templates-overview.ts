import {Component, input, InputSignal} from '@angular/core';
import {RouterLink} from "@angular/router";
import {InstructionTemplate} from '@api/model';
import {EmptyPipe} from '@components/empty-pipe';

@Component({
  selector: 'instruction-templates-overview',
  imports: [
    RouterLink,
    EmptyPipe
  ],
  templateUrl: './instruction-templates-overview.html'
})
export class InstructionTemplatesOverview {
  readonly templates: InputSignal<InstructionTemplate[]> = input.required()
}
