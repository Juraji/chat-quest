import {LOCALE_ID, Provider} from '@angular/core';
import {DATE_PIPE_DEFAULT_OPTIONS, registerLocaleData} from '@angular/common';
import localeNl from '@angular/common/locales/nl'

export function provideLocaleConfig(): Provider[] {
  registerLocaleData(localeNl);

  return [
    {
      provide: LOCALE_ID,
      useValue: navigator.language,
    },
    {
      provide: DATE_PIPE_DEFAULT_OPTIONS,
      useValue: {
        dateFormat: 'medium',
      }
    }
  ]
}
