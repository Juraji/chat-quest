import {ActivatedRoute, ActivatedRouteSnapshot, MaybeAsync, RedirectCommand} from '@angular/router';
import {Signal} from '@angular/core';
import {toSignal} from '@angular/core/rxjs-interop';
import {map} from 'rxjs';

export function routeParamSignal(route: ActivatedRoute, key: string): Signal<string> {
  const obs$ = route.params.pipe(
    map(data => {
      if (key in data) {
        return data[key];
      } else {
        const fromTree = findParamInRouteTree(route.snapshot, key);
        if (!!fromTree) {
          return fromTree;
        } else {
          throw new Error(`Route params does not contain the required key: ${key}`);
        }
      }
    })
  )

  return toSignal(obs$);
}

export function routeDataSignal<T, R = T>(route: ActivatedRoute, key: string, transform?: (value: T) => R): Signal<R> {
  const obs$ = route.data.pipe(
    map(data => {
      if (key in data) {
        return data[key];
      } else {
        const fromTree = findDataInRouteTree(route.snapshot, key);
        if (!!fromTree) {
          return fromTree;
        } else {
          throw new Error(`Route data does not contain the required key: ${key}`);
        }
      }
    }),
    map(data => {
      if (!!transform) return transform(data);
      return data;
    })
  )

  return toSignal(obs$);
}

export function resolveNewOrExisting<T>(
  route: ActivatedRouteSnapshot,
  idParam: string,
  onNew: () => MaybeAsync<T | RedirectCommand>,
  onExisting: (id: number) => MaybeAsync<T | RedirectCommand>,
): MaybeAsync<T | RedirectCommand> {
  if (!(idParam in route.params)) {
    throw new Error(`Parameter "${idParam}" not found in route.params`);
  }

  const idStr = findParamInRouteTree(route, idParam, '');
  if (idStr === 'new') return onNew()

  const idNum = paramAsNumber(idStr)
  return onExisting(idNum)
}

export function paramAsNumber(paramValue: string): number
export function paramAsNumber(route: ActivatedRouteSnapshot, param: string): number
export function paramAsNumber(...args: any): number {
  const arg0 = args[0]
  const arg1 = args[1]

  const idStr = arg0 instanceof ActivatedRouteSnapshot
    ? findParamInRouteTree(arg0, arg1)
    : arg0;

  const idNum = Number(idStr)
  if (isNaN(idNum)) {
    throw new Error(`Invalid id [${idStr}], expected a number`)
  }
  if (idNum <= 0) {
    throw new Error(`Invalid id [${idStr}], id should be a positive number`)
  }

  return idNum
}

export function findParamInRouteTree(
  routeSnapshot: ActivatedRouteSnapshot,
  paramName: string,
  fallback: string,
): string
export function findParamInRouteTree(
  routeSnapshot: ActivatedRouteSnapshot,
  paramName: string,
): string | null
export function findParamInRouteTree(
  routeSnapshot: ActivatedRouteSnapshot,
  paramName: string,
  fallback?: string,
): string | null {
  let current: ActivatedRouteSnapshot | null = routeSnapshot

  while (current != null) {
    if (paramName in current.params) {
      return current.params[paramName];
    }
    current = current.parent
  }

  return fallback || null
}

export function findQueryParamInRouteTree(
  routeSnapshot: ActivatedRouteSnapshot,
  paramName: string,
  fallback: string,
): string
export function findQueryParamInRouteTree(
  routeSnapshot: ActivatedRouteSnapshot,
  paramName: string,
): string | null
export function findQueryParamInRouteTree(
  routeSnapshot: ActivatedRouteSnapshot,
  paramName: string,
  fallback?: string,
): string | null {
  let current: ActivatedRouteSnapshot | null = routeSnapshot

  while (current != null) {
    if (paramName in current.queryParams) {
      return current.queryParams[paramName];
    }
    current = current.parent
  }

  return fallback || null
}

export function findDataInRouteTree(
  routeSnapshot: ActivatedRouteSnapshot,
  paramName: string,
  fallback: string,
): string
export function findDataInRouteTree(
  routeSnapshot: ActivatedRouteSnapshot,
  paramName: string,
): string | null
export function findDataInRouteTree(
  routeSnapshot: ActivatedRouteSnapshot,
  paramName: string,
  fallback?: string,
): string | null {
  let current: ActivatedRouteSnapshot | null = routeSnapshot

  while (current != null) {
    if (paramName in current.data) {
      return current.data[paramName];
    }
    current = current.parent
  }

  return fallback || null
}
