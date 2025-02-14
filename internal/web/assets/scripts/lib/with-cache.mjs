import { Queue } from "./queue.mjs";

/**
 * @template {Function} T
 * @param {T} fn
 * @param {{ cache?: Map<string, string>, maxSize?: number }} [options]
 * @returns {(...args: Parameters<T>) => ReturnType<T>}
 */
export function withCache(fn, options = { cache: new Map(), maxSize: 1_000 }) {
  const cache = options.cache ?? new Map();
  const maxSize = options.maxSize ?? 1_000;
  const q = new Queue();
  return (...args) => {
    const key = JSON.stringify(args);
    const cached = cache.get(key);
    if (cached !== undefined) {
      return cached;
    }
    const computed = fn(...args);
    cache.set(key, computed);
    q.enqueue(key);
    if (cache.size > maxSize) {
      cache.delete(q.dequeue());
    }
    return computed;
  };
}
