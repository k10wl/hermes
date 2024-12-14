import { describe } from "node:test";

import * as assert from "assert";

import { ShortcutManager } from "./shortcut-manager.mjs";

describe("ShortcutManager", () => {
  assert.equal(
    ShortcutManager.eventToNotation(
      /** @type {KeyboardEvent} */ ({ code: "KeyA" }),
    ),
    "<KeyA>",
    "should parse key in notation",
  );

  assert.equal(
    ShortcutManager.eventToNotation(
      /** @type {KeyboardEvent} */ ({
        altKey: true,
        metaKey: true,
        code: "KeyA",
      }),
    ),
    "<M-KeyA>",
    "should parse compound alt+meta key",
  );

  assert.equal(
    ShortcutManager.eventToNotation(
      /** @type {KeyboardEvent} */ ({ altKey: true, code: "KeyA" }),
    ),
    "<M-KeyA>",
    "should parse compound meta key",
  );

  assert.equal(
    ShortcutManager.eventToNotation(
      /** @type {KeyboardEvent} */ ({ metaKey: true, code: "KeyA" }),
    ),
    "<M-KeyA>",
    "should parse compound meta key",
  );

  assert.equal(
    ShortcutManager.eventToNotation(
      /** @type {KeyboardEvent} */ ({ ctrlKey: true, code: "KeyA" }),
    ),
    "<C-KeyA>",
    "should parse compound ctrl key",
  );

  assert.equal(
    ShortcutManager.eventToNotation(
      /** @type {KeyboardEvent} */ ({ shiftKey: true, code: "KeyA" }),
    ),
    "<S-KeyA>",
    "should parse compound shift key",
  );

  assert.equal(
    ShortcutManager.eventToNotation(
      /** @type {KeyboardEvent} */ ({
        metaKey: true,
        altKey: true,
        ctrlKey: true,
        shiftKey: true,
        code: "KeyA",
      }),
    ),
    "<M-C-S-KeyA>",
    "should keep modifiers in specific sequence",
  );
});
