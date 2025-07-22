import {Component, inject} from '@angular/core';
import {CharaCardParser} from '@util/chara-card';
import {map, mergeMap} from 'rxjs';
import {NewRecord} from '@db/core';
import {Character, Characters} from '@db/characters';
import {Notifications} from '@components/notifications';
import {ActivatedRoute, Router} from '@angular/router';
import {Tags} from '@db/tags';
import {BooleanSignal, booleanSignal} from '@util/ng';
import {CharacterImportExport, CharacterImportResult} from '@util/import-export';

@Component({
  selector: 'app-character-import-button',
  imports: [],
  templateUrl: './character-import-button.html'
})
export class CharacterImportButton {
  private readonly notifications = inject(Notifications)
  private readonly router = inject(Router);
  private readonly characters = inject(Characters);
  private readonly tags = inject(Tags);
  private readonly activatedRoute = inject(ActivatedRoute);
  private readonly charaCardParser = inject(CharaCardParser)
  private readonly characterImportExport = inject(CharacterImportExport)

  readonly showDropDown: BooleanSignal = booleanSignal(false)

  readonly acceptCharaCardTypes: string = this.charaCardParser.supportedContentTypes.join(',')

  async onImportFileSelected(e: Event, type: 'CharaCard' | 'ChatQuest') {
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

    try {
      let result: CharacterImportResult

      switch (type) {
        case 'CharaCard':
          result = await this.charaCardParser.parse(file)
          break
        case 'ChatQuest':
          result = await this.characterImportExport.importFromFile(file)
          break
      }

      const {character, tags: charTags} = result
      const doImport = confirm(`Do you want to import "${character.name}" from "${file.name}"?`)
      if (doImport) {
        this.tags
          .resolve(charTags)
          .pipe(
            map(tags => tags.map(t => t.id)),
            map(tagIds => ({...character, tagIds} as NewRecord<Character>)),
            mergeMap(char => this.characters.save(char))
          )
          .subscribe(character => {
            this.notifications.toast(`Character card "${character.name}" imported!`)
            this.router.navigate([character.id], {relativeTo: this.activatedRoute})
          })
      }
    } catch (error) {
      console.error(error)
      this.notifications.toast(`Error parsing card: ${error}`, 'DANGER')
    }
  }
}
