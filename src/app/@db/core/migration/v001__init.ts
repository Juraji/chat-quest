import {MigrationFn} from '../model';

export const v001__init: MigrationFn = async db => {
  db.createObjectStore('settings')
  db.createObjectStore('characters', {autoIncrement: true, keyPath: 'id'})
  db.createObjectStore('tags', {autoIncrement: true, keyPath: 'id'})
  db.createObjectStore('system-prompts', {autoIncrement: true, keyPath: 'id'})
  db.createObjectStore('scenarios', {autoIncrement: true, keyPath: 'id'});
  db.createObjectStore('worlds', {autoIncrement: true, keyPath: 'id'});
}
