import {Character, CharacterDetails} from '@api/characters';
import {Tag} from '@api/tags';

export interface CharacterFormData {
  character: Character
  characterDetails: CharacterDetails
  tags: Tag[]
  dialogueExamples: string[]
  greetings: string[]
  groupGreetings: string[]
}
