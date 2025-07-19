// Current highest version
export const CURRENT_VERSION = 1
const DEFAULT_STORE_OPTS: IDBObjectStoreParameters = {autoIncrement: true, keyPath: 'id'}

// {[target version]: [migrator]}
export const MIGRATIONS: Record<number, (db: IDBDatabase) => void> = {
  1: db => {
    db.createObjectStore('settings')
    db.createObjectStore('characters', DEFAULT_STORE_OPTS)
    db.createObjectStore('tags', DEFAULT_STORE_OPTS)
    db.createObjectStore('system-prompts', DEFAULT_STORE_OPTS)
  },
}
