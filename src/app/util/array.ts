/**
 * Adds a new item to an array and returns a new array with the added item.
 *
 * @template T - The type of elements in the array.
 * @param {T[]} items - The original array of items.
 * @param {T} itemToAdd - The item to be added to the array.
 * @param toStart - When true, adds the item to the start of the array, else adds it at the end
 * @returns {T[]} A new array containing all elements from the original array plus the new item.
 */
export function arrayAdd<T>(
  items: T[],
  itemToAdd: T,
  toStart: boolean = false,
): T[] {
  if (toStart) {
    return [itemToAdd, ...items];
  } else {
    return [...items, itemToAdd];
  }
}

/**
 * Adds or replaces an item in an array based on a condition.
 *
 * This function first attempts to find an existing item that matches the condition.
 * If found, it replaces that item with the new item. If not found,
 * it adds the new item to the end of the array.
 *
 * @template T - The type of elements in the array.
 * @param {T[]} items - The original array of items.
 * @param {T} newItem - The item to be added or used for updating.
 * @param {(item: T) => boolean} where - A predicate function to find the item to update.
 * @returns {T[]} A new array with either the specified item updated or the new item added.
 */
export function arrayReplace<T>(
  items: T[],
  newItem: T,
  where: (item: T) => boolean
): T[] {
  const index = items.findIndex(where)
  if (index === -1) return arrayAdd(items, newItem)
  return [...items.slice(0, index), newItem, ...items.slice(index + 1)];
}

/**
 * Updates an item in an array based on a condition and returns a new array.
 *
 * This function finds an existing item that matches the condition, applies an update function to it,
 * and returns a new array with the updated item. If no matching item is found, it throws an error.
 *
 * @template T - The type of elements in the array.
 * @param {T[]} items - The original array of items.
 * @param {(item: T) => T} update - A function that takes an item and returns its updated version.
 * @param {(item: T) => boolean} where - A predicate function to find the item to update.
 * @returns {T[]} A new array with the specified item updated.
 * @throws {Error} Throws an error if no matching item is found for updating.
 */
export function arrayUpdate<T>(
  items: T[],
  update: (item: T) => T,
  where: (item: T) => boolean
): T[] {
  const i = items.findIndex(where)
  if (i === -1) throw new Error("Update for unknown entity")
  const u: T = update(items[i])
  return [...items.slice(0, i), u, ...items.slice(i + 1)]
}

/**
 * Removes an item from an array and returns a new array without the removed item.
 *
 * @template T - The type of elements in the array.
 * @param {T[]} items - The original array of items.
 * @param {((item: T) => boolean) | number} where - A predicate function to find the item to remove or the index to remove.
 * @returns {T[]} A new array with the specified item removed, or the original array if no match was found.
 */
export function arrayRemove<T>(
  items: T[],
  where: ((item: T) => boolean) | number
): T[] {
  const index = typeof where === 'number' ? where : items.findIndex(where);
  if (index === -1) return items;
  return [...items.slice(0, index), ...items.slice(index + 1)];
}

/**
 * Merges two arrays into one, ensuring no duplicate items exist in the result.
 *
 * The merged array will contain all unique elements from both input arrays,
 * with preference given to items from the first array when duplicates are found.
 * Equality is determined using a custom equality function provided by the caller.
 *
 * @template T - The type of elements in the arrays.
 * @param {T[]} itemsA - The first array of items.
 * @param {T[]} itemsB - The second array of items to merge with the first.
 * @param {(a: T, b: T) => boolean} equalityFn - A function that determines if two items are equal.
 * @returns {T[]} A new array containing all unique elements from both input arrays,
 *                with preference given to items from the first array when duplicates exist.
 */
export function arrayMerge<T>(
  itemsA: T[],
  itemsB: T[],
  equalityFn: (a: T, b: T) => boolean
): T[] {
  return [...itemsA, ...itemsB.filter(itemB =>
    !itemsA.some(itemA => equalityFn(itemA, itemB))
  )];
}
