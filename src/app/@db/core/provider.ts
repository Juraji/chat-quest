import {CURRENT_VERSION, DATABASE_NAME, MIGRATIONS, POST_MIGRATIONS} from './migration';
import {IDBPDatabase, openDB} from 'idb';
import {InjectionToken, Provider} from '@angular/core';
import {HttpClient} from '@angular/common/http';

export const DATABASE = new InjectionToken<Promise<IDBPDatabase>>('DATABASE');

export function provideDatabase(): Provider {
  return {
    provide: DATABASE,
    useFactory: initializeDatabase,
    deps: [HttpClient]
  }
}

export async function initializeDatabase(http: HttpClient): Promise<IDBPDatabase> {
  console.log('Initializing database...', {DATABASE_NAME, CURRENT_VERSION});
  let previousVersion = CURRENT_VERSION
  let currentVersion = CURRENT_VERSION

  const db = await openDB(DATABASE_NAME, CURRENT_VERSION, {
    blocking: () => window.location.reload(),
    upgrade: async (db, oldVersion, newVersion) => {
      if (!newVersion) return
      console.log(`Migrating DB from ${oldVersion} to ${newVersion}`);

      previousVersion = oldVersion;
      currentVersion = newVersion;

      for (let i = oldVersion + 1; i <= newVersion; i++) {
        const migration = MIGRATIONS[i]
        if (!!migration) await migration(db)
      }
    }
  })

  if (currentVersion !== CURRENT_VERSION) {
    console.log(`Running post migrations from ${previousVersion} to ${currentVersion}`);
    for (let i = previousVersion + 1; i <= currentVersion; i++) {
      const postmigration = POST_MIGRATIONS[i]
      if (!!postmigration) await postmigration(db, http)
    }
  }

  return db
}

