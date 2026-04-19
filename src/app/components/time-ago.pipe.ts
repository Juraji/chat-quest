import {inject, LOCALE_ID, Pipe, PipeTransform} from '@angular/core';
import {interval, map, Observable, of, startWith} from 'rxjs';
import {DATE_PIPE_DEFAULT_OPTIONS, DatePipe} from '@angular/common';

const JUST_NOW_MAX = 5 // seconds
const SECONDS_MAX = 60 // 1 minute
const MINUTES_MAX = 3600 // 1 hour
const HOURS_MAX = 18000 // 5 hours
const UPDATE_INTERVAL = 5000

@Pipe({name: 'timeAgo'})
export class TimeAgoPipe implements PipeTransform {
  private readonly locale = inject(LOCALE_ID)
  private readonly dateOptions = inject(DATE_PIPE_DEFAULT_OPTIONS)
  private readonly datePipe = new DatePipe(this.locale, null, this.dateOptions)
  private readonly updateInterval = interval(UPDATE_INTERVAL).pipe(startWith(0));

  transform(date: Nullable<number | string | Date>, interactive: true): Observable<string>
  transform(date: Nullable<number | string | Date>, interactive: false): string
  transform(date: Nullable<number | string | Date>, interactive: boolean): Observable<string> | string {
    if (date === undefined || date === null) {
      return interactive ? of('') : ''
    }

    const baseDate = new Date(date)

    if (interactive) {
      const now = new Date()
      if (baseDate.getTime() < (new Date().getTime() - HOURS_MAX)) {
        return of(this.format(baseDate, now))
      }

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

    if (seconds < JUST_NOW_MAX) {
      return 'Just now';
    } else if (seconds < SECONDS_MAX) {
      return `${seconds}s ago`;
    } else if (seconds < MINUTES_MAX) { // 1 hour
      const minutes = Math.floor(seconds / 60);
      return `${minutes}m ago`;
    } else if (seconds < HOURS_MAX) { // 5 hours
      const hours = Math.floor(seconds / 3600);
      return `${hours}h ago`;
    } else {
      return this.datePipe.transform(date)!;
    }
  }
}
