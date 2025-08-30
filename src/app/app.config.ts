import {
  ApplicationConfig,
  provideAppInitializer,
  provideBrowserGlobalErrorListeners,
  provideZonelessChangeDetection
} from '@angular/core';
import {provideRouter, withRouterConfig} from '@angular/router';

import {routes} from './app.routes';
import {provideHttpClient, withFetch, withInterceptors} from '@angular/common/http';
import {provideChatQuestConfig} from '@config/config';
import {backendUriInterceptor} from '@config/backend-api-uri-interceptor';
import {sseInitializer} from '@config/sse-initializer';
import {provideLocaleConfig} from '@config/locale';

export const appConfig: ApplicationConfig = {
  providers: [
    provideBrowserGlobalErrorListeners(),
    provideZonelessChangeDetection(),
    provideRouter(routes, withRouterConfig({
      paramsInheritanceStrategy: "always",
      onSameUrlNavigation: "reload"
    })),
    provideLocaleConfig(),
    provideChatQuestConfig(null),
    provideHttpClient(
      withFetch(),
      withInterceptors([
        backendUriInterceptor
      ])
    ),
    provideAppInitializer(sseInitializer),
  ]
};
