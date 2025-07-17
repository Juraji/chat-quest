import {Observable, Subscriber} from "rxjs";
import {DATABASE_NAME} from "./init";
import {NewRecord, StoreRecord} from '@db/storeRecord';

function withDatabase<T>(operation: (db: IDBDatabase, observer: Subscriber<T>) => void): Observable<T> {
  return new Observable(observer => {
    const request = indexedDB.open(DATABASE_NAME);
    request.onerror = () => observer.error(request.error)

    request.onsuccess = () => {
      try {
        operation(request.result, observer);
        observer.complete();
      } catch (error) {
        observer.error(error);
      }
    };
  });
}

export function getRecord<T extends StoreRecord>(storeName: string, id: number): Observable<T> {
  return withDatabase<T>((db, observer) => {
    const request = db
      .transaction(storeName, 'readonly')
      .objectStore(storeName)
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

export function getAllRecords<T extends StoreRecord>(storeName: string): Observable<T> {
  return withDatabase<T>((db, observer) => {
    const cursorRequest = db
      .transaction(storeName, 'readonly')
      .objectStore(storeName)
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
  })
}

export function scanRecords<T extends StoreRecord>(storeName: string, indexName: string, key: any): Observable<T> {
  return withDatabase((db, observer) => {
    const keyRange = IDBKeyRange.only(key)
    const cursorRequest = db
      .transaction(storeName, 'readonly')
      .objectStore(storeName)
      .index(indexName)
      .openCursor(keyRange)

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
  })
}

export function createRecord<T extends StoreRecord>(storeName: string, record: NewRecord<T>): Observable<T> {
  return withDatabase<T>((db, observer) => {
    const request = db
      .transaction(storeName, 'readwrite')
      .objectStore(storeName)
      .add(record)

    request.onerror = () => observer.error(request.error);

    request.onsuccess = () => {
      const newRecord = {...record, id: request.result} as T;
      observer.next(newRecord);
      observer.complete();
    };
  });
}

export function updateRecord<T extends StoreRecord>(storeName: string, record: T): Observable<T> {
  return withDatabase((db, observer) => {
    const putRequest = db
      .transaction(storeName, 'readwrite')
      .objectStore(storeName)
      .put(record)

    putRequest.onerror = () => observer.error(putRequest.error)

    putRequest.onsuccess = () => {
      observer.next(record)
      observer.complete()
    }
  })
}

export function saveRecord<T extends StoreRecord>(storeName: string, record: T | NewRecord<T>): Observable<T> {
  if ('id' in record && !!record.id) {
    return updateRecord(storeName, record)
  } else {
    return createRecord(storeName, record)
  }
}

export function deleteRecord(storeName: string, id: number): Observable<void> {
  return withDatabase((db, observer) => {
    const deleteRequest = db
      .transaction(storeName, 'readwrite')
      .objectStore(storeName)
      .delete(id)

    deleteRequest.onerror = () => observer.error(deleteRequest.error)

    deleteRequest.onsuccess = () => {
      observer.next()
      observer.complete()
    }
  })
}
