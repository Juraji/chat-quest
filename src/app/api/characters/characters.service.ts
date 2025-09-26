import {computed, inject, Injectable, Signal, signal, WritableSignal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Character, CharacterCreated, CharacterDeleted, CharacterUpdated} from './characters.model';
import {isNew} from '@api/common';
import {SSE} from '@api/sse';
import {arrayAdd, arrayRemove, arrayUpdate} from '@util/array';

@Injectable({
  providedIn: 'root'
})
export class Characters {
  private http: HttpClient = inject(HttpClient)
  private readonly sse = inject(SSE)

  private readonly lvCache: WritableSignal<Character[]> = signal([]);
  readonly all: Signal<Character[]> = computed(() => this
    .lvCache().slice().sort((a, b) =>
      a.favorite === b.favorite
        ? a.name.localeCompare(b.name)
        : a.favorite ? -1 : 1))

  constructor() {
    this.setupLVCache()
  }

  listViewBy(idFn: () => Nullable<number>): Signal<Nullable<Character>> {
    return computed(() => {
      const id = idFn()
      if (!!id) {
        const chars = this.lvCache()
        return chars.find(char => char.id === id)
      } else {
        return null
      }
    })
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

  duplicate(characterId: number): Observable<Character> {
    return this.http.post<Character>(`/characters/${characterId}/duplicate`, null)
  }

  delete(characterId: number): Observable<void> {
    return this.http.delete<void>(`/characters/${characterId}`)
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

  private setupLVCache() {
    // Hydrate
    this.http
      .get<Character[]>('/characters')
      .subscribe(chars => this.lvCache.set(chars))

    // Listen for changes
    this.sse
      .on(CharacterCreated)
      .subscribe(char => this.lvCache.update(cache =>
        arrayAdd(cache, char)))
    this.sse
      .on(CharacterUpdated)
      .subscribe(char => this.lvCache.update(cache =>
        arrayUpdate(cache, c => ({...c, ...char}), c => c.id === char.id)))
    this.sse
      .on(CharacterDeleted)
      .subscribe(id => this.lvCache.update(cache =>
        arrayRemove(cache, id)))
  }
}
