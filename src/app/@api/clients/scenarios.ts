import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {isNew, Scenario} from '@api/model';
import {Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class Scenarios {
  constructor(private http: HttpClient) {
  }

  getAll(): Observable<Scenario[]> {
    return this.http.get<Scenario[]>(`/scenarios`)
  }

  get(scenarioId: number): Observable<Scenario> {
    return this.http.get<Scenario>(`/scenarios/${scenarioId}`)
  }

  save(scenario: Scenario): Observable<Scenario> {
    if (isNew(scenario)) {
      return this.http.post<Scenario>(`/scenarios/${scenario.id}`, scenario)
    } else {
      return this.http.put<Scenario>(`/scenarios`, scenario)
    }
  }

  delete(scenarioId: number): Observable<void> {
    return this.http.delete<void>(`/scenarios/${scenarioId}`)
  }
}
