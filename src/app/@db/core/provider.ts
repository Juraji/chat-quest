import {CURRENT_VERSION, DATABASE_NAME, MIGRATIONS} from './migration';
import {IDBPDatabase, openDB} from 'idb';
import {InjectionToken, Provider} from '@angular/core';

export const DATABASE = new InjectionToken<Promise<IDBPDatabase>>('DATABASE');

export function provideDatabase(): Provider {
  return {
    provide: DATABASE,
    useFactory: initializeDatabase
  }
}

export function initializeDatabase(): Promise<IDBPDatabase> {
  console.log('Initializing database...', {DATABASE_NAME, CURRENT_VERSION});
  return openDB(DATABASE_NAME, CURRENT_VERSION, {
    blocking: () => window.location.reload(),
    upgrade: async (database, oldVersion, newVersion) => {
      console.log(`Migrating DB from ${oldVersion} to ${newVersion}`);
      if (!newVersion) return
      for (let i = oldVersion + 1; i <= newVersion; i++) {
        const migration = MIGRATIONS[i]
        if (!!migration) await migration(database)
      }
    }
  })
}

