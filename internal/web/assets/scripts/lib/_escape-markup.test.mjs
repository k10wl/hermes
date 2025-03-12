import * as assert from "node:assert";
import test, { describe } from "node:test";

import { escapeMarkup } from "./escape-markup.mjs";

describe("escapeMarkup", () => {
  test("should escape markup", () => {
    assert.equal(
      escapeMarkup("<div>test&</div>"),
      "&lt;div&gt;test&amp;&lt;/div&gt;",
    );
  });

  test("should not escape markup", () => {
    assert.equal(
      escapeMarkup("&lt;div&gt;test&lt;/div&gt;"),
      "&amp;lt;div&amp;gt;test&amp;lt;/div&amp;gt;",
    );
  });
});
