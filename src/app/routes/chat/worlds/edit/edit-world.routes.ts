import {Routes} from '@angular/router';
import {WorldChatSessions} from './chat-sessions/world-chat-sessions';
import {WorldMemories} from './memories/world-memories';
import {scenariosResolver} from '@api/scenarios';
import {charactersResolver} from '@api/characters';

const routes: Routes = [
  {
    path: 'chat-sessions',
    component: WorldChatSessions,
    resolve: {
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
