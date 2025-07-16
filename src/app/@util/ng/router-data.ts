import {ActivatedRoute} from '@angular/router';
import {Signal} from '@angular/core';
import {toSignal} from '@angular/core/rxjs-interop';
import {map} from 'rxjs';

export function createRouteDataSignal<T>(route: ActivatedRoute, key: string): Signal<T> {
  const obs$ = route.data.pipe(
    map(data => {
      if (data[key] === undefined) {
        throw new Error(`Route data does not contain the required key: ${key}`);
      }
      return data[key];
    })
  )

  return toSignal(obs$);
}
