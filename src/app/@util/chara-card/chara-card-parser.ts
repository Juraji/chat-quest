import {Injectable} from '@angular/core';
import {CharaCard, CharaCardSpecVersion} from './charaCard';
import ExifReader, {ValueTag} from 'exifreader';

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

  parse(file: File): Promise<CharaCard> {
    switch (file.type) {
      case 'image/png':
      case 'image/jpeg':
      case 'image/tiff':
      case 'image/webp':
        return this.parseImage(file)
      case 'application/json':
        return this.parseJson(file)
      default:
        throw new Error(`Unsupported type "${file.type}"`)
    }
  }

  detectVersion(card: CharaCard): CharaCardSpecVersion {
    if (!('spec' in card)) {
      return 'chara_card_v1';
    } else {
      return card.spec;
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
    return new Promise((resolve, reject) => {
      const reader = new FileReader()
      reader.onerror = () => reject(reader.error)
      reader.onload = () => {
        const card: CharaCard = JSON.parse(reader.result as string)
        resolve(card)
      }
      reader.readAsText(file)
    })
  }
}
