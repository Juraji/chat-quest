import {CanActivateFn, RedirectCommand, Router} from '@angular/router';
import {inject} from '@angular/core';
import {Worlds} from '@api/worlds';
import {Memories} from '@api/memories';
import {forkJoin, map} from 'rxjs';

export const validatePreferencesGuard: CanActivateFn = () => {
  const router = inject(Router)

  const worlds = inject(Worlds)
  const memories = inject(Memories)

  const validations$ = {
    worlds: worlds.validatePreferences(),
    memories: memories.validatePreferences(),
  }

  return forkJoin(validations$)
    .pipe(
      map(results =>
        Object.values(results)
          .filter(result => result !== null)
          .flat()),
      map(messages => {
        if (messages.length > 0) {
          const urlTree = router.createUrlTree(['settings'], {queryParams: {validate: true}});
          return new RedirectCommand(urlTree);
        } else {
          return true
        }
      })
    )
};
