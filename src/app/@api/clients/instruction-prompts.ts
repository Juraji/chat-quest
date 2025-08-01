import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {InstructionPrompt, isNew} from '@api/model';

@Injectable({
  providedIn: 'root'
})
export class InstructionPrompts {
  private http: HttpClient = inject(HttpClient)

  getAll(): Observable<InstructionPrompt[]> {
    return this.http.get<InstructionPrompt[]>(`/instruction-prompts`)
  }

  get(promptId: number): Observable<InstructionPrompt> {
    return this.http.get<InstructionPrompt>(`/instruction-prompts/${promptId}`)
  }

  save(prompt: InstructionPrompt): Observable<InstructionPrompt> {
    if (isNew(prompt)) {
      return this.http.post<InstructionPrompt>(`/instruction-prompts`, prompt)
    } else {
      return this.http.put<InstructionPrompt>(`/instruction-prompts/${prompt.id}`, prompt)
    }
  }

  delete(promptId: number): Observable<void> {
    return this.http.delete<void>(`/instruction-prompts/${promptId}`)
  }
}
