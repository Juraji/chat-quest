import {Component, effect, input, InputSignal, output, OutputEmitterRef, signal, WritableSignal} from '@angular/core';
import {animate, state, style, transition, trigger} from '@angular/animations';

@Component({
  selector: '[appCollapse]',
  exportAs: 'collapse',
  animations: [
    trigger('collapseAnimation', [
      state('collapsed', style({
        height: '0px',
        visibility: 'hidden'
      })),
      state('expanded', style({
        height: '*',
        visibility: 'visible'
      })),
      transition('collapsed <=> expanded', [
        animate('350ms ease-in-out')
      ])
    ])
  ],
  host: {
    '[@collapseAnimation]': 'collapsed() ? "collapsed" : "expanded"'
  },
  template: `<ng-content></ng-content>`
})
export class Collapse {
  readonly collapsed: WritableSignal<boolean> = signal(true)

  readonly collapse: InputSignal<boolean> = input(true)
  readonly collapseChange: OutputEmitterRef<boolean> = output()

  constructor() {
    effect(() => {
      this.collapsed.set(this.collapse())
    });
  }

  toggle(setCollapsed?: boolean) {
    if (typeof setCollapsed === 'boolean') {
      this.collapsed.set(setCollapsed)
    } else {
      this.collapsed.update(state => !state)
    }
    this.collapseChange.emit(this.collapsed())
  }
}
