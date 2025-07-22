import {Injectable} from '@angular/core';
import {CharacterImportResult, ChatQuestCharacterExport} from './model';
import {NewRecord} from '@db/core';
import {Character} from '@db/characters';
import {dataUrlToBlob, readBlobAsDataUrl, readBlobAsJson, slicedBlobOf} from '@util/blobs';

@Injectable({
  providedIn: 'root'
})
export class CharacterImportExport {
  async importFromFile(file: File): Promise<CharacterImportResult> {
    const data: ChatQuestCharacterExport = await readBlobAsJson(file)

    if (!('spec' in data && data.spec === 'chat_quest_character_v1')) {
      throw new Error('Selected file is not a ChatQuest character')
    }

    const avatar = data.avatarData ? dataUrlToBlob(data.avatarData) : null

    const character: NewRecord<Character> = {
      ...data.character,
      tagIds: [],
      avatar
    }

    return {
      character,
      tags: data.tagStrings
    }
  }

  async exportToFile(character: Character, tags: string[]): Promise<Blob> {
    const {id, avatar, tagIds, ...characterExport} = character
    const avatarData = character.avatar ? await readBlobAsDataUrl(character.avatar) : null;

    const exportData: ChatQuestCharacterExport = {
      spec: 'chat_quest_character_v1',
      character: characterExport,
      avatarData,
      tagStrings: tags,
    }

    const strData = JSON.stringify(exportData)
    return slicedBlobOf(strData, 'application/json')
  }
}
