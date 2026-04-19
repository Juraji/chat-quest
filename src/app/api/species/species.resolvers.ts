import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {map} from 'rxjs';
import {resolveNewOrExisting} from '@util/ng';
import {NEW_ID} from '@api/common';
import {Species} from '@api/species/species.model';
import {SpeciesS} from '@api/species/species.service';

export const speciesResolver: ResolveFn<Species[]> = () => {
  const service = inject(SpeciesS)
  return service
    .getAll()
    .pipe(map(a => a.sort((a, b) => a.name.localeCompare(b.name))));
};

export function speciesResolverFactory(idParam: string): ResolveFn<Species> {
  return route => {
    const service = inject(SpeciesS)
    return resolveNewOrExisting(
      route,
      idParam,
      () => ({
        id: NEW_ID,
        name: '',
        description: '',
        avatarUrl: null,
      }),
      id => service.get(id)
    )
  }
}
