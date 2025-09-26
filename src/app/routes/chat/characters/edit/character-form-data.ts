import {Character} from '@api/characters';

export interface CharacterFormData {
  character: Character
  dialogueExamples: string[]
  greetings: string[]
}
