import {Routes} from '@angular/router';
import {ManagePage} from './manage-page';
import {ManageCharactersPage} from './characters/manage-characters-page';
import {ManageScenariosPage} from './scenarios/manage-scenarios-page';
import {manageCharactersResolver} from './characters/manage-characters.resolver';

const routes: Routes = [
  {
    path: '',
    component: ManagePage,
    children: [
      {
        path: 'characters',
        component: ManageCharactersPage,
        resolve: {
          characters: manageCharactersResolver
        }
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
