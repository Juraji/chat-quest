import {Routes} from '@angular/router';
import {WorldsOverviewPage} from './overview/worlds-overview-page';
import {worldResolverFactory, worldsResolver} from '@api/worlds';
import {ChatSessionPage, newChatSessionGuard, validatePreferencesGuard} from './session';
import {
  chatMessagesResolverFactory,
  chatParticipantsResolverFactory,
  chatSessionResolverFactory
} from '@api/chat-sessions';
import {scenariosResolver} from '@api/scenarios';
import {llmModelViewsResolver} from '@api/providers';
import {preferencesResolver} from '@api/preferences';
import {EditWorldPage} from './edit/edit-world-page';

const routes: Routes = [
  {
    path: '',
    component: WorldsOverviewPage,
    resolve: {
      worlds: worldsResolver
    }
  },
  {
    path: ':worldId/session/:chatSessionId',
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
      scenarios: scenariosResolver,
      llmModels: llmModelViewsResolver,
      preferences: preferencesResolver,
    }
  },
  {
    path: ':worldId',
    component: EditWorldPage,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    loadChildren: () => import("./edit/edit-world.routes"),
    resolve: {
      world: worldResolverFactory('worldId')
    }
  },
]

export default routes;
