import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {InstructionTemplate, isNew} from '@api/model';

@Injectable({
  providedIn: 'root'
})
export class InstructionTemplates {
  private http: HttpClient = inject(HttpClient)

  getAll(): Observable<InstructionTemplate[]> {
    return this.http.get<InstructionTemplate[]>(`/instruction-templates`)
  }

  get(promptId: number): Observable<InstructionTemplate> {
    return this.http.get<InstructionTemplate>(`/instruction-templates/${promptId}`)
  }

  save(prompt: InstructionTemplate): Observable<InstructionTemplate> {
    if (isNew(prompt)) {
      return this.http.post<InstructionTemplate>(`/instruction-templates`, prompt)
    } else {
      return this.http.put<InstructionTemplate>(`/instruction-templates/${prompt.id}`, prompt)
    }
  }

  delete(promptId: number): Observable<void> {
    return this.http.delete<void>(`/instruction-templates/${promptId}`)
  }
}
