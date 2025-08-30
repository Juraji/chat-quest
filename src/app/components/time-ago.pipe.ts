import {inject, LOCALE_ID, Pipe, PipeTransform} from '@angular/core';
import {interval, map, Observable, of, startWith} from 'rxjs';
import {DATE_PIPE_DEFAULT_OPTIONS, DatePipe} from '@angular/common';

@Pipe({name: 'timeAgo'})
export class TimeAgoPipe implements PipeTransform {
  private readonly locale = inject(LOCALE_ID)
  private readonly dateOptions = inject(DATE_PIPE_DEFAULT_OPTIONS)
  private readonly datePipe = new DatePipe(this.locale, null, this.dateOptions)
  private readonly updateInterval = interval(5000).pipe(startWith(0));

  transform(date: Nullable<number | string | Date>, interactive: true): Observable<string>
  transform(date: Nullable<number | string | Date>, interactive: false): string
  transform(date: Nullable<number | string | Date>, interactive: boolean): Observable<string> | string {
    if (date === undefined || date === null) {
      return interactive ? of('') : ''
    }

    const baseDate = new Date(date)

    if (interactive) {
      return this.updateInterval.pipe(map(() => {
        const now = new Date();
        return this.format(baseDate, now);
      }))
    } else {
      return this.format(baseDate, new Date())
    }
  }

  private format(date: Date, now: Date): string {
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    if (seconds < 5) {
      return 'Just now';
    } else if (seconds < 60) {
      return `${seconds}s ago`;
    } else if (seconds < 3600) { // 1 hour
      const minutes = Math.floor(seconds / 60);
      return `${minutes}m ago`;
    } else if (seconds < 86400) { // 24 hours
      const hours = Math.floor(seconds / 3600);
      return `${hours}h ago`;
    } else {
      return this.datePipe.transform(date)!;
    }
  }
}
