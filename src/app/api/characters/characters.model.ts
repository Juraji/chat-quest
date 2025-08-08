import {Tag} from '@api/tags';
import {ChatQuestModel} from '@api/common';

export interface Character extends ChatQuestModel {
  createdAt: Nullable<string>
  name: string
  favorite: boolean
  avatarUrl: Nullable<string>
}

export interface CharacterWithTags extends Character {
  tags: Tag[]
}

export interface CharacterDetails {
  characterId: number
  appearance: Nullable<string>
  personality: Nullable<string>
  history: Nullable<string>
  groupTalkativeness: number
}
