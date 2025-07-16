import {Routes} from '@angular/router';
import {SettingsPage} from './settings-page';
import {ChatSettingsPage} from './chat-settings/chat-settings-page';

const routes: Routes = [
  {
    path: '',
    component: SettingsPage,
    children: [
      {
        path: 'chat',
        component: ChatSettingsPage
      },
      {
        path: '**',
        redirectTo: 'chat'
      }
    ]
  }
]

export default routes
