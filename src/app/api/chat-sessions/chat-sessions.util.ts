import {ChatSession} from './chat-sessions.model';

export function chatSessionSortingTransformer(characters: ChatSession[]): ChatSession[] {
  return characters.sort((a, b) => {
    return b.createdAt!.localeCompare(a.createdAt!);
  })
}

export function sessionEntityFilter<T extends { chatSessionId: number }>(
  chatSessionId: () => number,
): (input: T) => boolean {
  return input => input.chatSessionId === chatSessionId()
}
