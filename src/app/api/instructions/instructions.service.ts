import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {Instruction} from './instructions.model';
import {isNew} from '@api/common';

@Injectable({
  providedIn: 'root'
})
export class Instructions {
  private http: HttpClient = inject(HttpClient)

  getAll(): Observable<Instruction[]> {
    return this.http.get<Instruction[]>(`/instruction-templates`)
  }

  get(promptId: number): Observable<Instruction> {
    return this.http.get<Instruction>(`/instruction-templates/${promptId}`)
  }

  save(prompt: Instruction): Observable<Instruction> {
    if (isNew(prompt)) {
      return this.http.post<Instruction>(`/instruction-templates`, prompt)
    } else {
      return this.http.put<Instruction>(`/instruction-templates/${prompt.id}`, prompt)
    }
  }

  delete(promptId: number): Observable<void> {
    return this.http.delete<void>(`/instruction-templates/${promptId}`)
  }
}
