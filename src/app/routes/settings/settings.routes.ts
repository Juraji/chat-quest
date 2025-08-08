import {Routes} from '@angular/router';
import {SettingsPage} from './settings-page';
import {
  connectionProfileResolverFactory,
  connectionProfilesResolver,
  connectionProfileTemplatesResolver,
  llmModelsResolverFactory,
  llmModelViewsResolver
} from '@api/providers/providers.resolvers';
import {EditConnectionProfile,} from "./connection-profiles"
import {instructionResolverFactory, instructionsResolver} from '@api/instructions/instructions.resolvers';
import {EditInstruction} from './edit-instruction-templates';
import {chatSettingsResolver} from '@api/worlds/worlds.resolvers';
import {memoryPreferencesResolver} from '@api/memories/memories.resolvers';

const routes: Routes = [
  {
    path: '',
    component: SettingsPage,
    resolve: {
      profiles: connectionProfilesResolver,
      templates: instructionsResolver,
      chatPreferences: chatSettingsResolver,
      memoryPreferences: memoryPreferencesResolver,
      llmModelViews: llmModelViewsResolver
    }
  },
  {
    path: 'connections/:profileId',
    component: EditConnectionProfile,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    resolve: {
      profile: connectionProfileResolverFactory('profileId'),
      models: llmModelsResolverFactory('profileId'),
      providers: connectionProfileTemplatesResolver,
    }
  },
  {
    path: 'instruction-templates/:instructionId',
    component: EditInstruction,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    resolve: {
      template: instructionResolverFactory('instructionId'),
    }
  }
]

export default routes
