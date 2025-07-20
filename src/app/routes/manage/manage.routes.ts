import {Routes} from '@angular/router';
import {ManagePage} from './manage-page';
import {ManageCharactersPage} from './characters/manage/manage-characters-page';
import {manageCharactersResolver} from './characters/manage/manage-characters.resolver';
import {CharacterEditPage} from './characters/edit/character-edit-page';
import {editCharacterResolver} from './characters/edit/edit-character.resolver';
import {ManageScenariosPage} from './scenarios/manage/manage-scenarios-page';
import {manageScenariosResolver} from './scenarios/manage/manage-scenarios.resolver';
import {EditScenarioPage} from './scenarios/edit/edit-scenario-page';
import {editScenarioResolver} from './scenarios/edit/edit-scenario.resolver';

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
        runGuardsAndResolvers: "paramsOrQueryParamsChange",
        resolve: {
          character: editCharacterResolver
        }
      },
      {
        path: 'scenarios',
        component: ManageScenariosPage,
        resolve: {
          scenarios: manageScenariosResolver
        }
      },
      {
        path: 'scenarios/:scenarioId',
        component: EditScenarioPage,
        runGuardsAndResolvers: "paramsOrQueryParamsChange",
        resolve: {
          scenario: editScenarioResolver
        }
      },
      {
        path: '**',
        redirectTo: 'characters'
      }
    ]
  }
]

export default routes
