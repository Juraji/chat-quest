import {
  ApplicationConfig,
  provideAppInitializer,
  provideBrowserGlobalErrorListeners,
  provideZonelessChangeDetection
} from '@angular/core';
import {provideRouter} from '@angular/router';

import {routes} from './app.routes';
import {provideHttpClient, withFetch, withInterceptors} from '@angular/common/http';
import {provideAnimationsAsync} from '@angular/platform-browser/animations/async';
import {backendUriInterceptor, provideChatQuestConfig, sseInitializer} from '@api/config';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideZonelessChangeDetection(),
    provideRouter(routes),
    provideChatQuestConfig(null),
    provideHttpClient(
      withFetch(),
      withInterceptors([
        backendUriInterceptor
      ])
    ),
    provideAppInitializer(sseInitializer),
    provideAnimationsAsync(),
  ]
};
