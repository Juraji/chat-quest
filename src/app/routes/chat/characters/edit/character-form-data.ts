import {Character, Tag} from '@api/characters';

export interface CharacterFormData {
  character: Character
  tags: Tag[]
  dialogueExamples: string[]
  greetings: string[]
}
