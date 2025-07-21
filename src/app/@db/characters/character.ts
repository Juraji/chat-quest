import {NewRecord, StoreRecord} from "@db/core";

export interface Character extends StoreRecord {
  //base
  name: string
  appearance: string
  personality: string
  avatar: Blob | null
  favorite: boolean
  tagIds: number[]

  // Extended
  history: string
  likelyActions: string[]
  unlikelyActions: string[]
  dialogueExamples: string[]

  // Chat Defaults
  scenario: string
  firstMessage: string
  alternateGreetings: string[]
  groupGreetings: string[]
  groupTalkativeness: number
}

export const NEW_CHARACTER: NewRecord<Character> = {
  name: '',
  appearance: '',
  personality: '',
  avatar: null,
  favorite: false,
  tagIds: [],

  // Extended
  history: '',
  likelyActions: [],
  unlikelyActions: [],
  dialogueExamples: [],

  // Chat Defaults
  scenario: '',
  firstMessage: '',
  alternateGreetings: [],
  groupGreetings: [],
  groupTalkativeness: 0.5
}
