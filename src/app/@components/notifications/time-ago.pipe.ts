import {Pipe, PipeTransform} from '@angular/core';
import {interval, map, Observable, startWith} from 'rxjs';

@Pipe({name: 'timeAgo'})
export class TimeAgoPipe implements PipeTransform {
  private updateInterval = interval(5000).pipe(startWith(0));

  transform(date: Date): Observable<string> {
    return this.updateInterval.pipe(map(() => {
      const now = new Date();
      return this.format(date, now);
    }))
  }

  private format(date: Date, now: Date): string {
    // Calculate the difference in seconds
    const seconds = Math.floor((now.getTime() - date.getTime()) / 1000);

    // Update every second using Angular's async pipe with RxJS interval
    // Note: This requires additional setup in your component

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
      // For older dates, use the full date format
      return date.toLocaleString();
    }
  }
}
