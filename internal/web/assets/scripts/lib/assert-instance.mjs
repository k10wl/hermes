/**
 * @template T
 * @returns {T}
 * @param {unknown} obj
 * @param {new (data: any) => T} type
 */
export function assertInstance(obj, type) {
  if (obj instanceof type) {
    /** @type {any} */
    const any = obj;
    /** @type {T} */
    const t = any;
    return t;
  }
  throw new Error(`Object ${obj} does not have the right type '${type}'!`);
}
