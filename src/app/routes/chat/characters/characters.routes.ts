import {Routes} from '@angular/router';
import {ManageCharactersPage} from './manage/manage-characters-page';
import {
  characterDialogExamplesResolverFactory,
  characterGreetingsResolverFactory,
  characterResolverFactory,
} from '@api/characters';
import {EditCharacterPage} from './edit/edit-character-page';

const routes: Routes = [
  {
    path: '',
    component: ManageCharactersPage,
  },
  {
    path: ':characterId',
    component: EditCharacterPage,
    runGuardsAndResolvers: "paramsOrQueryParamsChange",
    loadChildren: () => import("./edit/character-edit.routes"),
    resolve: {
      character: characterResolverFactory('characterId'),
      dialogueExamples: characterDialogExamplesResolverFactory('characterId'),
      greetings: characterGreetingsResolverFactory('characterId'),
    }
  },
]

export default routes;
