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
          templates: editConnectionProfileTemplatesResolver,
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
