import {Routes} from '@angular/router';
import {CharacterEditChatSettings} from './chat-settings/character-edit-chat-settings';
import {CharacterEditDescriptions} from './descriptions/character-edit-descriptions';
import {CharacterEditMemories} from './memories/character-edit-memories';
import {worldsResolver} from '@api/worlds';
import {CharacterBuilder} from './character-builder/character-builder';
import {instructionsResolver} from '@api/instructions';
import {llmModelViewsResolver} from '@api/providers';

const routes: Routes = [
  {
    path: "descriptions",
    component: CharacterEditDescriptions,
  },
  {
    path: "chat-settings",
    component: CharacterEditChatSettings,
  },
  {
    path: "memories",
    component: CharacterEditMemories,
    resolve: {
      worlds: worldsResolver
    }
  },
  {
    path: "character-builder",
    component: CharacterBuilder,
    resolve: {
      instructions: instructionsResolver,
      llmModels: llmModelViewsResolver,
      worlds: worldsResolver,
    }
  },
  {
    path: "**",
    redirectTo: "descriptions"
  }
]

export default routes;
