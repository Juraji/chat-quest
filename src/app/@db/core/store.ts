import {concat, defer, filter, map, merge, mergeMap, Observable, Subject} from "rxjs";
import {NewRecord, StoreRecord} from "./model";
import {IDBPDatabase} from 'idb';
import {inject} from '@angular/core';
import {DATABASE} from './provider';

export abstract class Store<T extends StoreRecord> {
  private readonly database: Promise<IDBPDatabase> = inject(DATABASE)
  private readonly addSubject: Subject<T> = new Subject();
  private readonly modifySubject: Subject<T> = new Subject();
  private readonly deleteSubject: Subject<number> = new Subject();

  protected constructor(readonly storeName: string) {
  }

  getAll(watch: boolean = false): Observable<T[]> {
    const getAll$: Observable<T[]> = defer(async () => {
      const db = await this.database
      let cursor = await db
        .transaction(this.storeName, 'readonly')
        .objectStore(this.storeName)
        .openCursor()

      const result: T[] = []

      while (cursor) {
        result.push(cursor.value)
        cursor = await cursor.continue()
      }

      return result;
    })

    if (watch) {
      const watch$ = merge(this.addSubject, this.modifySubject, this.deleteSubject)
        .pipe(mergeMap(() => this.getAll()))
      return concat(getAll$, watch$)
    } else {
      return getAll$;
    }
  }

  get(id: number): Observable<T>
  get<Watch extends boolean>(id: number, watch: Watch): Observable<Watch extends true ? (T | null) : T>
  get(id: number, watch: boolean = false): Observable<T | null> {
    const get$: Observable<T | null> = defer(async () => {
      const db = await this.database
      const item: T | undefined = await db
        .transaction(this.storeName, 'readonly')
        .objectStore(this.storeName)
        .get(id)

      return item ?? null
    })

    if (watch) {
      const watch$ = merge(
        this.modifySubject.pipe(filter(it => it.id === id)),
        this.deleteSubject.pipe(filter(it => it === id), map(() => null))
      )

      return concat(get$, watch$)
    } else {
      return get$
    }
  }

  save(record: NewRecord<T> | T): Observable<T> {
    return defer(async () => {
      const db = await this.database
      const store = db
        .transaction(this.storeName, 'readwrite')
        .objectStore(this.storeName)

      if ('id' in record && !!record.id) {
        // Existing record
        await store.put(record)
        this.modifySubject.next(record)
        return record
      } else {
        const id = await store.add(record)
        const newRecord = {...record, id} as T
        this.addSubject.next(newRecord)
        return newRecord
      }
    })
  }

  delete(id: number): Observable<void> {
    return defer(async () => {
      const db = await this.database
      await db
        .transaction(this.storeName, 'readwrite')
        .objectStore(this.storeName)
        .delete(id)

      this.deleteSubject.next(id)
    })
  }

  protected withDatabase<U>(block: (db: IDBPDatabase) => Promise<U>): Observable<U> {
    return defer(async () => {
      const db = await this.database
      return block(db)
    })
  }
}
