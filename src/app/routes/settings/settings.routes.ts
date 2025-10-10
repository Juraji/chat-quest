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
import {defaultInstructionTemplates, instructionResolverFactory, instructionsResolver} from '@api/instructions';
import {EditInstruction} from './edit-instruction-templates';
import {preferencesResolver} from '@api/preferences';

const routes: Routes = [
  {
    path: '',
    component: SettingsPage,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    resolve: {
      profiles: connectionProfilesResolver,
      instructions: instructionsResolver,
      instructionTemplates: defaultInstructionTemplates,
      preferences: preferencesResolver,
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
