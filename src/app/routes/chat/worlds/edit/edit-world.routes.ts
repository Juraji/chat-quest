import {Routes} from '@angular/router';
import {WorldChatSessions} from './chat-sessions/world-chat-sessions';
import {WorldMemories} from './memories/world-memories';
import {scenariosResolver} from '@api/scenarios';
import {charactersResolver} from '@api/characters';
import {chatSessionsResolverFactory} from '@api/chat-sessions/chat-sessions.resolvers';

const routes: Routes = [
  {
    path: 'chat-sessions',
    component: WorldChatSessions,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    resolve: {
      chatSessions: chatSessionsResolverFactory('worldId'),
      scenarios: scenariosResolver,
      characters: charactersResolver
    }
  },
  {
    path: 'memories',
    component: WorldMemories
  },
  {
    path: '**',
    redirectTo: 'chat-sessions'
  }
]

export default routes
