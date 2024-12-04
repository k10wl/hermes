/**
 * @template {Function} T
 * @param {T} fn
 * @param {number} delay
 * @returns {(...args: Parameters<T>) => void}
 */
export function debounce(fn, delay) {
  /** @type {ReturnType<typeof setTimeout>} */
  let timeout;
  return function (...args) {
    clearTimeout(timeout);
    timeout = setTimeout(() => {
      fn(...args);
    }, delay);
  };
}
