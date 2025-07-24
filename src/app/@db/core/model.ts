import {IDBPDatabase} from 'idb';
import {HttpClient} from '@angular/common/http';

export interface StoreRecord {
  id: number
}

export type NewRecord<T extends StoreRecord> = Omit<T, 'id'>
export type MigrationFn = (db: IDBPDatabase) => Promise<void>
export type PostMigrationFn = (db: IDBPDatabase, http: HttpClient) => Promise<void>
