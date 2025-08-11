import {Component, input, InputSignal, output, OutputEmitterRef, viewChild} from '@angular/core';
import {CropperOptions, ImageCropperComponent} from 'ngx-image-cropper';

@Component({
  selector: 'app-avatar-image-crop',
  imports: [ImageCropperComponent],
  templateUrl: './avatar-image-crop.html',
})
export class AvatarImageCrop {
  protected readonly imageCropper = viewChild(ImageCropperComponent);
  protected readonly cropperOptions: Partial<CropperOptions> = {
    autoCrop: false,
    maintainAspectRatio: true,
    format: 'jpeg',
    aspectRatio: 1,
    resizeToWidth: 500,
    onlyScaleDown: true,
    output: 'base64',
  };

  readonly imageDataUrl: InputSignal<string> = input.required()
  readonly onCropComplete: OutputEmitterRef<string> = output()
  readonly onCropCanceled: OutputEmitterRef<void> = output()

  onAcceptCrop() {
    const e = this.imageCropper()?.crop()
    if (!!e && !!e.base64) {
      this.onCropComplete.emit(e.base64)
    }
  }

  onCancel() {
    this.onCropCanceled.emit()
  }
}
