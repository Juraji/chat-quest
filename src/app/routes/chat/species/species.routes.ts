import {Routes} from '@angular/router';
import {speciesResolver, speciesResolverFactory} from '@api/species';
import {SpeciesOverview} from './overview/species-overview';
import {EditSpeciesPage} from './edit/edit-species-page';

const routes: Routes = [
  {
    path: '',
    component: SpeciesOverview,
    resolve: {
      species: speciesResolver
    }
  },
  {
    path: ':speciesId',
    component: EditSpeciesPage,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    resolve: {
      species: speciesResolverFactory('speciesId'),
    }
  }
]

export default routes;
