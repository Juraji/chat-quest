import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {World} from './worlds.model';
import {isNew} from '@api/common';

@Injectable({
  providedIn: 'root'
})
export class Worlds {
  private http: HttpClient = inject(HttpClient)

  getAll(): Observable<World[]> {
    return this.http.get<World[]>("/worlds");
  }

  get(worldId: number): Observable<World> {
    return this.http.get<World>(`/worlds/${worldId}`);
  }

  save(world: World): Observable<World> {
    if (isNew(world)) {
      return this.http.post<World>(`/worlds`, world)
    } else {
      return this.http.put<World>(`/worlds/${world.id}`, world)
    }
  }

  delete(worldId: number): Observable<void> {
    return this.http.delete<void>(`/worlds/${worldId}`)
  }
}
