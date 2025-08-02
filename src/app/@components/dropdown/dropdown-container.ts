import {
  booleanAttribute,
  contentChild,
  Directive,
  effect,
  ElementRef,
  HostListener,
  input,
  InputSignalWithTransform,
  linkedSignal
} from '@angular/core';
import {DropdownToggle} from '@components/dropdown/dropdown-toggle';
import {DropdownMenu} from '@components/dropdown/dropdown-menu';

@Directive({
    selector: '[dropdown]',
    host: {
        '[class.dropdown]': 'true'
    }
})
export class DropdownContainer {
    private readonly toggle = contentChild(DropdownToggle)
    private readonly toggleElement = contentChild(DropdownToggle, {read: ElementRef})
    private readonly menu = contentChild(DropdownMenu)

    readonly open: InputSignalWithTransform<boolean, unknown> = input(false, {transform: booleanAttribute});
    readonly isOpen = linkedSignal(() => this.open());

    constructor() {
        effect(() => {
            const toggle = this.toggle();
            toggle?.hostClicked?.subscribe(() => this.isOpen.update(s => !s))
        });

        effect(() => {
            const isOpen = this.isOpen();
            this.menu()?.show?.set(isOpen);
        });
    }

    @HostListener("window:click", ["$event"])
    onWindowClick(event: MouseEvent) {
        const toggleEl = this.toggleElement()?.nativeElement;
        const clickedElement = event?.target as HTMLElement;
        if (toggleEl && !toggleEl.contains(clickedElement)) {
            this.isOpen.set(false)
        }
    }

}
