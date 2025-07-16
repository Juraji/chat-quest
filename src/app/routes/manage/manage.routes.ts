import {Routes} from '@angular/router';
import {ManagePage} from './manage-page';
import {ManageCharactersPage} from './characters/manage-characters-page';
import {ManageScenariosPage} from './scenarios/manage-scenarios-page';

const routes: Routes = [
  {
    path: '',
    component: ManagePage,
    children: [
      {
        path: 'characters',
        component: ManageCharactersPage
      },
      {
        path: 'scenarios',
        component: ManageScenariosPage
      },
      {
        path: '**',
        redirectTo: 'characters'
      }
    ]
  }
]

export default routes
