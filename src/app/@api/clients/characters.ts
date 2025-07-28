import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Character, CharacterDetails, CharacterTextBlock, isNew, Tag} from '@api/model';

@Injectable({
  providedIn: 'root'
})
export class Characters {
  constructor(private http: HttpClient) {
  }

  getAll(): Observable<Character[]> {
    return this.http.get<Character[]>('/characters')
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

  addTag(characterId: number, tagId: number): Observable<void> {
    return this.http.post<void>(`/characters/${characterId}/tags/${tagId}`, null)
  }

  removeTag(characterId: number, tagId: number): Observable<void> {
    return this.http.delete<void>(`/characters/${characterId}/tags/${tagId}`)
  }

  getDialogueExamples(characterId: number): Observable<CharacterTextBlock[]> {
    return this.http.get<CharacterTextBlock[]>(`/characters/${characterId}/dialogue-examples`)
  }

  saveDialogueExamples(characterId: number, dialogueExamples: string[]): Observable<void> {
    return this.http.post<void>(`/characters/${characterId}/dialogue-examples`, dialogueExamples)
  }

  getGreetings(characterId: number): Observable<CharacterTextBlock[]> {
    return this.http.get<CharacterTextBlock[]>(`/characters/${characterId}/greetings`)
  }

  saveGreetings(characterId: number, greetings: string[]): Observable<void> {
    return this.http.post<void>(`/characters/${characterId}/greetings`, greetings)
  }

  getGroupGreetings(characterId: number): Observable<CharacterTextBlock[]> {
    return this.http.get<CharacterTextBlock[]>(`/characters/${characterId}/group-greetings`)
  }

  saveGroupGreetings(characterId: number, greetings: string[]): Observable<void> {
    return this.http.post<void>(`/characters/${characterId}/group-greetings`, greetings)
  }
}
