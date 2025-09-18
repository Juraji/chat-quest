import {Component, computed, input, InputSignal, InputSignalWithTransform} from '@angular/core';

type RendererOptions = {
  enableActions: boolean
  enableQuoted: boolean
  enableVariables: boolean
  enableOOC: boolean
  enableMD: boolean
}

const DEFAULT_OPTIONS: RendererOptions = {
  enableActions: true,
  enableQuoted: true,
  enableVariables: true,
  enableOOC: true,
  enableMD: false,
}

@Component({
  selector: 'app-rendered-message',
  imports: [],
  templateUrl: './rendered-message.html',
  styleUrl: './rendered-message.scss'
})
export class RenderedMessage {

  readonly renderOptions: InputSignalWithTransform<RendererOptions, Partial<RendererOptions>> =
    input(DEFAULT_OPTIONS, {transform: v => ({...DEFAULT_OPTIONS, ...v})})

  readonly message: InputSignal<string> = input.required()

  readonly formatted = computed(() => {
    const template = this.message()
    const opts = this.renderOptions()
    if (!template) return template;
    return this.render(template, opts)
  })

  render(value: string, opts: RendererOptions): string {
    let result = this.escapeHtml(value)

    if (opts.enableActions) {
      result = this.wrap(result, 'action', '*', '*');
    }
    if (opts.enableQuoted) {
      result = this.wrap(result, 'quoted', '&quot;', '&quot;');
      result = this.wrap(result, 'quoted', '&ldquo;', '&rdquo;');
    }
    if (opts.enableVariables) {
      result = this.wrap(result, 'variable-block', '{{r', '}}');
      result = this.wrap(result, 'variable-block', '{{i', '}}');
      result = this.wrap(result, 'variable-block', '{{e', '}}');
      result = this.wrap(result, 'variable', '{{$', '}}');
      result = this.wrap(result, 'variable', '{{.', '}}');
    }
    if (opts.enableOOC) {
      result = this.wrap(result, 'out-of-character ', '[OOC:', ']');
      result = this.wrap(result, 'out-of-character ', '[System note:', ']');
    }
    if (opts.enableMD) {
      result = this.wrap(result, 'md-block', '```', '```');
    }

    return result
  }

  private wrap(
    text: string,
    className: string,
    startSeq: string,
    endSeq: string,
    trimSeq: boolean = false,
    allowMissingEndSeq: boolean = false,
  ): string {
    let nextStartPos = 0
    while (nextStartPos < text.length) {
      const startVarIdx = text.indexOf(startSeq, nextStartPos);
      if (startVarIdx === -1) break;

      let endVarIdx = startVarIdx + startSeq.length; // after startSeq
      while (endVarIdx < text.length && !text.startsWith(endSeq, endVarIdx)) {
        endVarIdx++;
      }
      if (endVarIdx >= text.length || !text.startsWith(endSeq, endVarIdx)) {
        if (allowMissingEndSeq) {
          const before = text.substring(0, startVarIdx)
          const varContent = trimSeq
            ? text.substring(startVarIdx + startSeq.length).trim()
            : text.substring(startVarIdx);
          const wrappedVar = `<span class="${className}">${varContent}</span>`;
          text = before + wrappedVar
        }

        break;
      }

      endVarIdx += endSeq.length // include endSeq
      const before = text.substring(0, startVarIdx)
      const after = text.substring(endVarIdx)
      const varContent = trimSeq
        ? text.substring(startVarIdx + startSeq.length, endVarIdx - endSeq.length).trim()
        : text.substring(startVarIdx, endVarIdx);
      const wrappedVar = `<span class="${className}">${varContent}</span>`;
      text = before + wrappedVar + after;
      nextStartPos = startVarIdx + wrappedVar.length // Including added chars
    }
    return text;
  }

  private escapeHtml(unsafe: string): string {
    return unsafe
      .replaceAll('&', '&amp;')
      .replaceAll('<', '&lt;')
      .replaceAll('>', '&gt;')
      .replaceAll('"', '&quot;')
      .replaceAll('“', '&ldquo;')
      .replaceAll('”', '&rdquo;')
      .replaceAll("'", '&#039;');
  };
}
