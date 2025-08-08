import {ResolveFn} from '@angular/router';
import {inject} from '@angular/core';
import {Memories} from '@api/memories/memories.service';
import {MemoryPreferences} from '@api/memories/memories.model';

export const memoryPreferencesResolver: ResolveFn<MemoryPreferences> = () => {
  const service = inject(Memories)
  return service.getPreferences();
};

