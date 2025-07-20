import {IDBPDatabase} from 'idb';

export interface StoreRecord {
  id: number
}

export type NewRecord<T extends StoreRecord> = Omit<T, 'id'>
export type MigrationFn = (db: IDBPDatabase) => Promise<void>
