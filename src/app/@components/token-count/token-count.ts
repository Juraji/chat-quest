import {Component, inject, input, InputSignal, signal, WritableSignal} from '@angular/core';
import {System} from '@api/clients';
import {takeUntilDestroyed, toObservable} from '@angular/core/rxjs-interop';
import {debounce, defer, iif, mergeMap} from 'rxjs';

const DEBOUNCE_TIMES = [0, 500]

@Component({
  selector: 'app-token-count',
  imports: [],
  templateUrl: './token-count.html'
})
export class TokenCount {
  private readonly system = inject(System);

  readonly text: InputSignal<Nullable<string>> = input.required()
  protected readonly tokenCount: WritableSignal<number> = signal(0)

  constructor() {
    toObservable(this.text)
      .pipe(
        takeUntilDestroyed(),
        debounce(() => DEBOUNCE_TIMES),
        mergeMap(t => iif(
          () => !!t?.trim(),
          defer(() => this.system.countTokens(t!)),
          [0]
        ))
      )
      .subscribe(count => this.tokenCount.set(count));
  }
}
