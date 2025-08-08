import {ActivatedRouteSnapshot, MaybeAsync, RedirectCommand} from '@angular/router';

export function resolveNewOrExisting<T>(
  route: ActivatedRouteSnapshot,
  idParam: string,
  onNew: () => MaybeAsync<T | RedirectCommand>,
  onExisting: (id: number) => MaybeAsync<T | RedirectCommand>,
): MaybeAsync<T | RedirectCommand> {
  if (!(idParam in route.params)) {
    throw new Error(`Parameter "${idParam}" not found in route.params`);
  }

  const idStr = route.params[idParam];

  if (idStr === 'new') {
    return onNew()
  }

  const idNum = Number(idStr)
  if (isNaN(idNum)) {
    throw new Error(`Invalid id ${idStr}, expected a number`)
  }
  if (idNum <= 0) {
    throw new Error(`Invalid id ${idStr}, id should be a positive number`)
  }

  return onExisting(idNum)
}
