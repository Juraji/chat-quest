import {ActivatedRoute, ActivatedRouteSnapshot, MaybeAsync, RedirectCommand} from '@angular/router';
import {Signal} from '@angular/core';
import {toSignal} from '@angular/core/rxjs-interop';
import {map} from 'rxjs';

export function routeQueryParamSignal(
  route: ActivatedRoute,
  key: string,
): Signal<string | null>
export function routeQueryParamSignal<R>(
  route: ActivatedRoute,
  key: string,
  transform: (value: string | null) => R
): Signal<R>
export function routeQueryParamSignal<R = string>(
  route: ActivatedRoute,
  key: string,
  transform?: (value: string | null) => R,
): Signal<R | null> {
  const obs$ = route.queryParamMap.pipe(
    map(queryParamMap => queryParamMap.get(key)),
    map(data => !!transform ? transform(data) : data as R)
  )

  return toSignal(obs$, {initialValue: null});
}

export function routeDataSignal<T>(
  route: ActivatedRoute,
  key: string
): Signal<T>
export function routeDataSignal<T, R>(
  route: ActivatedRoute,
  key: string,
  transform: (value: T) => R
): Signal<R>
export function routeDataSignal<T, R = T>(
  route: ActivatedRoute,
  key: string,
  transform?: (value: T) => R
): Signal<R> {
  const obs$ = route.data.pipe(
    map(data => {
      if (key in data) {
        return data[key];
      } else {
        throw new Error(`Route data does not contain the required key: ${key}`);
      }
    }),
    map(data => !!transform ? transform(data) : data)
  )

  return toSignal(obs$);
}

export function resolveNewOrExisting<T>(
  route: ActivatedRouteSnapshot,
  idParam: string,
  onNew: () => MaybeAsync<T | RedirectCommand>,
  onExisting: (id: number) => MaybeAsync<T | RedirectCommand>,
): MaybeAsync<T | RedirectCommand> {
  const idStr = route.paramMap.get(idParam);
  if (idStr == null) throw new Error(`Parameter "${idParam}" not found in route.params`);
  if (idStr === 'new') return onNew()
  return onExisting(paramAsId(idStr))
}

export function paramAsId(paramValue: string): number
export function paramAsId(route: ActivatedRouteSnapshot, param: string): number
export function paramAsId(...args: any): number {
  const arg0 = args[0]
  const arg1 = args[1]

  const idStr = arg0 instanceof ActivatedRouteSnapshot
    ? arg0.paramMap.get(arg1)
    : arg0;

  const idNum = Number(idStr)
  if (isNaN(idNum)) throw new Error(`Invalid id [${idStr}], expected a number`)
  if (idNum <= 0) throw new Error(`Invalid id [${idStr}], id should be a positive number`)

  return idNum
}
