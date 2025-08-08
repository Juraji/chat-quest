import {ObservableInput} from 'rxjs';

export type ForkJoinSource<T> = {
  [P in keyof T]: ObservableInput<T[P]>
}
