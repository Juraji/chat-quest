import {inject, Injectable} from '@angular/core';
import {HttpClient, HttpParams} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Memory} from './memories.model';
import {isNew} from '@api/common';

@Injectable({
  providedIn: 'root'
})
export class Memories {
  private http: HttpClient = inject(HttpClient)

  getAll(worldId: number): Observable<Memory[]> {
    return this.http.get<Memory[]>(`/worlds/${worldId}/memories`)
  }

  getAllByCharacter(worldId: number, characterId: number): Observable<Memory[]> {
    return this.http.get<Memory[]>(`/worlds/${worldId}/memories/by-character/${characterId}`)
  }

  save(worldId: number, memory: Memory): Observable<Memory> {
    if (isNew(memory)) {
      return this.http.post<Memory>(`/worlds/${worldId}/memories`, memory)
    } else {
      return this.http.put<Memory>(`/worlds/${worldId}/memories/${memory.id}`, memory)
    }
  }

  delete(worldId: number, memoryId: number): Observable<void> {
    return this.http.delete<void>(`/worlds/${worldId}/memories/${memoryId}`)
  }

  generateMemoriesForMessage(worldId: number, messageId: number, includeNPreceding: number | null): Observable<void> {
    let params = new HttpParams()
    if (!!includeNPreceding) {
      params = params.set("includeNPreceding", includeNPreceding)
    }

    return this.http.post<void>(`/worlds/${worldId}/memories/generate-for-message/${messageId}`, null, {params})
  }
}
