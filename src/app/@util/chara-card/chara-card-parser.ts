import {Injectable} from '@angular/core';
import {CharaCard} from './charaCard';
import ExifReader, {ValueTag} from 'exifreader';
import {NEW_CHARACTER} from '@db/characters';
import {CharaCardV1} from '@util/chara-card/charaCardV1';
import {CharaCardV2} from '@util/chara-card/charaCardV2';
import {CharaCardV3} from '@util/chara-card/charaCardV3';
import {CharacterImportResult} from '@util/import-export';
import {readBlobAsJson} from '@util/blobs';

@Injectable({
  providedIn: 'root'
})
export class CharaCardParser {
  readonly supportedContentTypes: string[] = [
    'image/png',
    'image/jpeg',
    'image/tiff',
    'image/webp',
    'application/json'
  ];

  async parse(file: File): Promise<CharacterImportResult> {
    let card: CharaCard;
    let avatarImage: File | null = null;

    switch (file.type) {
      case 'image/png':
      case 'image/jpeg':
      case 'image/tiff':
      case 'image/webp':
        card = await this.parseImage(file)
        avatarImage = file
        break
      case 'application/json':
        card = await this.parseJson(file)
        break
      default:
        throw new Error(`Unsupported type "${file.type}"`)
    }

    if (!('spec' in card)) {
      return this.characterFromV1Card(card as CharaCardV1, avatarImage);
    } else if (card.spec === 'chara_card_v2') {
      return this.characterFromV2Card(card as CharaCardV2, avatarImage);
    } else if (card.spec === 'chara_card_v3') {
      return this.characterFromV3Card(card as CharaCardV3, avatarImage);
    } else {
      throw new Error(`Unsupported CharaCard version "${card.spec_version}"`);
    }
  }

  private async parseImage(file: File): Promise<CharaCard> {
    const tags = await ExifReader.load(file)
    if ('chara' in tags) {
      const tag: ValueTag = tags['chara']
      return JSON.parse(atob(tag.value))
    } else {
      throw new Error("Image does not contain character data")
    }
  }

  private parseJson(file: File): Promise<CharaCard> {
    return readBlobAsJson(file)
  }

  private characterFromV1Card(card: CharaCardV1, avatarImage: File | null): CharacterImportResult {
    return {
      character: {
        ...NEW_CHARACTER,
        name: card.name,
        personality: card.personality,
        avatar: avatarImage,
        scenario: card.scenario,
        firstMessage: card.first_mes,
        dialogueExamples: this.parseDialogExamples(card.mes_example),
      },
      tags: [],
    };
  }

  private characterFromV2Card(card: CharaCardV2, avatarImage: File | null): CharacterImportResult {
    const data = card.data

    return {
      character: {
        ...NEW_CHARACTER,
        name: data.name,
        personality: data.personality,
        avatar: avatarImage,
        scenario: data.scenario,
        firstMessage: data.first_mes,
        alternateGreetings: data.alternate_greetings,
        dialogueExamples: this.parseDialogExamples(data.mes_example),
        groupTalkativeness: this.inferTalkativenessFromExtensions(data.extensions),
      },
      tags: data.tags,
    };
  }

  private characterFromV3Card(card: CharaCardV3, avatarImage: File | null): CharacterImportResult {
    const data = card.data

    return {
      character: {
        ...NEW_CHARACTER,
        name: data.name,
        personality: data.personality,
        avatar: avatarImage,
        scenario: data.scenario,
        firstMessage: data.first_mes,
        alternateGreetings: data.alternate_greetings,
        dialogueExamples: this.parseDialogExamples(data.mes_example),
        groupTalkativeness: this.inferTalkativenessFromExtensions(data.extensions),
        groupGreetings: data.group_only_greetings
      },
      tags: data.tags,
    };
  }

  private parseDialogExamples(message: string): string[] {
    const trimmed = message.trim();
    if (trimmed.length === 0) return []

    return trimmed
      .split('{{user}}')
      .filter(v => !!v)
      .map((v) => `{{user}}${v}`)
  }

  private inferTalkativenessFromExtensions(extensions: Record<string, any>) {
    if ('talkativeness' in extensions) {
      const parsed = Number(extensions['talkativeness'])
      if (!isNaN(parsed)) return parsed;
    }
    return NEW_CHARACTER.groupTalkativeness
  }
}
