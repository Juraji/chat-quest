import {signal, WritableSignal} from '@angular/core';

export interface BooleanSignal extends WritableSignal<boolean> {
  toggle(): void;
}

export function booleanSignal(initial: boolean): BooleanSignal {
  const base = signal(initial);
  Object.defineProperty(base, 'toggle', {
    value: () => base.update(s => !s),
    writable: false,
    configurable: false,
    enumerable: false
  });
  return base as BooleanSignal
}
