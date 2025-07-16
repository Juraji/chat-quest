import {Routes} from '@angular/router';

export const routes: Routes = [
  {
    path: 'home',
    loadChildren: () => import("./routes/home/home.routes")
  },
  {
    path: 'settings',
    loadChildren: () => import("./routes/settings/settings.routes")
  },
  {
    path: '**',
    redirectTo: '/home',
  }
];
