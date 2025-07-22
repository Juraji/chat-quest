import {Injectable} from '@angular/core';
import {Tag} from './tag';
import {NewRecord, Store} from '@db/core';
import {firstValueFrom, Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class Tags extends Store<Tag> {
  constructor() {
    super('tags')
  }

  override save(record: NewRecord<Tag> | Tag): Observable<Tag> {
    // Only allow update to existing tags, else use resolve().
    if ('id' in record && !!record.id) {
      return super.save(record);
    } else {
      throw new Error("Unsupported operation, use resolve()")
    }
  }

  resolve(labels: string[]): Observable<Tag[]> {
    return this.withDatabase(async db => {

      const result: Tag[] = [];
      for (const label of labels) {
        const lowercase = label.toLowerCase();
        const existing: Tag | null = await db
          .transaction(this.storeName, 'readonly')
          .objectStore(this.storeName)
          .index('lowercase')
          .get(lowercase)

        if (!!existing) {
          result.push(existing);
        } else {
          const t = await firstValueFrom(super.save({label, lowercase}))
          result.push(t)
        }
      }

      return result
    })
  }
}
