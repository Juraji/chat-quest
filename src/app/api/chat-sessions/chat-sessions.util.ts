import {ChatSession} from './chat-sessions.model';

export function chatSessionSortingTransformer(characters: ChatSession[]): ChatSession[] {
  return characters.sort((a, b) => {
    return b.createdAt!.localeCompare(a.createdAt!);
  })
}
