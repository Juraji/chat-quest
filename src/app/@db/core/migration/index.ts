import {MigrationFn} from '../model';
import {v001__init} from './v001__init';

// {[target version]: [migrator]}
export const MIGRATIONS: Record<number, MigrationFn> = {
  1: v001__init,
}

// Current highest version
export const CURRENT_VERSION = 1
export const DATABASE_NAME = 'ChatQuestStore'
