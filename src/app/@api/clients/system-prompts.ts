import {Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {isNew, SystemPrompt} from '@api/model';

@Injectable({
  providedIn: 'root'
})
export class SystemPrompts {
  constructor(private http: HttpClient) {
  }

  getAll(): Observable<SystemPrompt[]> {
    return this.http.get<SystemPrompt[]>(`/system-prompts`)
  }

  get(promptId: number): Observable<SystemPrompt> {
    return this.http.get<SystemPrompt>(`/system-prompts/${promptId}`)
  }

  save(prompt: SystemPrompt): Observable<SystemPrompt> {
    if (isNew(prompt)) {
      return this.http.post<SystemPrompt>(`/system-prompts`, prompt)
    } else {
      return this.http.put<SystemPrompt>(`/system-prompts/${prompt.id}`, prompt)
    }
  }

  delete(promptId: number): Observable<void> {
    return this.http.delete<void>(`/system-prompts/${promptId}`)
  }
}
