import {Routes} from '@angular/router';
import {SettingsPage} from './settings-page';
import {connectionProfilesOverviewResolver} from './components/connection-profiles';
import {
  EditConnectionProfile,
  editConnectionProfileLlmModelsResolver,
  editConnectionProfileResolver,
  editConnectionProfileTemplatesResolver
} from "./connection-profiles"
import {instructionTemplatesOverviewResolver} from './components/instruction-templates';
import {EditInstructionTemplate, editInstructionTemplateResolver} from './edit-instruction-templates';
import {chatSettingsResolver} from './components/chat-settings';
import {memorySettingsResolver} from './components/memory-settings';

const routes: Routes = [
  {
    path: '',
    component: SettingsPage,
    resolve: {
      profiles: connectionProfilesOverviewResolver,
      templates: instructionTemplatesOverviewResolver,
      chatPreferences: chatSettingsResolver,
      memoryPreferences: memorySettingsResolver,
    }
  },
  {
    path: 'connections/:profileId',
    component: EditConnectionProfile,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    resolve: {
      profile: editConnectionProfileResolver,
      models: editConnectionProfileLlmModelsResolver,
      providers: editConnectionProfileTemplatesResolver,
    }
  },
  {
    path: 'instruction-templates/:templateId',
    component: EditInstructionTemplate,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    resolve: {
      template: editInstructionTemplateResolver,
    }
  }
]

export default routes
