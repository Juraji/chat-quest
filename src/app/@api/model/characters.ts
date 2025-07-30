import {ChatQuestModel} from './model';

export interface Character extends ChatQuestModel {
  createdAt: Nullable<string>
  name: string
  favorite: boolean
  avatarUrl: Nullable<string>
}

export interface CharacterDetails {
  characterId: number
  appearance: Nullable<string>
  personality: Nullable<string>
  history: Nullable<string>
  groupTalkativeness: number
}
