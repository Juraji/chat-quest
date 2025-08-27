import {CanActivateFn} from '@angular/router';

export const validatePreferencesGuard: CanActivateFn = () => {
  return true
  // const router = inject(Router)
  //
  // const worlds = inject(Worlds)
  // const memories = inject(Memories)
  //
  // const validations$ = {
  //   worlds: worlds.validatePreferences(),
  //   memories: memories.validatePreferences(),
  // }
  //
  // return forkJoin(validations$)
  //   .pipe(
  //     map(results =>
  //       Object.values(results)
  //         .filter(result => result !== null)
  //         .flat()),
  //     map(messages => {
  //       if (messages.length > 0) {
  //         const urlTree = router.createUrlTree(['settings'], {queryParams: {validate: true}});
  //         return new RedirectCommand(urlTree);
  //       } else {
  //         return true
  //       }
  //     })
  //   )
};
