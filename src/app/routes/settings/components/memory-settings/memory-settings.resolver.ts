import {ResolveFn} from '@angular/router';
import {MemoryPreferences} from '@api/model';
import {inject} from '@angular/core';
import {Memories} from '@api/clients/memories';

export const memorySettingsResolver: ResolveFn<MemoryPreferences> = () => {
  const service = inject(Memories)
  return service.getPreferences();
};

