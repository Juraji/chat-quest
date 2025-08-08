export interface ChatQuestModel {
  id: number
}

export const NEW_ID = 0

export function isNew(m: ChatQuestModel | null | undefined): boolean {
  return !m || m.id === NEW_ID;
}
