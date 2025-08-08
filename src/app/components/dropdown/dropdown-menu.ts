import {Directive} from '@angular/core';
import {BooleanSignal, booleanSignal} from '@util/ng';

@Directive({
  selector: '[dropdownMenu]',
  host: {
    '[class.dropdown-menu]': 'true',
    '[class.show]': 'show()',
  }
})
export class DropdownMenu {
  readonly show: BooleanSignal = booleanSignal(false)
}
