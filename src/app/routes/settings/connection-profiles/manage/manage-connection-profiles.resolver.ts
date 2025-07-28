import {ResolveFn} from '@angular/router';
import {ConnectionProfile} from '@api/model';
import {inject} from '@angular/core';
import {ConnectionProfiles} from '@api/clients';

export const manageConnectionProfilesResolver: ResolveFn<ConnectionProfile[]> = () => {
  const service = inject(ConnectionProfiles)
  return service.getAll()
}
