import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Scenario} from './scenarios.model';
import {isNew} from '@api/common';

@Injectable({
  providedIn: 'root'
})
export class Scenarios {
  private http: HttpClient = inject(HttpClient)

  getAll(): Observable<Scenario[]> {
    return this.http.get<Scenario[]>(`/scenarios`)
  }

  get(promptId: number): Observable<Scenario> {
    return this.http.get<Scenario>(`/scenarios/${promptId}`)
  }

  save(prompt: Scenario): Observable<Scenario> {
    if (isNew(prompt)) {
      return this.http.post<Scenario>(`/scenarios`, prompt)
    } else {
      return this.http.put<Scenario>(`/scenarios/${prompt.id}`, prompt)
    }
  }

  delete(promptId: number): Observable<void> {
    return this.http.delete<void>(`/scenarios/${promptId}`)
  }
}
