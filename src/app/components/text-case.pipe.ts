import {Pipe, PipeTransform} from '@angular/core';

type CaseType = "TITLE" | "LOWER" | "UPPER"

@Pipe({name: 'textCase'})
export class TextCasePipe implements PipeTransform {
  transform(value: Nullable<string>, caseType: CaseType = "TITLE"): string {
    if (!value) return ""

    return value
      .split("_")
      .map(word => this.applyCase(word, caseType))
      .join(" ")
  }

  private applyCase(word: string, caseType: CaseType) {
    switch (caseType) {
      case "TITLE":
        return word[0].toUpperCase() + word.slice(1).toLocaleLowerCase()
      case "LOWER":
        return word.toLocaleLowerCase()
      case "UPPER":
        return word.toUpperCase()
      default:
        throw new Error("Unknown case type: " + caseType)
    }
  }
}
