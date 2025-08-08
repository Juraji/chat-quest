import {Routes} from '@angular/router';
import {ChatPage} from './chat-page';
import {ManageCharactersPage} from './characters/manage/manage-characters-page';
import {EditCharacterPage} from './characters/edit/edit-character-page';
import {WorldsOverviewPage} from './worlds/overview/worlds-overview-page';
import {ScenariosOverview} from './scenarios/overview/scenarios-overview';
import {EditScenarioPage} from './scenarios/edit/edit-scenario-page';
import {
  characterDetailsResolverFactory,
  characterDialogExamplesResolverFactory,
  characterGreetingsResolverFactory,
  characterGroupGreetingsResolverFactory,
  characterResolverFactory,
  charactersResolver,
  characterTagsResolverFactory
} from '@api/characters/characters.resolvers';
import {scenarioResolverFactory, scenariosResolver} from '@api/scenarios/scenarios.resolvers';

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
          characters: charactersResolver,
        }
      },
      {
        path: 'characters/:characterId',
        component: EditCharacterPage,
        runGuardsAndResolvers: "paramsOrQueryParamsChange",
        loadChildren: () => import("./characters/edit/character-edit.routes"),
        resolve: {
          character: characterResolverFactory('characterId'),
          characterDetails: characterDetailsResolverFactory('characterId'),
          tags: characterTagsResolverFactory('characterId'),
          dialogueExamples: characterDialogExamplesResolverFactory('characterId'),
          greetings: characterGreetingsResolverFactory('characterId'),
          groupGreetings: characterGroupGreetingsResolverFactory('characterId'),
        }
      },
      {
        path: 'scenarios',
        component: ScenariosOverview,
        resolve: {
          scenarios: scenariosResolver
        }
      },
      {
        path: 'scenarios/:scenarioId',
        component: EditScenarioPage,
        runGuardsAndResolvers: "paramsOrQueryParamsChange",
        resolve: {
          scenario: scenarioResolverFactory('scenarioId'),
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
