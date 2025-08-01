import {
  booleanAttribute,
  Component,
  computed,
  inject,
  input,
  InputSignal,
  InputSignalWithTransform,
  SecurityContext
} from '@angular/core';
import {DomSanitizer} from '@angular/platform-browser';

@Component({
  selector: 'app-rendered-message',
  imports: [],
  templateUrl: './rendered-message.html',
  styleUrl: './rendered-message.scss'
})
export class RenderedMessage {

  private readonly sanitizer = inject(DomSanitizer)

  readonly enableVariables: InputSignalWithTransform<boolean, unknown> = input(false, {transform: booleanAttribute});

  readonly message: InputSignal<string> = input.required()
  readonly formatted = computed(() => {
    const template = this.message()
    const enableVars = this.enableVariables()
    return this.render(template, enableVars) ?? ''
  })

  render(value: string | null, enableVars: boolean): string | null {
    if (!value) return value;
    let result = this.wrap(value, 'action', '*', '*');
    if (enableVars) {
      result = this.wrap(result, 'variable', '{{.', '}}');
    }
    return this.sanitizer.sanitize(SecurityContext.HTML, result)
  }

  private wrap(text: string, className: string, startSeq: string, endSeq: string): string {
    let nextStartPos = 0
    while (nextStartPos < text.length) {
      const startVarIdx = text.indexOf(startSeq, nextStartPos);
      if (startVarIdx === -1) break;

      let endVarIdx = startVarIdx + startSeq.length; // after startSeq
      while (endVarIdx < text.length && !text.startsWith(endSeq, endVarIdx)) {
        endVarIdx++;
      }
      if (endVarIdx >= text.length || !text.startsWith(endSeq, endVarIdx)) break;

      endVarIdx += endSeq.length // include endSeq
      const before = text.substring(0, startVarIdx)
      const after = text.substring(endVarIdx)
      const varContent = text.substring(startVarIdx, endVarIdx);
      const wrappedVar = `<span class="${className}">${varContent}</span>`;
      text = before + wrappedVar + after;
      nextStartPos = startVarIdx + wrappedVar.length // Including added chars
    }
    return text;
  }
}
