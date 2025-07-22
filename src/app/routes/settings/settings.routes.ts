import {Routes} from '@angular/router';
import {SettingsPage} from './settings-page';
import {ChatSettingsPage} from './chat-settings/chat-settings-page';
import {MetaDataSettingsPage} from './meta-data-settings/meta-data-settings-page';

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
        path: 'meta-data',
        component: MetaDataSettingsPage
      },
      {
        path: '**',
        redirectTo: 'chat'
      }
    ]
  }
]

export default routes
