import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {CQPreferences} from './preferences.model';

@Injectable({providedIn: 'root'})
export class Preferences {
  private http: HttpClient = inject(HttpClient)

  get(): Observable<CQPreferences> {
    return this.http.get<CQPreferences>('/preferences')
  }

  save(prefs: CQPreferences): Observable<CQPreferences> {
    return this.http.put<CQPreferences>(`/preferences`, prefs)
  }
}
