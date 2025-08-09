import {ActivatedRoute} from '@angular/router';
import {Signal} from '@angular/core';
import {toSignal} from '@angular/core/rxjs-interop';
import {map} from 'rxjs';

export function routeParamSignal(route: ActivatedRoute, paramName: string): Signal<string> {
  const obs$ = route.params.pipe(
    map(data => {
      if (paramName in data) {
        return data[paramName];
      } else {
        throw new Error(`Route params does not contain the required key: ${paramName}`);
      }
    })
  )

  return toSignal(obs$);
}

export function routeDataSignal<T>(route: ActivatedRoute, key: string): Signal<T> {
  const obs$ = route.data.pipe(
    map(data => {
      if (key in data) {
        return data[key];
      }
      throw new Error(`Route data does not contain the required key: ${key}`);
    })
  )

  return toSignal(obs$);
}
