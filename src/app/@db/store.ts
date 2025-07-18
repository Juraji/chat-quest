import {Observable, Subscriber, toArray} from "rxjs";
import {NewRecord, StoreRecord} from "@db/model";
import {DATABASE_NAME} from '@db/init';

export abstract class Store<T extends StoreRecord> {

  protected constructor(readonly storeName: string) {
  }

  getAll(): Observable<T[]> {
    return this.withDatabase<T>((db, observer) => {
      const cursorRequest = db
        .transaction(this.storeName, 'readonly')
        .objectStore(this.storeName)
        .openCursor()

      cursorRequest.onerror = () => observer.error(cursorRequest.error)
      cursorRequest.onsuccess = () => {
        const cursor = cursorRequest.result
        if (!!cursor) {
          observer.next(cursor.value)
          cursor.continue()
        } else {
          observer.complete()
        }
      }
    }).pipe(toArray())
  }

  get(id: number): Observable<T> {
    return this.withDatabase<T>((db, observer) => {
      const request = db
        .transaction(this.storeName, 'readonly')
        .objectStore(this.storeName)
        .get(id);

      request.onerror = () => observer.error(request.error);
      request.onsuccess = () => {
        if (request.result) {
          observer.next(request.result as T);
        }
        observer.complete();
      };
    });
  }

  save(record: NewRecord<T> | T): Observable<T> {
    if ('id' in record && !!record.id) {
      return this.updateRecord(record)
    } else {
      return this.createRecord(record)
    }
  }

  delete(id: number): Observable<void> {
    return this.withDatabase((db, observer) => {
      const deleteRequest = db
        .transaction(this.storeName, 'readwrite')
        .objectStore(this.storeName)
        .delete(id)

      deleteRequest.onerror = () => observer.error(deleteRequest.error)

      deleteRequest.onsuccess = () => {
        observer.next()
        observer.complete()
      }
    })
  }


  protected withDatabase<T>(operation: (db: IDBDatabase, observer: Subscriber<T>) => void): Observable<T> {
    return new Observable(observer => {
      const request = indexedDB.open(DATABASE_NAME);
      request.onerror = () => observer.error(request.error)

      request.onsuccess = () => {
        try {
          operation(request.result, observer);
        } catch (error) {
          observer.error(error);
        }
      };
    });
  }

  private createRecord(record: NewRecord<T>): Observable<T> {
    return this.withDatabase<T>((db, observer) => {
      const request = db
        .transaction(this.storeName, 'readwrite')
        .objectStore(this.storeName)
        .add(record)

      request.onerror = () => observer.error(request.error);

      request.onsuccess = () => {
        const newRecord = {...record, id: request.result} as T;
        console.log(newRecord)
        observer.next(newRecord);
        observer.complete();
      };
    });
  }

  private updateRecord(record: T): Observable<T> {
    return this.withDatabase((db, observer) => {
      const putRequest = db
        .transaction(this.storeName, 'readwrite')
        .objectStore(this.storeName)
        .put(record)

      putRequest.onerror = () => observer.error(putRequest.error)

      putRequest.onsuccess = () => {
        observer.next(record)
        observer.complete()
      }
    })
  }
}
