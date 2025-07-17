import {Routes} from '@angular/router';
import {ManagePage} from './manage-page';
import {ManageCharactersPage} from './characters/manage/manage-characters-page';
import {ManageScenariosPage} from './scenarios/manage-scenarios-page';
import {manageCharactersResolver} from './characters/manage/manage-characters.resolver';
import {CharacterEditPage} from './characters/edit/character-edit-page';
import {editCharacterResolver} from './characters/edit/edit-character.resolver';

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
        path: 'characters/:characterId',
        component: CharacterEditPage,
        resolve: {
          character: editCharacterResolver
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
