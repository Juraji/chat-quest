import {AfterViewInit, Component, ElementRef, inject, Signal, signal, viewChild, WritableSignal} from '@angular/core';
import {debounceTime, fromEvent} from "rxjs";
import {takeUntilDestroyed} from "@angular/core/rxjs-interop";

@Component({
    selector: 'scalable',
    imports: [],
    templateUrl: './scalable.html',
    styleUrl: './scalable.scss',
    host: {
        '[style.width]': 'hostWidth()',
        '[style.height]': 'hostHeight()',
    }
})
export class Scalable implements AfterViewInit {
    private readonly hostRef: ElementRef<HTMLElement> = inject(ElementRef)
    private readonly contentRef: Signal<ElementRef<HTMLElement> | undefined> =
        viewChild("content", {read: ElementRef})

    protected readonly contentScale: WritableSignal<string> = signal('scale(0)')
    protected readonly hostWidth: WritableSignal<string | null> = signal(null)
    protected readonly hostHeight: WritableSignal<string | null> = signal(null)

    constructor() {
        fromEvent(window, 'resize')
            .pipe(
                takeUntilDestroyed(),
                debounceTime(100)
            )
            .subscribe(() => this.updateScale())
    }

    ngAfterViewInit() {
        this.updateScale()
    }

    private updateScale() {
        const wrapper = this.hostRef.nativeElement;
        const content = this.contentRef()?.nativeElement;

        if (!wrapper || !content) return;

        this.hostWidth.set(null)
        this.hostHeight.set(null)

        requestAnimationFrame(() => {
            const wrapperWidth = wrapper.clientWidth;
            const wrapperHeight = wrapper.clientHeight;
            const contentWidth = content.scrollWidth;
            const contentHeight = content.scrollHeight;

            // Calculate the scale factor to fit the content within the wrapper
            const widthScale = wrapperWidth / contentWidth;
            const heightScale = wrapperHeight / contentHeight;

            const scale = Math.min(widthScale, heightScale)

            const scaledWidth = contentWidth * scale;
            const scaledHeight = contentHeight * scale;

            // Set parent to exactly fit the scaled content + 1px
            this.contentScale.set(`scale(${scale})`);
            this.hostWidth.set(`${scaledWidth}px`)
            this.hostHeight.set(`${scaledHeight}px`)
        })
    }
}
