import {Component, effect, inject, input, InputSignal, signal, WritableSignal} from '@angular/core';
import {Memories, Memory} from '@api/memories';
import {MemoryListItem} from './memory-list-item';
import {arrayRemove, arrayReplace} from '@util/array';
import {NEW_ID} from '@api/common';

@Component({
  selector: 'memory-list',
  imports: [
    MemoryListItem
  ],
  templateUrl: './memory-list.html',
})
export class MemoryList {
  private readonly memoriesStore = inject(Memories)

  readonly worldId: InputSignal<number> = input.required()
  readonly characterId: InputSignal<Nullable<number>> = input()
  protected readonly memories: WritableSignal<Memory[]> = signal([])
  protected readonly newMemories: WritableSignal<Memory[]> = signal([])

  constructor() {
    effect(() => {
      const worldId = this.worldId()
      const characterId = this.characterId()

      if (characterId === undefined) {
        this.memoriesStore
          .getAll(worldId)
          .subscribe(memories => this.memories.set(memories))
      } else if (!!characterId) {
        this.memoriesStore
          .getAllByCharacter(worldId, characterId)
          .subscribe(memories => this.memories.set(memories))
      }
    });
  }

  onAddMemory() {
    const newMemory: Memory = {
      id: NEW_ID,
      worldId: this.worldId(),
      characterId: null,
      createdAt: null,
      content: ''
    }

    this.newMemories.update(memories => [newMemory, ...memories])
  }

  onSaveMemory(memory: Memory, newMemoryIndex: Nullable<number>) {
    this.memoriesStore
      .save(this.worldId(), memory)
      .subscribe(updated => {
        this.memories.update(memories =>
          arrayReplace(memories, updated, m => m.id === memory.id));
        if (!!newMemoryIndex) {
          this.onDeleteMemory(newMemoryIndex)
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
      this.memoriesStore
        .delete(this.worldId(), id)
        .subscribe(() => this.memories.update(memories =>
          arrayRemove(memories, m => m.id === id)))
    }
  }
}
