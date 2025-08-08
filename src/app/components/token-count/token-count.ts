import {Component, inject, input, InputSignal, signal, WritableSignal} from '@angular/core';
import {takeUntilDestroyed, toObservable} from '@angular/core/rxjs-interop';
import {debounce, defer, iif, mergeMap, timer} from 'rxjs';
import {System} from '@api/system';
import {BooleanSignal, booleanSignal} from '@util/ng';

const DEBOUNCE_TIME = 1000

@Component({
  selector: 'app-token-count',
  imports: [],
  templateUrl: './token-count.html'
})
export class TokenCount {
  private readonly system = inject(System);

  private readonly initialEmitted: BooleanSignal = booleanSignal(false)

  readonly text: InputSignal<Nullable<string>> = input.required()
  protected readonly tokenCount: WritableSignal<number> = signal(0)

  constructor() {
    toObservable(this.text)
      .pipe(
        takeUntilDestroyed(),
        debounce(() => {
          if (!this.initialEmitted()) {
            this.initialEmitted.set(true)
            return timer(0)
          } else {
            return timer(DEBOUNCE_TIME)
          }
        }),
        mergeMap(t => iif(
          () => !!t?.trim(),
          defer(() => this.system.countTokens(t!)),
          [0]
        ))
      )
      .subscribe(count => this.tokenCount.set(count));
  }
}
