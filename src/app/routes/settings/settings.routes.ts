import {Routes} from '@angular/router';
import {SettingsPage} from './settings-page';
import {ManageConnectionProfiles} from './connection-profiles/manage/manage-connection-profiles';
import {manageConnectionProfilesResolver} from './connection-profiles/manage/manage-connection-profiles.resolver';
import {EditConnectionProfile} from "./connection-profiles/edit/edit-connection-profile"
import {
  editConnectionProfileLlmModelsResolver,
  editConnectionProfileResolver,
  editConnectionProfileTemplatesResolver
} from './connection-profiles/edit/edit-connection-profile.resolver';
import {InstructionTemplatesOverviewPage} from './instruction-templates/overview/instruction-templates-overview-page';
import {
  instructionTemplatesOverviewResolver
} from './instruction-templates/overview/instruction-templates-overview.resolver';
import {EditInstructionTemplate} from './instruction-templates/edit/edit-instruction-template';
import {editInstructionTemplateResolver} from './instruction-templates/edit/edit-instruction-template.resolver';

const routes: Routes = [
  {
    path: '',
    component: SettingsPage,
    children: [
      {
        path: 'connections',
        component: ManageConnectionProfiles,
        resolve: {
          profiles: manageConnectionProfilesResolver
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
        path: 'instruction-templates',
        component: InstructionTemplatesOverviewPage,
        resolve: {
          templates: instructionTemplatesOverviewResolver
        }
      },
      {
        path: 'instruction-templates/:templateId',
        component: EditInstructionTemplate,
        runGuardsAndResolvers: "paramsOrQueryParamsChange",
        resolve: {
          template: editInstructionTemplateResolver,
        }
      },
      {
        path: '**',
        redirectTo: 'connections'
      }
    ]
  }
]

export default routes
