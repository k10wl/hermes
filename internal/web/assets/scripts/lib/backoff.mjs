/**
 * @param {number} initial
 * @param {(calls: number) => number} calculation
 */
export function backoff(initial, calculation) {
  let calls = 0;
  return () => {
    calls++;
    if (calls === 1) {
      return initial;
    }
    return initial * calculation(calls);
  };
}

/** @param {number} calls  */
export function exponent(calls) {
  return Math.pow(2, calls - 1);
}
