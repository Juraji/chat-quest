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
import {worldResolverFactory, worldsResolver} from '@api/worlds';
import {EditWorldPage} from './worlds/edit/edit-world-page';
import {ChatSessionPage, newChatSessionGuard, validatePreferencesGuard} from './worlds/chat';
import {
  chatMessagesResolverFactory,
  chatParticipantsResolverFactory,
  chatSessionResolverFactory
} from '@api/chat-sessions';

const routes: Routes = [
  {
    path: '',
    component: ChatPage,
    children: [
      {
        path: 'worlds',
        component: WorldsOverviewPage,
        resolve: {
          worlds: worldsResolver
        }
      },
      {
        path: 'worlds/:worldId/chat/:chatSessionId',
        component: ChatSessionPage,
        canActivate: [
          validatePreferencesGuard,
          newChatSessionGuard
        ],
        resolve: {
          world: worldResolverFactory('worldId'),
          chatSession: chatSessionResolverFactory('worldId', 'chatSessionId'),
          participants: chatParticipantsResolverFactory('worldId', 'chatSessionId'),
          messages: chatMessagesResolverFactory('worldId', 'chatSessionId'),
          allCharacters: charactersResolver
        }
      },
      {
        path: 'worlds/:worldId',
        component: EditWorldPage,
        runGuardsAndResolvers: "paramsOrQueryParamsChange",
        loadChildren: () => import("./worlds/edit/edit-world.routes"),
        resolve: {
          world: worldResolverFactory('worldId')
        }
      },
      {
        path: 'characters',
        component: ManageCharactersPage,
        resolve: {
          characters: charactersResolver,
          worlds: worldsResolver
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
