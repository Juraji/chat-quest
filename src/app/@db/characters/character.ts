import {StoreRecord} from "@db/model";

export interface Character extends StoreRecord {
  name: string
  appearance: string
  personality: string
  history: string
  likelyActions: string[]
  unlikelyActions: string[]
  dialogueExamples: string[]
  extraTraits: Record<string, string>
  avatar: Blob | null
  favorite: boolean,
  tagIds: number[]
}
