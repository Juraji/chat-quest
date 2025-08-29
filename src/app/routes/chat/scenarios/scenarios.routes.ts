import {Routes} from '@angular/router';
import {ScenariosOverview} from './overview/scenarios-overview';
import {scenarioResolverFactory, scenariosResolver} from '@api/scenarios';
import {EditScenarioPage} from './edit/edit-scenario-page';

const routes: Routes = [
  {
    path: '',
    component: ScenariosOverview,
    resolve: {
      scenarios: scenariosResolver
    }
  },
  {
    path: ':scenarioId',
    component: EditScenarioPage,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    resolve: {
      scenario: scenarioResolverFactory('scenarioId'),
    }
  },
]

export default routes;
