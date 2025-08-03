import {Routes} from '@angular/router';
import {ChatPage} from './chat-page';
import {ManageCharactersPage} from './characters/manage/manage-characters-page';
import {manageCharactersResolver} from './characters/manage/manage-characters.resolver';
import {EditCharacterPage} from './characters/edit/edit-character-page';
import {editCharacterResolver,} from './characters/edit/edit-character.resolver';
import {WorldsOverviewPage} from './worlds/overview/worlds-overview-page';
import {ScenariosOverview} from './scenarios/overview/scenarios-overview';
import {scenariosOverviewResolver} from './scenarios/overview/scenarios-overview.resolver';
import {EditScenarioPage} from './scenarios/edit/edit-scenario-page';
import {editScenarioResolver} from './scenarios/edit/edit-scenario.resolver';

const routes: Routes = [
  {
    path: '',
    component: ChatPage,
    children: [
      {
        path: 'worlds',
        component: WorldsOverviewPage,
      },
      {
        path: 'characters',
        component: ManageCharactersPage,
        resolve: {
          characters: manageCharactersResolver,
        }
      },
      {
        path: 'characters/:characterId',
        component: EditCharacterPage,
        runGuardsAndResolvers: "paramsOrQueryParamsChange",
        loadChildren: () => import("./characters/edit/character-edit.routes"),
        resolve: {
          characterFormData: editCharacterResolver,
        }
      },
      {
        path: 'scenarios',
        component: ScenariosOverview,
        resolve: {
          scenarios: scenariosOverviewResolver
        }
      },
      {
        path: 'scenarios/:scenarioId',
        component: EditScenarioPage,
        runGuardsAndResolvers: "paramsOrQueryParamsChange",
        resolve: {
          scenario: editScenarioResolver,
        }
      },
      {
        path: '**',
        redirectTo: 'worlds'
      }
    ]
  }
]

export default routes
