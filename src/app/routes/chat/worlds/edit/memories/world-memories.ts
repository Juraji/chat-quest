import {Component, inject} from '@angular/core';
import {routeParamSignal} from '@util/ng';
import {ActivatedRoute} from '@angular/router';
import {MemoryList} from '@components/memory-list';

@Component({
  selector: 'world-memories',
  imports: [
    MemoryList
  ],
  templateUrl: './world-memories.html'
})
export class WorldMemories {
  private readonly activatedRoute = inject(ActivatedRoute)

  readonly worldId = routeParamSignal(this.activatedRoute, 'worldId', p => parseInt(p ?? '0'));
}
