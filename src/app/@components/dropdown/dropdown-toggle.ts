import {Directive, HostListener, output} from '@angular/core';

@Directive({
  selector: '[dropdownToggle]',
  host: {
    '[class.dropdown-toggle]': 'true'
  }
})
export class DropdownToggle {

  readonly hostClicked = output()

  @HostListener("click")
  onHostClick() {
    this.hostClicked.emit()
  }

}
