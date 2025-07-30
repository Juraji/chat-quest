import {Component, ElementRef, HostListener, viewChild} from '@angular/core';
import {booleanSignal, BooleanSignal} from '@util/ng';

@Component({
  selector: 'app-dropdown',
  imports: [],
  templateUrl: './dropdown.html',
  styleUrl: './dropdown.scss'
})
export class Dropdown {
  readonly opened: BooleanSignal = booleanSignal(false)

  readonly toggleBtn = viewChild("dropdownToggle", {read: ElementRef});

  @HostListener("window:click", ['$event'])
  onWindowClick(event: MouseEvent): void {
    if (event.target == this.toggleBtn()?.nativeElement) return;
    this.opened.set(false)
  }
}
