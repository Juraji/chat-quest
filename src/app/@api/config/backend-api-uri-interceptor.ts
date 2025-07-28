import {HttpInterceptorFn} from '@angular/common/http';

export const backendUriInterceptor: HttpInterceptorFn = (req, next) => {
  if (!req.url.startsWith('http')) {
    const url = req.url.startsWith('/') ? req.url.substring(1) : req.url;
    req = req.clone({
      url: `http://localhost:8080/api/${url}`,
    })
  }

  return next(req);
};
