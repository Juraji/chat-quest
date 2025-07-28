import {ChatQuestModel} from './model';

export interface Character extends ChatQuestModel {
  createdAt: Nullable<number>
  name: string
  favorite: boolean
  avatarUrl: Nullable<string>
}

export interface CharacterDetails {
  characterId: number
  appearance: Nullable<string>
  personality: Nullable<string>
  history: Nullable<string>
  scenario: Nullable<string>
  groupTalkativeness: number
}

export interface CharacterTextBlock {
  characterId: number
  text: string
}
