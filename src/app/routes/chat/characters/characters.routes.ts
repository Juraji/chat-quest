import {Routes} from '@angular/router';
import {ManageCharactersPage} from './manage/manage-characters-page';
import {
  characterDialogExamplesResolverFactory,
  characterGreetingsResolverFactory,
  characterGroupGreetingsResolverFactory,
  characterResolverFactory,
  charactersResolver,
  characterTagsResolverFactory
} from '@api/characters';
import {worldsResolver} from '@api/worlds';
import {EditCharacterPage} from './edit/edit-character-page';

const routes: Routes = [
  {
    path: '',
    component: ManageCharactersPage,
    resolve: {
      characters: charactersResolver,
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
      groupGreetings: characterGroupGreetingsResolverFactory('characterId'),
    }
  },
]

export default routes;
