import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Worlds} from '@api/clients/worlds';
import {ChatPreferences} from '@api/model';

export const chatSettingsResolver: ResolveFn<ChatPreferences> = () => {
  const service = inject(Worlds)
  return service.getChatPreferences();
};
