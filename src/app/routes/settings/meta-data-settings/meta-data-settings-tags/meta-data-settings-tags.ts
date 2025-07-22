import {Component, inject} from '@angular/core';
import {Tag, Tags} from '@db/tags';
import {toSignal} from '@angular/core/rxjs-interop';
import {Notifications} from '@components/notifications';

@Component({
  selector: 'app-meta-data-settings-tags',
  imports: [],
  templateUrl: './meta-data-settings-tags.html',
})
export class MetaDataSettingsTags {
  private readonly tags = inject(Tags)
  private readonly notifications = inject(Notifications)

  readonly availableTags = toSignal(this.tags.getAll(true))

  onDeleteTag(tag: Tag) {
    const doDelete = confirm(`Are you sure you want to delete the tag "${tag.label}"?`)
    if (doDelete) {
      this.tags
        .delete(tag.id)
        .subscribe(()=> this.notifications.toast("Tag deleted."))
    }
  }
}
