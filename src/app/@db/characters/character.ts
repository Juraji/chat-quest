import {StoreRecord} from "@db/core";

export interface Character extends StoreRecord {
  name: string
  appearance: string
  personality: string
  history: string
  likelyActions: string[]
  unlikelyActions: string[]
  dialogueExamples: string[]
  avatar: Blob | null
  favorite: boolean,
  tagIds: number[]
}
