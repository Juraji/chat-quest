import {Component, inject, Signal} from '@angular/core';
import {ActivatedRoute, RouterLink} from "@angular/router";
import {EmptyPipe} from '@components/empty.pipe';
import {Instruction} from '@api/instructions';
import {routeDataSignal} from '@util/ng';
import {DropdownContainer, DropdownMenu, DropdownToggle} from '@components/dropdown';
import {KeyValuePipe} from '@angular/common';

@Component({
  selector: 'instruction-overview',
  imports: [
    RouterLink,
    EmptyPipe,
    DropdownContainer,
    DropdownToggle,
    DropdownMenu,
    KeyValuePipe
  ],
  templateUrl: './instructions-overview.html'
})
export class InstructionOverview {
  private readonly activatedRoute = inject(ActivatedRoute)

  readonly instructionList: Signal<Instruction[]> = routeDataSignal(this.activatedRoute, 'instructions')
  readonly instructionTemplates: Signal<Record<string, Instruction>> =
    routeDataSignal(this.activatedRoute, 'instructionTemplates')
}
