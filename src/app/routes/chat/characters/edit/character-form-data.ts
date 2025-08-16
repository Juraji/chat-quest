import {Character} from '@api/characters';
import {Tag} from '@api/tags';

export interface CharacterFormData {
  character: Character
  tags: Tag[]
  dialogueExamples: string[]
  greetings: string[]
  groupGreetings: string[]
}
