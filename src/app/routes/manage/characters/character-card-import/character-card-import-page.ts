import {Component, inject, signal, WritableSignal} from '@angular/core';
import {PageHeader} from '@components/page-header/page-header';
import {Notifications} from '@components/notifications';
import {JsonPipe, NgTemplateOutlet} from '@angular/common';
import {CharaCard, CharaCardParser} from '@util/chara-card';

@Component({
  selector: 'app-character-card-import-page',
  imports: [
    PageHeader,
    NgTemplateOutlet,
    JsonPipe
  ],
  templateUrl: './character-card-import-page.html'
})
export class CharacterCardImportPage {
  private readonly notifications = inject(Notifications)
  private readonly parser = inject(CharaCardParser)

  readonly card: WritableSignal<CharaCard | null> = signal(null)
  readonly accept: string = this.parser.supportedContentTypes.join(',')

  async onFileSelected(e: Event) {
    e.preventDefault();

    const input = e.target as HTMLInputElement;
    const file = input.files?.item(0)

    if (!file) return;

    const fileType = file.type;

    if (!(fileType.startsWith('image/') || fileType === 'application/json')) {
      this.notifications.toast('Please select a valid image or json file.', 'DANGER')
      return
    }

    input.value = ''

    try{
      const card = await this.parser.parse(file)
    this.card.set(card)
    } catch (error) {
      this.notifications.toast(`Error parsing card: ${error}`, 'DANGER')
    }
  }
}
