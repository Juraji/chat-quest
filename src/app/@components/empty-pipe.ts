import {Pipe, PipeTransform} from '@angular/core';

@Pipe({
  name: 'empty',
})
export class EmptyPipe implements PipeTransform {
  transform(value: any[] | null | undefined): boolean {
    return Array.isArray(value) ? value.length === 0 : true;
  }
}
