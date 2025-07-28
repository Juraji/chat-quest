import {Routes} from '@angular/router';
import {ManagePage} from './manage-page';
import {ManageCharactersPage} from './characters/manage/manage-characters-page';
import {manageCharactersResolver} from './characters/manage/manage-characters.resolver';
import {EditCharacterPage} from './characters/edit/edit-character-page';
import {
  editCharacterDetailsResolver,
  editCharacterDialogueExamplesResolver,
  editCharacterGreetingsResolver,
  editCharacterGroupGreetingsResolver,
  editCharacterResolver,
  editCharacterTagsResolver
} from './characters/edit/edit-character.resolver';

const routes: Routes = [
  {
    path: '',
    component: ManagePage,
    children: [
      {
        path: 'characters',
        component: ManageCharactersPage,
        resolve: {
          characters: manageCharactersResolver,
        }
      },
      {
        path: 'characters/:characterId',
        component: EditCharacterPage,
        runGuardsAndResolvers: "paramsOrQueryParamsChange",
        resolve: {
          character: editCharacterResolver,
          characterDetails: editCharacterDetailsResolver,
          tags: editCharacterTagsResolver,
          dialogueExamples: editCharacterDialogueExamplesResolver,
          greetings: editCharacterGreetingsResolver,
          groupGreetings: editCharacterGroupGreetingsResolver,
        }
      },
      {
        path: '**',
        redirectTo: 'characters'
      }
    ]
  }
]

export default routes
