import {v001__init} from './v001__init';
import {MigrationFn} from '../model';

// {[target version]: [migrator]}
export const MIGRATIONS: Record<number, MigrationFn> = {
  1: v001__init,
}

// Current highest version
export const CURRENT_VERSION = 1
