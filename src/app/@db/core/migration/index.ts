import {MigrationFn, PostMigrationFn} from '../model';
import {v001__init, v001__POST__init} from './v001__init';

// {[target version]: [migrator]}
export const MIGRATIONS: Record<number, MigrationFn> = {
  1: v001__init,
}

export const POST_MIGRATIONS: Record<number, PostMigrationFn> = {
  1: v001__POST__init,
}

// Current highest version
export const CURRENT_VERSION = 1
export const DATABASE_NAME = 'ChatQuestStore'
