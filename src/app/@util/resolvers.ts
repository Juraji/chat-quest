import {MaybeAsync, RedirectCommand} from '@angular/router';

export function resolveNewOrExisting<T>(
  id: string | number,
  onNew: () => MaybeAsync<T | RedirectCommand>,
  onExisting: (id: number) => MaybeAsync<T | RedirectCommand>,
): MaybeAsync<T | RedirectCommand> {
  if (id === "new") return onNew()

  const idN = typeof id === "number" ? id : Number(id)
  if (isNaN(idN)) throw Error(`Invalid id ${id}, expected a number`)

  return onExisting(idN)
}
