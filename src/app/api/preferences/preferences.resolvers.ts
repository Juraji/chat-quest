import {ResolveFn} from '@angular/router';
import {CQPreferences} from '@api/preferences/preferences.model';
import {inject} from '@angular/core';
import {Preferences} from '@api/preferences/preferences.service';

export const preferencesResolver: ResolveFn<CQPreferences> = () => {
  const service = inject(Preferences)
  return service.get()
}
