import {Routes} from '@angular/router';
import {ChatPage} from './chat-page';

const routes: Routes = [
  {
    path: '',
    component: ChatPage,
    children: [
      {
        path: 'worlds',
        loadChildren: () => import("./worlds/worlds.routes")
      },
      {
        path: 'characters',
        loadChildren: () => import("./characters/characters.routes")
      },
      {
        path: 'scenarios',
        loadChildren: () => import("./scenarios/scenarios.routes")
      },
      {
        path: '**',
        redirectTo: 'worlds'
      }
    ]
  }
]

export default routes
