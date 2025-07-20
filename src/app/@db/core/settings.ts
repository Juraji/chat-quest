import {Injectable} from '@angular/core';
import {filter, Observable, Subject, Subscriber} from 'rxjs';
import {DATABASE_NAME} from '@db/core/migration';

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

  get<T>(name: string, watch: boolean = false): Observable<T | null> {
    return this.withDatabase<T>((db, observer) => {
      const getRequest = db
        .transaction(this.storeName, 'readonly')
        .objectStore(this.storeName)
        .get(name)

      getRequest.onerror = () => observer.error(getRequest.error)
      getRequest.onsuccess = () => {
        // Emit initial data or complete if not found
        if (getRequest.result) {
          observer.next(getRequest.result)
        } else {
          observer.complete()
          return
        }

        // Setup watch or complete
        if (watch) {
          const modifySub = this.modifySubject
            .pipe(filter(it => it.name === name))
            .subscribe(record => observer.next(record.value))

          observer.add((() => modifySub.unsubscribe()))
        } else {
          observer.complete()
        }
      }
    })
  }

  set(name: string, value: any): Observable<void> {
    return this.withDatabase((db, observer) => {
      const transaction = db.transaction(this.storeName, 'readwrite')
      const store = transaction.objectStore(this.storeName)

      const handleRequest: (r: IDBRequest) => void = r => {
        r.onerror = () => observer.error(r.error)
        r.onsuccess = () => {
          observer.next()
          observer.complete()
          this.modifySubject.next({name, value})
        }
      }

      // In order to upsert, we need to query for existing first, then add or put.
      const existsRequest = store.get(name)
      existsRequest.onerror = () => observer.error(existsRequest.error)
      existsRequest.onsuccess = () => {
        if (existsRequest.result) {
          // Existing record, update
          handleRequest(store.put(value, name))
        } else {
          // New record, add
          handleRequest(store.add(value, name))
        }
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
