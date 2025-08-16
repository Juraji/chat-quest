import {BaseCharacter} from './characters.model';

export function characterSortingTransformer<T extends BaseCharacter>(characters: T[]): T[] {
  return characters.sort((a, b) => {
    if (a.favorite !== b.favorite) return a.favorite ? -1 : 1;
    return a.name.localeCompare(b.name);
  })
}
