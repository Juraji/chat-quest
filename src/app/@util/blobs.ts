export function readBlobAsText(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onerror = () => reject(reader.error)
    reader.onload = () => resolve(reader.result as string)
    reader.readAsText(blob)
  })
}

export async function readBlobAsJson<T>(blob: Blob): Promise<T> {
  if (blob.type !== 'application/json') throw new Error('Blob type "application/json" is required')
  const data = await readBlobAsText(blob);
  return JSON.parse(data);
}

export function readBlobAsDataUrl(blob: Blob): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader()
    reader.onerror = () => reject(reader.error)
    reader.onload = () => resolve(reader.result as string)
    reader.readAsDataURL(blob)
  })
}

export function downloadBlob(blob: Blob, fileName: string) {
  const url = URL.createObjectURL(blob);
  const body = document.body

  const anchor = document.createElement('a');
  anchor.href = url;
  anchor.download = fileName;

  body.appendChild(anchor);
  anchor.click()
  body.removeChild(anchor);

  URL.revokeObjectURL(url);
}

export function dataUrlToBlob(dataUrl: string): Blob {
  const contentType = dataUrl.match(/^data:(.*?);base64,/)![1]
  const base64Data = dataUrl.split(',')[1]

  const byteCharacters = atob(base64Data);
  return slicedBlobOf(byteCharacters, contentType)
}

export function slicedBlobOf(data: string, contentType: string, sliceSize: number = 1024) {
  const bytesLength = data.length;
  const slicesCount = Math.ceil(bytesLength / sliceSize);
  const byteArrays = new Array(slicesCount);

  for (let sliceIndex = 0; sliceIndex < slicesCount; ++sliceIndex) {
    const begin = sliceIndex * sliceSize;
    const end = Math.min(begin + sliceSize, bytesLength);

    const bytes = new Array(end - begin);
    for (let offset = begin, i = 0; offset < end; ++i, ++offset) {
      bytes[i] = data[offset].charCodeAt(0);
    }
    byteArrays[sliceIndex] = new Uint8Array(bytes);
  }
  return new Blob(byteArrays, {type: contentType});
}
