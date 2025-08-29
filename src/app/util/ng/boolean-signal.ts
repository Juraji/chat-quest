import {linkedSignal, signal, WritableSignal} from '@angular/core';

export interface BooleanSignal extends WritableSignal<boolean> {
  toggle(): void;

  setFalse(): void;

  setTrue(): void;
}

export function booleanSignal(initial: boolean | (() => boolean)): BooleanSignal {
  if (typeof initial === 'boolean') {
    return wrapSignal(signal(initial))
  } else {
    return wrapSignal(linkedSignal(initial))
  }
}

function wrapSignal(signal: WritableSignal<boolean>): BooleanSignal {
  Object.defineProperty(signal, 'toggle', {
    value: function (this: BooleanSignal) {
      this.update(s => !s)
    },
    writable: false,
    configurable: false,
    enumerable: false
  });


  Object.defineProperty(signal, 'setFalse', {
    value: function (this: BooleanSignal) {
      this.set(false)
    },
    writable: false,
    configurable: false,
    enumerable: false
  });


  Object.defineProperty(signal, 'setTrue', {
    value: function (this: BooleanSignal) {
      this.set(true)
    },
    writable: false,
    configurable: false,
    enumerable: false
  });



  return signal as BooleanSignal
}
