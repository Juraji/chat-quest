import {Component, inject, Signal} from '@angular/core';
import {ActivatedRoute, Router, RouterLink} from "@angular/router";
import {EmptyPipe} from '@components/empty.pipe';
import {Instruction, Instructions} from '@api/instructions';
import {NEW_ID} from '@api/common';
import {Notifications} from '@components/notifications';
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
  private readonly instructions = inject(Instructions)
  private readonly notifications = inject(Notifications)
  private readonly router = inject(Router)

  readonly instructionList: Signal<Instruction[]> = routeDataSignal(this.activatedRoute, 'instructions')

  onDuplicateInstruction(event: Event, instruction: Instruction) {
    event.stopPropagation();

    instruction = {
      ...instruction,
      id: NEW_ID,
      name: instruction.name + ' (copy)',
    }

    this.instructions
      .save(instruction)
      .subscribe(res => {
        this.notifications.toast(`Instruction copied as "${res.name}"!`)
        this.router.navigate([], {
          queryParams: {u: Date.now()},
          replaceUrl: true,
          relativeTo: this.activatedRoute
        })
      })
  }
}
