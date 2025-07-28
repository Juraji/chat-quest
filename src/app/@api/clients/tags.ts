import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {isNew, Tag} from '@api/model';

@Injectable({
  providedIn: 'root'
})
export class Tags {
  constructor(private http: HttpClient) {
  }

  getAll(): Observable<Tag[]> {
    return this.http.get<Tag[]>(`/tags`)
  }

  get(tagId: number): Observable<Tag> {
    return this.http.get<Tag>(`/tags/${tagId}`)
  }

  save(tag: Tag): Observable<Tag> {
    if (isNew(tag)) {
      return this.http.post<Tag>(`/system-prompts`, tag)
    } else {
      return this.http.put<Tag>(`/system-prompts/${tag.id}`, tag)
    }
  }

  delete(tagId: number): Observable<void> {
    return this.http.delete<void>(`/system-prompts/${tagId}`)
  }
}
