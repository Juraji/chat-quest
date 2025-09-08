import {Pipe, PipeTransform} from '@angular/core';
import {LlmModel, LlmModelView} from '@api/providers';

@Pipe({
  name: 'llmLabel'
})
export class LlmLabelPipe implements PipeTransform {

  transform(value: LlmModel | LlmModelView): string {
    let t: string
    switch (value.modelType) {
      case "CHAT_MODEL":
        t = 'Chat'
        break
      case "EMBEDDING_MODEL":
        t = 'Embedding'
        break
      default:
        t = 'Unknown'
    }

    if ('profileName' in value) {
      return `${value.modelId} \u2022 ${value.profileName} (${t})`
    } else {
      return `${value.modelId} (${t})`
    }
  }

}
