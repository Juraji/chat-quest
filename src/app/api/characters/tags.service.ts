import {inject, Injectable, Signal, signal, WritableSignal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {defer, Observable, of, throwError} from 'rxjs';
import {Tag, TagCreated, TagDeleted, TagUpdated} from './tags.model';
import {isNew} from '@api/common';
import {SSE} from '@api/sse';
import {arrayAdd, arrayRemove, arrayUpdate} from '@util/array';

@Injectable({
  providedIn: 'root'
})
export class Tags {
  private http: HttpClient = inject(HttpClient)
  private readonly sse = inject(SSE)

  private readonly tagCache: WritableSignal<Tag[]> = signal([]);
  readonly all: Signal<Tag[]> = this.tagCache

  constructor() {
    this.setupTagCache()
  }

  /** @deprecated */
  getAll(): Observable<Tag[]> {
    return defer(() => [this.tagCache()])
  }

  get(tagId: number): Observable<Tag> {
    return defer(() => {
      const tags = this.tagCache()
      const tag = tags.find(tag => tag.id === tagId);
      if (tag) {
        return of(tag)
      } else {
        return throwError(() => new Error(`No tag with id ${tagId}`));
      }
    })
  }

  save(tag: Tag): Observable<Tag> {
    if (isNew(tag)) {
      return this.http
        .post<Tag>(`/tags`, tag)
    } else {
      return this.http
        .put<Tag>(`/tags/${tag.id}`, tag)
    }
  }

  delete(tagId: number): Observable<void> {
    return this.http
      .delete<void>(`/tags/${tagId}`)
  }

  private setupTagCache() {
    // Hydrate
    this.http
      .get<Tag[]>(`/tags`)
      .subscribe(tags => this.tagCache.set(tags));

    // Listen for changes
    this.sse
      .on(TagCreated)
      .subscribe(tag => this.tagCache.update(cache =>
        arrayAdd(cache, tag)))
    this.sse
      .on(TagUpdated)
      .subscribe(tag => {
        this.tagCache.update(cache =>
          arrayUpdate(cache, t => ({...t, ...tag}), t => t.id === tag.id));
      })
    this.sse
      .on(TagDeleted)
      .subscribe(tagId => {
        this.tagCache.update(cache =>
          arrayRemove(cache, t => t.id === tagId));
      })
  }
}
