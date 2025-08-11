import {booleanAttribute, Directive, input, InputSignalWithTransform, output, OutputEmitterRef} from '@angular/core';
import {booleanSignal, BooleanSignal} from '@util/ng';

@Directive({
  selector: '[collapse]',
  exportAs: 'collapse',
  host: {
    '[class.d-none]': 'collapsed()'
  },
})
export class Collapse {
  readonly collapse: InputSignalWithTransform<boolean, unknown> = input(true, {transform: booleanAttribute})
  readonly collapsed: BooleanSignal = booleanSignal(() => this.collapse());
  readonly collapseChange: OutputEmitterRef<boolean> = output()

  toggle(setCollapsed?: boolean) {
    if (typeof setCollapsed === 'boolean') {
      this.collapsed.set(setCollapsed)
    } else {
      this.collapsed.toggle()
    }
    this.collapseChange.emit(this.collapsed())
  }
}
