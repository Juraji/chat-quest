import {booleanAttribute, Component, effect, inject, input, InputSignal, signal, WritableSignal} from '@angular/core';
import {Memories, Memory} from '@api/memories';
import {MemoryListItem} from './memory-list-item';
import {arrayRemove, arrayReplace} from '@util/array';
import {NEW_ID} from '@api/common';
import {defer, map} from 'rxjs';

@Component({
  selector: 'memory-list',
  imports: [
    MemoryListItem
  ],
  templateUrl: './memory-list.html',
})
export class MemoryList {
  private readonly memoriesService = inject(Memories)

  readonly worldId: InputSignal<number> = input.required()
  readonly characterId: InputSignal<Nullable<number>> = input()
  readonly disabled = input(false, {transform: booleanAttribute})

  protected readonly memories: WritableSignal<Memory[]> = signal([])
  protected readonly newMemories: WritableSignal<Memory[]> = signal([])

  constructor() {
    effect(() => {
      const worldId = this.worldId()
      const characterId = this.characterId()

      const memories$ = defer(() => {
        if (characterId === undefined) {
          return this.memoriesService.getAll(worldId)
        } else if (!!characterId) {
          return this.memoriesService.getAllByCharacter(worldId, characterId)
        } else {
          return []
        }
      })

      memories$
        .pipe(map(memories => memories
          .sort((a, b) => b.createdAt!.localeCompare(a.createdAt!))))
        .subscribe(memories => this.memories.set(memories))
    });
  }

  onAddMemory() {
    const newMemory: Memory = {
      id: NEW_ID,
      worldId: this.worldId(),
      characterId: this.characterId(),
      createdAt: null,
      content: ''
    }

    this.newMemories.update(memories => [newMemory, ...memories])
  }

  onSaveMemory(memory: Memory, newMemoryIndex: Nullable<number>) {
    this.memoriesService
      .save(this.worldId(), memory)
      .subscribe(updated => {
        this.memories.update(memories =>
          arrayReplace(memories, updated, m => m.id === memory.id));
        if (typeof newMemoryIndex === 'number') {
          this.onDeleteNewMemory(newMemoryIndex)
        }
      })
  }

  onDeleteNewMemory(atIndex: number) {
    this.newMemories.update(memories =>
      arrayRemove(memories, atIndex))
  }

  onDeleteMemory(id: number) {
    const doDelete = confirm('Are you sure you want to delete this memory?');

    if (doDelete) {
      this.memoriesService
        .delete(this.worldId(), id)
        .subscribe(() => this.memories.update(memories =>
          arrayRemove(memories, m => m.id === id)))
    }
  }
}
