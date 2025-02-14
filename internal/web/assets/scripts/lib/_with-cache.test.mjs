import { describe, test } from "node:test";

import * as assert from "assert";

import { withCache } from "./with-cache.mjs";

describe("withCache", () => {
  test("with primitive", () => {
    /**
     * @param {string} str
     * @returns {string}
     */
    const fn = (str) => {
      return `computed: ${str}`;
    };
    /** @type {Map<string, string>} */
    const cache = new Map([['["cached"]', "cached"]]);
    const fnWithCache = withCache(fn, { cache });
    assert.equal(fnWithCache("cached"), "cached", "should return cached value");
    assert.equal(
      fnWithCache("new"),
      "computed: new",
      "should compute new value",
    );
    assert.equal(
      cache.get('["new"]'),
      "computed: new",
      "should store computed value in cache",
    );
  });

  test("with multiple primitives", () => {
    /**
     * @param {string} a
     * @param {string} b
     * @returns {string}
     */
    const fn = (a, b) => {
      return `computed: ${a} ${b}`;
    };
    const cache = new Map([['["a","b"]', "cached"]]);
    const fnWithCache = withCache(fn, { cache });
    assert.equal(fnWithCache("a", "b"), "cached", "should return cached value");
    assert.equal(
      fnWithCache("new a", "new b"),
      "computed: new a new b",
      "should compute new value",
    );
    assert.equal(
      cache.get('["new a","new b"]'),
      "computed: new a new b",
      "should store computed value in cache",
    );
  });

  test("should limit size of cache", () => {
    /**
     * @param {string} str
     * @returns {string}
     */
    const fn = (str) => {
      return `computed: ${str}`;
    };
    /** @type {Map<string, string>} */
    const cache = new Map();
    const maxSize = 2;
    const fnWithCache = withCache(fn, { cache, maxSize });
    for (let i = 0; i < maxSize + 1; i++) {
      fnWithCache(`${i}`);
    }
    assert.equal(cache.size, maxSize, "should limit cache size");
    assert.equal(
      JSON.stringify(cache.entries().toArray()),
      JSON.stringify([
        ['["1"]', "computed: 1"],
        ['["2"]', "computed: 2"],
      ]),
      "should work as FIFO",
    );
  });
});
