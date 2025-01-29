import test, { describe } from "node:test";

import * as assert from "assert";

import { stringMatching } from "./string-matching.mjs";

describe("stringMatching", () => {
  test("should match on words", () => {
    assert.deepStrictEqual(
      stringMatching("hello", "hello"),
      { ok: true, matches: [true, true, true, true, true] },
      "full word match",
    );
    assert.deepStrictEqual(
      stringMatching("hello", "hell"),
      { ok: true, matches: [true, true, true, true, false] },
      "hell match",
    );
    assert.deepStrictEqual(
      stringMatching("hello", "h"),
      { ok: true, matches: [true, false, false, false, false] },
      "h match",
    );
    //assert.deepStrictEqual(
    //  stringMatching("hello", "lo"),
    //  { ok: true, matches: [false, false, false, true, true] },
    //  "lo match",
    //);
    //assert.deepStrictEqual(
    //  stringMatching("hello", "elo"),
    //  { ok: true, matches: [false, true, false, true, true] },
    //  "elo match",
    //);
    //assert.deepStrictEqual(
    //  stringMatching("hello", "helo"),
    //  { ok: true, matches: [true, true, true, false, true] },
    //  "helo match",
    //);
    //assert.deepStrictEqual(
    //  stringMatching("hello", "hlo"),
    //  { ok: true, matches: [true, false, false, true, true] },
    //  "hlo match",
    //);
    assert.deepStrictEqual(
      stringMatching("hello", "miss"),
      { ok: false, matches: [false, false, false, false, false] },
      "full miss",
    );
    assert.deepStrictEqual(
      stringMatching("hello", "wow"),
      { ok: false, matches: [false, false, false, false, false] },
      "full miss but with matching substr",
    );
    assert.deepStrictEqual(
      stringMatching("hello", ""),
      { ok: true, matches: [false, false, false, false, false] },
      "match on empty search",
    );
    //assert.deepStrictEqual(
    //  stringMatching("chats history", "history"),
    //  {
    //    ok: true,
    //    matches: [
    //      false,
    //      false,
    //      false,
    //      false,
    //      false,
    //      false,
    //      true,
    //      true,
    //      true,
    //      true,
    //      true,
    //      true,
    //      true,
    //    ],
    //  },
    //  "should highlight",
    //);
  });
});
