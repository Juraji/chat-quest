import {CanActivateFn, RedirectCommand, Router} from '@angular/router';
import {inject} from '@angular/core';
import {Preferences} from '@api/preferences';
import {map} from 'rxjs';

export const validatePreferencesGuard: CanActivateFn = () => {
  const preferences = inject(Preferences)
  const router = inject(Router)

  return preferences
    .validate()
    .pipe(map(validationErrors => {
      if (validationErrors) {
        const urlTree = router.createUrlTree(['settings'], {queryParams: {validate: true}});
        return new RedirectCommand(urlTree);
      } else {
        return true
      }
    }))
};
