import {Routes} from '@angular/router';
import {WorldChatSessions} from './chat-sessions/world-chat-sessions';
import {WorldMemories} from './memories/world-memories';

const routes: Routes = [
  {
    path: 'chat-sessions',
    component: WorldChatSessions
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
