import {Injectable} from '@angular/core';

@Injectable({
  providedIn: 'root'
})
export class SettingsStore {
  private static readonly PREFIX: string = 'rp_tavern:';

  get<T>(name: string): T | null {
    const value = localStorage.getItem(SettingsStore.PREFIX + name);
    return value ? JSON.parse(value) : null;
  }

  getOrElse<T>(name: string, defaultValue: T): T {
    return this.get(name) ?? defaultValue
  }

  set<T>(name: string, value: T): void {
    const ser = JSON.stringify(value)
    localStorage.setItem(SettingsStore.PREFIX + name, ser);
  }
}
