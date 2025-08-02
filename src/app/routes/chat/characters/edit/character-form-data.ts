import {Character, CharacterDetails, Tag} from "@api/model";

export interface CharacterFormData {
  character: Character
  characterDetails: CharacterDetails
  tags: Tag[]
  dialogueExamples: string[]
  greetings: string[]
  groupGreetings: string[]
}
