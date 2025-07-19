import {CURRENT_VERSION, MIGRATIONS} from './migrations';

export const DATABASE_NAME = 'ChatQuestStore'

/**
 * Initializes the IndexedDB database and applies necessary migrations.
 *
 * This function opens the database with the specified name and current version. If the database doesn't exist,
 * it will be created. The onupgradeneeded event handler is crucial for database migrations:
 *
 * When a new version of the database is requested (either by explicitly setting a higher version or
 * when the user clears their browser data), this function will execute all migration functions from
 * the old version to the new version in ascending order. Each migration function receives the database
 * object and should perform schema changes like creating new object stores, adding indexes, etc.
 *
 * The MIGRATIONS constant defines these migration functions as an object where keys are target versions
 * and values are functions that perform the migration for that specific version.
 *
 * @returns A promise that resolves when the database is successfully initialized or rejects on error.
 */
export function initializeDatabase(): Promise<void> {
    return new Promise((resolve, reject) => {
        const request = indexedDB.open(DATABASE_NAME, CURRENT_VERSION);
        request.onerror = () => reject(request.error);
        request.onsuccess = () => resolve()

        request.onupgradeneeded = e => {
            if (!e.newVersion) return
            const db: IDBDatabase = request.result;
            for (let i = e.oldVersion + 1; i <= e.newVersion; i++) {
                const migration = MIGRATIONS[i]
                if (!!migration) migration(db)
            }
        }
    });
}

