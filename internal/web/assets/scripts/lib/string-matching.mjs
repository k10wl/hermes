import { AssertString } from "./assert.mjs";

/**
 * @param {string} source
 * @param {string} lookup
 * @returns {{matches: boolean[], ok: boolean}}
 * @throws {Error} on bad argument types
 */
export function stringMatching(source, lookup) {
  AssertString.check(source, "source must be a string");
  AssertString.check(lookup, "lookup must be a string");
  if (lookup.length === 0) {
    return { ok: true, matches: source.split("").map(() => false) };
  }
  /** @type {boolean[]} */
  const matches = [];
  let p = 0;
  for (let i = 0; i < source.length; i++) {
    let match = false;
    if (p >= lookup.length) {
      matches[i] = match;
    }
    const char = source[i];
    if (!char) {
      break;
    }
    if (char === lookup[p]) {
      p++;
      match = true;
    }
    matches[i] = match;
  }
  return { ok: p === lookup.length, matches };
}
