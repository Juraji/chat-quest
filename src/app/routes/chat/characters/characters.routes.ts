import {Routes} from '@angular/router';
import {ManageCharactersPage} from './manage/manage-characters-page';
import {
  characterDialogExamplesResolverFactory,
  characterGreetingsResolverFactory,
  characterResolverFactory,
  characterTagsResolverFactory
} from '@api/characters';
import {worldsResolver} from '@api/worlds';
import {EditCharacterPage} from './edit/edit-character-page';

const routes: Routes = [
  {
    path: '',
    component: ManageCharactersPage,
    resolve: {
      worlds: worldsResolver
    }
  },
  {
    path: ':characterId',
    component: EditCharacterPage,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    loadChildren: () => import("./edit/character-edit.routes"),
    resolve: {
      character: characterResolverFactory('characterId'),
      tags: characterTagsResolverFactory('characterId'),
      dialogueExamples: characterDialogExamplesResolverFactory('characterId'),
      greetings: characterGreetingsResolverFactory('characterId'),
    }
  },
]

export default routes;
