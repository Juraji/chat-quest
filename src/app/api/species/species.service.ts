import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {isNew} from '@api/common';
import {Species} from './species.model';

@Injectable({
  providedIn: 'root'
})
export class SpeciesS {
  private http: HttpClient = inject(HttpClient)

  getAll(): Observable<Species[]> {
    return this.http.get<Species[]>(`/species`)
  }

  get(speciesId: number): Observable<Species> {
    return this.http.get<Species>(`/species/${speciesId}`)
  }

  save(species: Species): Observable<Species> {
    if (isNew(species)) {
      return this.http.post<Species>(`/species`, species)
    } else {
      return this.http.put<Species>(`/species/${species.id}`, species)
    }
  }

  delete(speciesId: number): Observable<void> {
    return this.http.delete<void>(`/species/${speciesId}`)
  }
}
