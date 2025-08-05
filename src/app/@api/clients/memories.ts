import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {isNew, Memory} from '@api/model';
import {Observable} from 'rxjs';

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
}
