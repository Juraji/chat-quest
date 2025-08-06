import {ResolveFn} from "@angular/router";
import {LlmModelView} from "@api/model";
import {inject} from "@angular/core";
import {ConnectionProfiles} from "@api/clients";

export const llmModelViewsResolver: ResolveFn<LlmModelView[]> = () => {
  const service = inject(ConnectionProfiles)
  return service.getAllLlmModelViews()
}
