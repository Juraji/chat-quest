import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {map, Observable} from 'rxjs';

@Injectable({
  providedIn: 'root'
})
export class System {
  private http: HttpClient = inject(HttpClient)

  countTokens(text: string): Observable<number> {
    return this.http
      .post<{ count: number }>('/system/tokenizer/count', text)
      .pipe(map(res => res.count))
  }

  shutdown() {
    return this.http.post<void>('/system/shutdown', null)
  }
}
