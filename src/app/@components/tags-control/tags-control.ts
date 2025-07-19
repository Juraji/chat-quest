import {
    Component,
    computed,
    effect,
    forwardRef,
    inject,
    input,
    InputSignal,
    linkedSignal,
    signal,
    Signal,
    WritableSignal
} from '@angular/core';
import {ControlValueAccessor, FormsModule, NG_VALUE_ACCESSOR} from '@angular/forms';
import {Tag, Tags} from '@db/tags';
import {toSignal} from '@angular/core/rxjs-interop';

@Component({
    selector: 'app-tags-control',
    templateUrl: './tags-control.html',
    styleUrls: ['./tags-control.scss'],
    imports: [
        FormsModule
    ],
    providers: [
        {
            provide: NG_VALUE_ACCESSOR,
            useExisting: forwardRef(() => TagsControl),
            multi: true
        }
    ]
})
export class TagsControl implements ControlValueAccessor {
    private readonly tags = inject(Tags)

    private onChange: (value: number[]) => void = () => null;
    private onTouched: () => void = () => null;

    readonly tagIds: InputSignal<number[]> = input<number[]>([])
    readonly disabled: InputSignal<boolean> = input<boolean>(false)

    readonly availableTags: Signal<Tag[]> = toSignal(this.tags.getAll(true), {initialValue: []})
    readonly currentTagIds: WritableSignal<number[]> = linkedSignal(() => this.tagIds())
    readonly currentTags: Signal<Tag[]> = computed(() => {
        const ids = this.currentTagIds()
        const tags = this.availableTags()
        return ids.map(id => tags
            .find(t => t.id === id))
            .filter(t => !!t)
    })

    readonly inputText: WritableSignal<string> = signal('')
    readonly isDisabled: WritableSignal<boolean> = linkedSignal(() => this.disabled())

    constructor() {
        effect(() => {
            // When inputText changes we are touched
            this.inputText()
            this.onTouched()
        });
    }

    writeValue(obj: number[]): void {
        if (!Array.isArray(obj)) {
            throw new Error('writeValue expects an array of tags');
        }

        this.currentTagIds.set(obj)
    }

    registerOnChange(fn: any): void {
        this.onChange = fn;
    }

    registerOnTouched(fn: any): void {
        this.onTouched = fn;
    }

    setDisabledState?(isDisabled: boolean): void {
        this.isDisabled.set(isDisabled);
    }

    addTag() {
        this.onTouched()
        const inputText = this.inputText().trim()
        if (!inputText) return

        const tagsToAdd = inputText
            .split(',')
            .map(t => t.trim())

        for (const tagName of tagsToAdd) {
            if (tagName) this.addSingleTag(tagName)
        }

        this.inputText.set('')
    }

    removeTag(tagId: number): void {
        this.onTouched()
        this.currentTagIds.update(tags => tags.filter(tId => tId !== tagId))
        this.onChange(this.currentTagIds())
    }

    onInputKeyDown(event: KeyboardEvent) {
        if (event.key === 'Enter') {
            this.addTag()
            event.stopPropagation()
            event.preventDefault()
        }
    }

    private addSingleTag(tagName: string) {
        const alreadyHasTag = this.currentTags().some(t => t.label.toLowerCase() === tagName.toLowerCase())
        if (alreadyHasTag) return

        const existingTag = this.availableTags().find(t => t.label.toLowerCase() === tagName.toLowerCase())
        if (existingTag) {
            // Use existing tag
            this.currentTagIds.update(tIds => ([...tIds, existingTag.id]))
            this.onChange(this.currentTagIds())
        } else {
            // Create new tag
            this.tags
                .save({label: tagName})
                .subscribe(tag => {
                    this.currentTagIds.update(tIds => ([...tIds, tag.id]))
                    this.onChange(this.currentTagIds())
                })
        }
    }
}
