import {inject, Injectable, signal, Signal, WritableSignal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, tap} from 'rxjs';
import {Tag} from './tags.model';
import {isNew} from '@api/common';
import {arrayAddItem, arrayRemoveItem, arrayUpsertItem} from '@util/array';

@Injectable({
  providedIn: 'root'
})
export class Tags {
  private http: HttpClient = inject(HttpClient)

  private readonly _cachedTags: WritableSignal<Tag[]> = signal([]);
  readonly cachedTags: Signal<Tag[]> = this._cachedTags;

  constructor() {
    this.getAll().subscribe()
  }

  getAll(): Observable<Tag[]> {
    return this.http
      .get<Tag[]>(`/tags`)
      .pipe(tap(tags => this._cachedTags.set(tags)))
  }

  get(tagId: number): Observable<Tag> {
    return this.http
      .get<Tag>(`/tags/${tagId}`)
      .pipe(tap(tag => this._cachedTags
        .update(prev => arrayUpsertItem(prev, tag, ({id}) => id === tag.id))))
  }

  save(tag: Tag): Observable<Tag> {
    if (isNew(tag)) {
      return this.http
        .post<Tag>(`/tags`, tag)
        .pipe(tap(tag => this._cachedTags
          .update(prev => arrayAddItem(prev, tag))))
    } else {
      return this.http
        .put<Tag>(`/tags/${tag.id}`, tag)
        .pipe(tap(tag => this._cachedTags
          .update(prev => arrayUpsertItem(prev, tag, ({id}) => id === tag.id))))
    }
  }

  delete(tagId: number): Observable<void> {
    return this.http
      .delete<void>(`/tags/${tagId}`)
      .pipe(tap(() => this._cachedTags
        .update(prev => arrayRemoveItem(prev, ({id}) => id === tagId))))
  }
}
