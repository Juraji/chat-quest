import {CharaCard} from './charaCard';

/**
 * @see https://github.com/kwaroran/character-card-spec-v3/blob/main/SPEC_V3.md
 */

export interface CharaCardV3 extends CharaCard {
  spec: 'chara_card_v3'
  spec_version: '3.0'
  data: CharaCardV3Data
}

export interface CharaCardV3Data {
  name: string
  description: string
  tags: Array<string>
  creator: string
  character_version: string
  mes_example: string
  extensions: Record<string, any>
  system_prompt: string
  post_history_instructions: string
  first_mes: string
  alternate_greetings: Array<string>
  personality: string
  scenario: string
  creator_notes: string
  character_book?: CharaCardV3Lorebook
  assets?: CharaCardV3Asset[]
  nickname?: string
  creator_notes_multilingual?: Record<string, string>
  source?: string[]
  group_only_greetings: Array<string>
  creation_date?: number
  modification_date?: number
}

export interface CharaCardV3Asset {
  type: string
  uri: string
  name: string
  ext: string
}

export interface CharaCardV3Lorebook {
  name?: string
  description?: string
  scan_depth?: number
  token_budget?: number
  recursive_scanning?: boolean
  extensions: Record<string, any>
  entries: CharaCardV3LorebookEntry[]
}

export interface CharaCardV3LorebookEntry {
  keys: Array<string>
  content: string
  extensions: Record<string, any>
  enabled: boolean
  insertion_order: number
  case_sensitive?: boolean
  use_regex: boolean
  constant?: boolean
  name?: string
  priority?: number
  id?: number | string
  comment?: string
  selective?: boolean
  secondary_keys?: Array<string>
  position?: 'before_char' | 'after_char'
}
