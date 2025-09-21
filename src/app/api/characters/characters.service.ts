import {computed, inject, Injectable, Signal, signal, WritableSignal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {map, mergeMap, Observable} from 'rxjs';
import {Character, CharacterCreated, CharacterDeleted, CharacterListView, CharacterUpdated} from './characters.model';
import {isNew} from '@api/common';
import {CharacterTagAdded, CharacterTagRemoved, Tag, TagDeleted, TagUpdated} from './tags.model';
import {SSE} from '@api/sse';
import {arrayAdd, arrayRemove, arrayReplace, arrayUpdate} from '@util/array';
import {Tags} from '@api/characters/tags.service';

@Injectable({
  providedIn: 'root'
})
export class Characters {
  private http: HttpClient = inject(HttpClient)
  private readonly tags = inject(Tags)
  private readonly sse = inject(SSE)

  private readonly lvCache: WritableSignal<CharacterListView[]> = signal([]);
  readonly all: Signal<CharacterListView[]> = computed(() => this
    .lvCache().slice().sort((a, b) =>
      a.favorite === b.favorite
        ? a.name.localeCompare(b.name)
        : a.favorite ? -1 : 1))

  constructor() {
    this.setupLVCache()
  }

  listViewBy(idFn: () => Nullable<number>): Signal<Nullable<CharacterListView>> {
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

  private setupLVCache() {
    // Hydrate
    this.http
      .get<CharacterListView[]>('/characters')
      .subscribe(chars => this.lvCache.set(chars))

    // Listen for changes
    this.sse
      .on(CharacterCreated)
      .subscribe(char => {
        const lv: CharacterListView = {...char, tags: []}
        this.lvCache.update(cache =>
          arrayAdd(cache, lv))
      })
    this.sse
      .on(CharacterUpdated)
      .subscribe(char => {
        this.lvCache.update(cache =>
          arrayUpdate(cache, c => ({...c, ...char}), c => c.id === char.id))
      })
    this.sse
      .on(CharacterDeleted)
      .subscribe(id => this.lvCache.update(cache =>
        arrayRemove(cache, id)))

    this.sse
      .on(CharacterTagAdded)
      .pipe(mergeMap(([charId, tagId]) => this.tags
        .get(tagId)
        .pipe(map(t => ([charId, t]) as [number, Tag]))
      ))
      .subscribe(([charId, tag]) => this.lvCache.update(cache =>
        arrayUpdate(cache, c => ({...c, tags: arrayAdd(c.tags, tag)}), c => c.id === charId)))
    this.sse
      .on(CharacterTagRemoved)
      .subscribe(([charId, tagId]) => this.lvCache.update(cache =>
        arrayUpdate(cache, c => ({...c, tags: arrayRemove(c.tags, t => t.id === tagId)}), c => c.id === charId)))

    this.sse
      .on(TagUpdated)
      .subscribe(tag =>
        this.lvCache.update(cache => cache
          .map(c => ({...c, tags: arrayReplace(c.tags, tag, t => t.id === tag.id)}),)))
    this.sse
      .on(TagDeleted)
      .subscribe(tagId =>
        this.lvCache.update(cache => cache
          .map(c => ({...c, tags: arrayRemove(c.tags, t => t.id === tagId)}),)))
  }
}
