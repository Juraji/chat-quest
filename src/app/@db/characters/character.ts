import {StoreRecord} from "@db/core";

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
