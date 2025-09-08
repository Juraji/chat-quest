import {inject, Injectable} from '@angular/core';
import {HttpClient} from '@angular/common/http';
import {Observable} from 'rxjs';
import {AiProviders, ConnectionProfile, LlmModel, LlmModelView} from './providers.model';
import {isNew} from '@api/common';

@Injectable({
  providedIn: 'root'
})
export class Providers {
  private http: HttpClient = inject(HttpClient)


  getTemplates(): Observable<AiProviders> {
    return this.http.get<AiProviders>(`${location.origin}/data/ai-providers.json`);
  }

  getAllLlmModelViews(): Observable<LlmModelView[]> {
    return this.http.get<LlmModelView[]>(`/connection-profiles/model-views`)
  }

  getAll(): Observable<ConnectionProfile[]> {
    return this.http.get<ConnectionProfile[]>(`/connection-profiles`)
  }

  get(profileId: number): Observable<ConnectionProfile> {
    return this.http.get<ConnectionProfile>(`/connection-profiles/${profileId}`)
  }

  save(profile: ConnectionProfile): Observable<ConnectionProfile> {
    if (isNew(profile)) {
      return this.http.post<ConnectionProfile>(`/connection-profiles`, profile)
    } else {
      return this.http.put<ConnectionProfile>(`/connection-profiles/${profile.id}`, profile)
    }
  }

  delete(profileId: number): Observable<void> {
    return this.http.delete<void>(`/connection-profiles/${profileId}`)
  }

  getModels(profileId: number): Observable<LlmModel[]> {
    return this.http.get<LlmModel[]>(`/connection-profiles/${profileId}/models`)
  }

  refreshModels(profileId: number): Observable<LlmModel[]> {
    return this.http.post<LlmModel[]>(`/connection-profiles/${profileId}/models/refresh`, null)
  }

  saveModel(llmModel: LlmModel): Observable<LlmModel> {
    return this.http.put<LlmModel>(`/connection-profiles/${llmModel.profileId}/models/${llmModel.id}`, llmModel)
  }
}
