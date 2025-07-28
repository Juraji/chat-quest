import {Routes} from '@angular/router';

export const routes: Routes = [
  {
    path: 'chat',
    loadChildren: () => import("./routes/chat/chat.routes")
  },
  // {
  //   path: 'manage',
  //   loadChildren: () => import("./routes/manage/manage.routes")
  // },
  {
    path: 'settings',
    loadChildren: () => import("./routes/settings/settings.routes")
  },
  {
    path: '**',
    redirectTo: '/chat',
  }
];
