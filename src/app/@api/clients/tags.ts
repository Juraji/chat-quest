import {inject, Injectable, signal, Signal, WritableSignal} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable, tap} from 'rxjs';
import {isNew, Tag} from '@api/model';

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
      .pipe(tap(tag => this._cachedTags.update(prev => {
        return prev.some(t => t.id === tag.id)
          ? prev
          : [...prev, tag]
      })))
  }

  save(tag: Tag): Observable<Tag> {
    if (isNew(tag)) {
      return this.http
        .post<Tag>(`/tags`, tag)
        .pipe(tap(tag => this._cachedTags.update(prev => [...prev, tag])))
    } else {
      return this.http
        .put<Tag>(`/tags/${tag.id}`, tag)
        .pipe(tap(tag => this._cachedTags.update(prev => {
          const idx = prev.findIndex(t => t.id === tag.id);
          return idx === -1
            ? [...prev, tag]
            : [...prev.slice(0, idx), tag, ...prev.slice(idx + 1)]
        })))
    }
  }

  delete(tagId: number): Observable<void> {
    return this.http
      .delete<void>(`/tags/${tagId}`)
      .pipe(tap(() => this._cachedTags.update(prev => {
        const idx = prev.findIndex(t => t.id === tagId);
        return idx === -1
          ? prev
          : prev.slice().splice(idx, 1)
      })))
  }
}
