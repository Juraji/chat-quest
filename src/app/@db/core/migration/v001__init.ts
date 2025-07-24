import {MigrationFn, PostMigrationFn} from '../model';
import {firstValueFrom} from 'rxjs';

export const v001__init: MigrationFn = async db => {
  db.createObjectStore('settings')
  db.createObjectStore('characters', {autoIncrement: true, keyPath: 'id'})

  const tagsStore = db.createObjectStore('tags', {autoIncrement: true, keyPath: 'id'})
  tagsStore.createIndex('lowercase', 'lowercase', {unique: true})

  db.createObjectStore('system-prompts', {autoIncrement: true, keyPath: 'id'})
  db.createObjectStore('scenarios', {autoIncrement: true, keyPath: 'id'});
  db.createObjectStore('worlds', {autoIncrement: true, keyPath: 'id'});
}

export const v001__POST__init: PostMigrationFn = async (db, http) => {
  const defaultPrompts = await firstValueFrom(http.get<any[]>("/data/default-system-prompts.json"))
  const store = db
    .transaction('system-prompts', 'readwrite')
    .objectStore('system-prompts')

  for (const prompt of defaultPrompts) {
    store.add(prompt)
  }
}
