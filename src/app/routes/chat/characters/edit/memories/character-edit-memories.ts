import {Component, inject, Signal} from '@angular/core';
import {MemoryList} from '@components/memory-list';
import {ActivatedRoute} from '@angular/router';
import {paramAsId, routeDataSignal, routeParamSignal} from '@util/ng';
import {World} from '@api/worlds';

@Component({
  selector: 'app-character-edit-memories',
  imports: [
    MemoryList
  ],
  templateUrl: './character-edit-memories.html'
})
export class CharacterEditMemories {
  protected readonly activatedRoute = inject(ActivatedRoute)

  readonly worlds: Signal<World[]> = routeDataSignal(this.activatedRoute, 'worlds')
  readonly characterId: Signal<number> = routeParamSignal(this.activatedRoute, 'characterId', paramAsId)
}
