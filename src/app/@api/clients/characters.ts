import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Character, CharacterDetails, CharacterWithTags, isNew, Tag} from '@api/model';

@Injectable({
  providedIn: 'root'
})
export class Characters {
  constructor(private http: HttpClient) {
  }

  getAll(): Observable<Character[]> {
    return this.http.get<Character[]>('/characters')
  }

  getAllWithTags(): Observable<CharacterWithTags[]> {
    return this.http.get<CharacterWithTags[]>('/characters/with-tags')
  }

  get(characterId: number): Observable<Character> {
    return this.http.get<Character>(`/characters/${characterId}`)
  }

  save(character: Character): Observable<Character> {
    if (isNew(character)) {
      return this.http.post<Character>(`/characters`, character)
    } else {
      return this.http.put<Character>(`/characters/${character.id}`, character)
    }
  }

  delete(characterId: number): Observable<void> {
    return this.http.delete<void>(`/characters/${characterId}`)
  }

  getDetails(characterId: number): Observable<CharacterDetails> {
    return this.http.get<CharacterDetails>(`/characters/${characterId}/details`)
  }

  saveDetails(details: CharacterDetails): Observable<CharacterDetails> {
    return this.http.put<CharacterDetails>(`/characters/${details.characterId}/details`, details)
  }

  getTags(characterId: number): Observable<Tag[]> {
    return this.http.get<Tag[]>(`/characters/${characterId}/tags`)
  }

  saveTags(characterId: number, tagIds: number[]): Observable<void> {
    return this.http.post<void>(`/characters/${characterId}/tags`, tagIds)
  }

  getDialogueExamples(characterId: number): Observable<string[]> {
    return this.http.get<string[]>(`/characters/${characterId}/dialogue-examples`)
  }

  saveDialogueExamples(characterId: number, dialogueExamples: string[]): Observable<void> {
    return this.http.post<void>(`/characters/${characterId}/dialogue-examples`, dialogueExamples)
  }

  getGreetings(characterId: number): Observable<string[]> {
    return this.http.get<string[]>(`/characters/${characterId}/greetings`)
  }

  saveGreetings(characterId: number, greetings: string[]): Observable<void> {
    return this.http.post<void>(`/characters/${characterId}/greetings`, greetings)
  }

  getGroupGreetings(characterId: number): Observable<string[]> {
    return this.http.get<string[]>(`/characters/${characterId}/group-greetings`)
  }

  saveGroupGreetings(characterId: number, greetings: string[]): Observable<void> {
    return this.http.post<void>(`/characters/${characterId}/group-greetings`, greetings)
  }
}
