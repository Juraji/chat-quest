import {inject, Injectable} from '@angular/core';
import {concat, defer, filter, map, Observable, Subject} from 'rxjs';
import {IDBPDatabase} from 'idb';
import {DATABASE} from './provider';

interface SettingUpdate {
  name: string;
  value: any;
}

@Injectable({
  providedIn: 'root'
})
export class Settings {
  private readonly storeName = "settings";
  private readonly modifySubject: Subject<SettingUpdate> = new Subject();
  private readonly database: Promise<IDBPDatabase> = inject(DATABASE)

  get<T>(name: string, watch: boolean = false): Observable<T | null> {
    const get$: Observable<T | null> = defer(async () => {
      const db = await this.database
      const record: T | undefined = await db
        .transaction(this.storeName, 'readonly')
        .objectStore(this.storeName)
        .get(name)

      return record ?? null;
    })

    if (watch) {
      const watch$ = this.modifySubject.pipe(
        filter(it => it.name === name),
        map(it => it.value)
      )

      return concat(get$, watch$)
    } else {
      return get$
    }
  }

  set(name: string, value: any): Observable<void> {
    return defer(async () => {
      const db = await this.database
      const store = db
        .transaction(this.storeName, 'readwrite')
        .objectStore(this.storeName)

      const existing = await store.get(name)

      if (!!existing) {
        await store.put(value, name)
      } else {
        await store.add(value, name)
      }

      this.modifySubject.next({name, value})
    })
  }
}
