// Current highest version
export const CURRENT_VERSION = 1
const DEFAULT_STORE_OPTS: IDBObjectStoreParameters = {autoIncrement: true, keyPath: 'id'}

type MigrationFn = (db: IDBDatabase) => Promise<void>

// {[target version]: [migrator]}
export const MIGRATIONS: Record<number, MigrationFn> = {
  1: async db => {
    db.createObjectStore('settings')
    db.createObjectStore('characters', DEFAULT_STORE_OPTS)
    db.createObjectStore('tags', DEFAULT_STORE_OPTS)
    db.createObjectStore('system-prompts', DEFAULT_STORE_OPTS)
  },
}
