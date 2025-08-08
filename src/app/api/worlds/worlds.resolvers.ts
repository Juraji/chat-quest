import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Worlds} from '@api/worlds/worlds.service';
import {ChatPreferences} from '@api/worlds/index';

export const chatSettingsResolver: ResolveFn<ChatPreferences> = () => {
  const service = inject(Worlds)
  return service.getChatPreferences();
};
