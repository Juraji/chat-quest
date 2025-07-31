import {Routes} from '@angular/router';
import {CharacterEditChatSettings} from './chat-settings/character-edit-chat-settings';
import {CharacterEditDescriptions} from './descriptions/character-edit-descriptions';
import {CharacterEditMemories} from './memories/character-edit-memories';

const routes: Routes = [
  {
    path: "chat-settings",
    component: CharacterEditChatSettings,
  },
  {
    path: "descriptions",
    component: CharacterEditDescriptions,
  },
  {
    path: "memories",
    component: CharacterEditMemories,
  },
  {
    path: "**",
    redirectTo: "chat-settings"
  }
]

export default routes;
