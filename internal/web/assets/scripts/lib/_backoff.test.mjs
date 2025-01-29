import * as assert from "node:assert";
import { describe } from "node:test";

import { backoff, exponent } from "./backoff.mjs";

describe("backoff test", () => {
  const next = backoff(1, (n) => n++);
  assert.equal(next(), 1, "did not return inital value on first call");
  assert.equal(next(), 2, "incremented value by one to 2");
  assert.equal(next(), 3, "incremented value by one to 3");
});

describe("exponent test", () => {
  const next = backoff(1, exponent);
  assert.equal(next(), 1, "did not return inital value on first call");
  assert.equal(next(), 2, "incremented value by one to 2");
  assert.equal(next(), 4, "incremented value by one to 4");
  assert.equal(next(), 8, "incremented value by one to 8");
  assert.equal(next(), 16, "incremented value by one to 16");
});
