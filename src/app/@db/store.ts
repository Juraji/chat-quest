import {filter, Observable, Subject, Subscriber} from "rxjs";
import {NewRecord, StoreRecord} from "@db/model";
import {DATABASE_NAME} from '@db/init';

export abstract class Store<T extends StoreRecord> {
  private readonly addSubject: Subject<T> = new Subject();
  private readonly modifySubject: Subject<T> = new Subject();
  private readonly deleteSubject: Subject<number> = new Subject();

  protected constructor(readonly storeName: string) {
  }

  getAll(watch: boolean): Observable<T[]> {
    return this.withDatabase<T[]>((db, observer) => {
      const store = db
        .transaction(this.storeName, 'readonly')
        .objectStore(this.storeName)

      const cursorRequest = store.openCursor()
      const result: T[] = []

      cursorRequest.onerror = () => observer.error(cursorRequest.error)
      cursorRequest.onsuccess = () => {
        const cursor = cursorRequest.result
        if (!!cursor) {
          result.push(cursor.value as T)
          cursor.continue()
        } else {
          observer.next(result)

          // Setup watch or complete
          if (watch) {
            const refresh: () => void = () => this
              .getAll(false)
              .subscribe(observer.next.bind(observer))

            const addSub = this.addSubject.subscribe(refresh)
            const modifySub = this.modifySubject.subscribe(refresh)
            const deleteSub = this.deleteSubject.subscribe(refresh)

            observer.add(() => {
              addSub.unsubscribe()
              modifySub.unsubscribe()
              deleteSub.unsubscribe()
            })
          } else {
            observer.complete()
          }
        }
      }
    })
  }

  get<Watch extends boolean>(id: number, watch: Watch): Observable<Watch extends true ? (T | null) : T>
  get(id: number, watch: boolean): Observable<T | null> {
    return this.withDatabase<T | null>((db, observer) => {
      const request = db
        .transaction(this.storeName, 'readonly')
        .objectStore(this.storeName)
        .get(id)

      request.onerror = () => observer.error(request.error)
      request.onsuccess = () => {
        // Emit initial data or complete if not found
        if (request.result) {
          observer.next(request.result as T)
        } else {
          observer.complete()
          return
        }

        // Setup watch or complete
        if (watch) {
          const modifySub = this.modifySubject
            .pipe(filter(it => it.id === id))
            .subscribe(record => observer.next(record))

          const deleteSub = this.deleteSubject
            .pipe(filter(it => it === id))
            .subscribe(() => observer.next(null))

          observer.add(() => {
            modifySub.unsubscribe()
            deleteSub.unsubscribe()
          })
        } else {
          observer.complete()
        }
      }
    })
  }

  save(record: NewRecord<T> | T): Observable<T> {
    return this.withDatabase<T>((db, observer) => {
      const transaction = db.transaction(this.storeName, 'readwrite')
      const store = transaction.objectStore(this.storeName)

      if ('id' in record && !!record.id) {
        // Existing record
        const request = store.put(record)
        request.onerror = () => observer.error(request.error)

        request.onsuccess = () => {
          observer.next(record)
          observer.complete()
          this.modifySubject.next(record)
        }
      } else {
        // New record
        const request = store.add(record)
        request.onerror = () => observer.error(request.error)

        request.onsuccess = () => {
          const newRecord = {...record, id: request.result} as T
          observer.next(newRecord)
          observer.complete()
          this.addSubject.next(newRecord)
        }
      }
    })
  }

  delete(id: number): Observable<void> {
    return this.withDatabase((db, observer) => {
      const transaction = db.transaction(this.storeName, 'readwrite')
      const store = transaction.objectStore(this.storeName)

      // Delete the record
      const deleteRequest = store.delete(id)
      deleteRequest.onerror = () => observer.error(deleteRequest.error)
      deleteRequest.onsuccess = () => {
        this.deleteSubject.next(id)
        observer.next()
        observer.complete()
      }
    })
  }

  protected withDatabase<T>(operation: (db: IDBDatabase, observer: Subscriber<T>) => void): Observable<T> {
    return new Observable(observer => {
      const request = indexedDB.open(DATABASE_NAME)
      request.onerror = () => observer.error(request.error)

      request.onsuccess = () => {
        try {
          operation(request.result, observer)
        } catch (error) {
          observer.error(error)
        }
      }
    })
  }
}
