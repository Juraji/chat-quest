import {HttpInterceptorFn} from '@angular/common/http';
import {inject} from '@angular/core';
import {ChatQuestConfig} from './config';

export const backendUriInterceptor: HttpInterceptorFn = (req, next) => {
  const baseUrl = inject(ChatQuestConfig).apiBaseUrl

  if (!req.url.startsWith('http')) {
    const url = req.url.startsWith('/') ? req.url.substring(1) : req.url;
    req = req.clone({
      url: `${baseUrl}/${url}`,
    })
  }

  return next(req);
};
