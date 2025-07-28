import {ChatQuestModel} from './model';

export interface Character extends ChatQuestModel {
  createdAt: number
  name: string
  favorite: boolean
  avatarUrl: string
}

export interface CharacterDetails {
  characterId: number
  appearance: string
  personality: string
  history: string
  scenario: string
  groupTalkativeness: number
}

export interface CharacterTextBlock {
  characterId: number
  text: string
}
