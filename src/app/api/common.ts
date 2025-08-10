export interface ChatQuestModel {
  id: number
}

export const NEW_ID = 0

export function isNew(m: ChatQuestModel | null | undefined): boolean {
  return !m || m.id === NEW_ID;
}

export function entityIdFilter<T extends ChatQuestModel | number>(
  entityId: () => number
): (input: T) => boolean {
  return input => typeof input === "number"
    ? input === entityId()
    : input.id === entityId()
}
