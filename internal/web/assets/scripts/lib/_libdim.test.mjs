import { describe, it } from "node:test";

import * as assert from "assert";

import { Bind, html, Signal } from "./libdim.mjs";

// html function name allows prettier format. Cool feature, but messes up tests
const _html = html;

describe("html", () => {
  it("should return string", () => {
    const testsContainer = document.createElement("div");
    testsContainer.replaceChildren(_html`<div>test</div>`);
    assert.equal(
      testsContainer.innerHTML,
      "<div>test</div>",
      "expected to return string",
    );

    testsContainer.replaceChildren(_html`<div>first</div><div>second</div>`);
    assert.equal(
      testsContainer.innerHTML,
      "<div>first</div><div>second</div>",
      "expected to render multiple elements at same level",
    );

    testsContainer.replaceChildren(
      _html`<style>.div{color: red;}</style><div>second</div>`,
    );
    assert.equal(
      testsContainer.innerHTML,
      "<style>.div{color: red;}</style><div>second</div>",
      "expected to return style elemetn",
    );
  });

  it("common primitives with template input", () => {
    const testsContainer = document.createElement("div");
    testsContainer.replaceChildren(_html`<div>${"foo"}</div>`);
    assert.equal(
      testsContainer.innerHTML,
      "<div>foo</div>",
      "expected to insert string value",
    );

    testsContainer.replaceChildren(_html`<div>${null}</div>`);
    assert.equal(
      testsContainer.innerHTML,
      "<div></div>",
      "expected empty on null value",
    );

    testsContainer.replaceChildren(_html`<div>${42}</div>`);
    assert.equal(
      testsContainer.innerHTML,
      "<div>42</div>",
      "expected string interpritation of number on number",
    );

    testsContainer.replaceChildren(_html`<div>${true}</div>`);
    assert.equal(
      testsContainer.innerHTML,
      "<div>true</div>",
      "expected true to be in final",
    );

    testsContainer.replaceChildren(_html`<div>${false}</div>`);
    assert.equal(
      testsContainer.innerHTML,
      "<div>false</div>",
      "expected false to be in final",
    );

    testsContainer.replaceChildren(_html`foo`);
    assert.equal(
      testsContainer.innerHTML,
      "foo",
      "expected to render plain string",
    );
  });

  it("should be able to render nested collecections", () => {
    const el = document.createElement("div");
    document.body.append(el);
    el.append(
      _html`<div id="wrapper">${_html`<div>${_html`<span>nested</span>`}</div>`}</div>`,
    );

    assert.equal(
      el.innerHTML,
      `<div id="wrapper"><div><span>nested</span></div></div>`,
      "should render nested collections",
    );
  });

  it("should render maps", () => {
    const content = ["foo", "bar", "baz"];
    const el = document.createElement("div");
    document.body.append(el);
    el.append(
      _html`<div>${content.map((data) => _html`<span>${data}</span>`)}</div>`,
    );

    assert.equal(
      el.innerHTML,
      `<div><span>foo</span><span>bar</span><span>baz</span></div>`,
      "nested(_html calls should be correctly placed in DOM",
    );
  });

  it("should allow element binding", () => {
    const el = document.createElement("div");
    const div = new Bind();
    el.append(_html`<div id="target" bind="${div}"></div>`);
    const target = el.querySelector("#target");
    assert.equal(target !== null, true, "target should be rendered");
    assert.equal(
      div.current,
      target,
      "variable and document element should be same",
    );
  });

  it("should handle empty template input", () => {
    const testsContainer = document.createElement("div");
    testsContainer.replaceChildren(_html``);
    assert.equal(
      testsContainer.innerHTML,
      "",
      "expected to return empty string for empty template",
    );
  });

  it("should handle undefined values in template", () => {
    const testsContainer = document.createElement("div");
    testsContainer.replaceChildren(_html`<div>${undefined}</div>`);
    assert.equal(
      testsContainer.innerHTML,
      "<div></div>",
      "expected empty on undefined value",
    );
  });

  it("should handle nested templates with different types", () => {
    const el = document.createElement("div");
    document.body.append(el);
    el.append(_html`<div>${_html`<span>${"nested"}</span>`}</div>`);

    assert.equal(
      el.innerHTML,
      `<div><span>nested</span></div>`,
      "should render nested templates correctly",
    );
  });

  it("should remain rendering order in complex cases", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<div>1</div>`;
    const fragment3 = _html`<div>3</div>`;
    const fragment5 = _html`<div>5</div>`;

    const main = _html`<div>${fragment1}<div>2</div>${fragment3}<div>4</div>${fragment5}</div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><div>1</div><div>2</div><div>3</div><div>4</div><div>5</div></div>`,
      "should render nested templates correctly",
    );
  });

  it("should handle deeply nested structures", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<div><span>1</span></div>`;
    const fragment2 = _html`<div><span>2</span><span>2.1</span></div>`;
    const fragment3 = _html`<div><span>3</span><div><span>3.1</span></div></div>`;

    const main = _html`<div>${fragment1}${fragment2}${fragment3}</div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><div><span>1</span></div><div><span>2</span><span>2.1</span></div><div><span>3</span><div><span>3.1</span></div></div></div>`,
      "should render deeply nested templates correctly",
    );
  });

  it("should handle mixed content types", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<div>Text</div>`;
    const fragment2 = _html`<div><span>Element</span></div>`;
    const fragment3 = _html`<div>${42}</div>`;

    const main = _html`<div>${fragment1}${fragment2}${fragment3}</div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><div>Text</div><div><span>Element</span></div><div>42</div></div>`,
      "should render mixed content types correctly",
    );
  });

  it("should handle fragments without parent elements", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<span>Fragment 1</span>`;
    const fragment2 = _html`<span>Fragment 2</span>`;
    const fragment3 = _html`<span>Fragment 3</span>`;

    el.append(fragment1, fragment2, fragment3);

    assert.equal(
      el.innerHTML,
      `<span>Fragment 1</span><span>Fragment 2</span><span>Fragment 3</span>`,
      "should render fragments without parent elements correctly",
    );
  });

  it("should handle inline elements", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<span>Inline 1</span>`;
    const fragment2 = _html`<span>Inline 2</span>`;
    const fragment3 = _html`<span>Inline 3</span>`;

    el.append(_html`${fragment1}${fragment2}${fragment3}`);

    assert.equal(
      el.innerHTML,
      `<span>Inline 1</span><span>Inline 2</span><span>Inline 3</span>`,
      "should render inline elements correctly",
    );
  });

  it("should handle complex nested and inline structures", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<div><span>Nested 1</span></div>`;
    const fragment2 = _html`<span>Inline 2</span>`;
    const fragment3 = _html`<div><span>Nested 3</span><span>Inline 3.1</span></div>`;

    const main = _html`<div>${fragment1}${fragment2}${fragment3}</div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><div><span>Nested 1</span></div><span>Inline 2</span><div><span>Nested 3</span><span>Inline 3.1</span></div></div>`,
      "should render complex nested and inline structures correctly",
    );
  });

  it("should handle elements with attributes", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<div class="class1">1</div>`;
    const fragment2 = _html`<div id="id2">2</div>`;
    const fragment3 = _html`<div data-test="test3">3</div>`;

    const main = _html`<div>${fragment1}${fragment2}${fragment3}</div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><div class="class1">1</div><div id="id2">2</div><div data-test="test3">3</div></div>`,
      "should render elements with attributes correctly",
    );
  });

  it("should handle empty elements", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<div></div>`;
    const fragment2 = _html`<span></span>`;
    const fragment3 = _html`<p></p>`;

    const main = _html`<div>${fragment1}${fragment2}${fragment3}</div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><div></div><span></span><p></p></div>`,
      "should render empty elements correctly",
    );
  });

  it("should handle elements with text and children", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<div>Text<div>Child</div></div>`;
    const fragment2 = _html`<span>Text<span>Child</span></span>`;
    const fragment3 = _html`<div>Text<div>Child</div></div>`;

    const main = _html`<div>${fragment1}${fragment2}${fragment3}</div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><div>Text<div>Child</div></div><span>Text<span>Child</span></span><div>Text<div>Child</div></div></div>`,
      "should render elements with text and children correctly",
    );
  });

  it("should handle elements with mixed content and attributes", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<div class="class1">Text<div>Child</div></div>`;
    const fragment2 = _html`<span id="id2">Text<span>Child</span></span>`;
    const fragment3 = _html`<div data-test="test3">Text<div>Child</div></div>`;

    const main = _html`<div>${fragment1}${fragment2}${fragment3}</div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><div class="class1">Text<div>Child</div></div><span id="id2">Text<span>Child</span></span><div data-test="test3">Text<div>Child</div></div></div>`,
      "should render elements with mixed content and attributes correctly",
    );
  });

  it("should handle complex inline and block elements", () => {
    const el = document.createElement("div");

    const fragment1 = _html`<span>Inline 1</span>`;
    const fragment2 = _html`<div>Block 2</div>`;
    const fragment3 = _html`<span>Inline 3</span>`;

    const main = _html`<div>${fragment1}${fragment2}${fragment3}</div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><span>Inline 1</span><div>Block 2</div><span>Inline 3</span></div>`,
      "should render complex inline and block elements correctly",
    );
  });

  it("should handle rendering fragment within another fragment", () => {
    const el = document.createElement("div");

    const fragment3 = _html`<div>3</div>`;
    const fragment5 = _html`${fragment3}`;

    const main = _html`<div><div>1</div><div>2</div>${fragment5}<div>4</div><div>5</div></div>`;

    el.append(main);

    assert.equal(
      el.innerHTML,
      `<div><div>1</div><div>2</div><div>3</div><div>4</div><div>5</div></div>`,
      "should render fragment within another fragment correctly",
    );
  });
});

describe("Bind", () => {
  const tries = [
    null,
    undefined,
    "foo",
    42,
    true,
    false,
    new ArrayBuffer(0),
    new Blob(),
  ];

  it("should allow any bindings if no assertion is provided", () => {
    const binding = new Bind();
    tries.forEach((value) => {
      binding.current = value;
      assert.equal(binding.current, value, "should allow any bindings");
    });
  });

  it("should allow any bindings if assertion is provided", () => {
    const binding = new Bind((el) => {
      if (el instanceof ArrayBuffer) {
        return el;
      }
      throw new Error("expected to return HTMLElement");
    });
    const el = new ArrayBuffer(0);
    binding.current = el;
    assert.equal(
      binding.current,
      el,
      "should allow bindings that match assertion",
    );
  });

  it("should throw if binding does not match assertion", () => {
    const binding = new Bind((el) => {
      if (el instanceof Node) {
        return el;
      }
      throw new Error("expected to return HTMLElement");
    });
    tries.forEach((value) => {
      assert.throws(() => {
        // @ts-expect-error false positive when everything is fine
        binding.current = value;
      }, "should throw if binding does not match assertion");
    });
  });
});

describe("Signal", () => {
  it("should be able to subscribe to signal", () => {
    const signal = new Signal(0);
    const testContainer = document.createElement("div");
    const fragment = _html`<div>${signal}</div>`;
    testContainer.replaceChildren(fragment);
    assert.equal(testContainer.innerHTML, "<div>0</div>");
    signal.value = 1;
    assert.equal(testContainer.innerHTML, "<div>1</div>");
  });

  it("should be able to subscribe and unsubscribe from signal", () => {
    const signal = new Signal(0);
    const testContainer = document.createElement("div");
    const fragment = _html`<div>${signal}</div>`;
    testContainer.replaceChildren(fragment);
    assert.equal(testContainer.innerHTML, "<div>0</div>");
    let subscription = null;
    const teardown = signal.subscribe((value) => {
      subscription = value;
    });
    signal.value = 1;
    teardown();
    assert.equal(testContainer.innerHTML, "<div>1</div>");
    assert.equal(subscription, 1);
    signal.value = 2;
    assert.equal(testContainer.innerHTML, "<div>2</div>");
    assert.equal(subscription, 1);
  });

  it("should be able to change attributes", () => {
    const signal = new Signal(0);
    const fragment = _html`<div data-value="${signal}"></div>`;
    const holder = fragment.querySelector("[data-value]");
    if (!holder) {
      throw new Error("expected to find element");
    }
    assert.equal(holder.getAttribute("data-value"), "0");
    signal.value = 1;
    assert.equal(holder.getAttribute("data-value"), "1");
  });

  it("should render signals in the middle of contents", () => {
    const signal = new Signal(0);
    const fragment = _html`<div>${signal}foo${signal}</div>`;
    const testContainer = document.createElement("div");
    testContainer.replaceChildren(fragment);
    assert.equal(testContainer.innerHTML, "<div>0foo0</div>");
    signal.value = 1;
    assert.equal(testContainer.innerHTML, "<div>1foo1</div>");
  });

  it("should be able to signal into css", () => {
    const signal = new Signal("red");
    const fragment = _html`<style>.foo{color: ${signal}}</style><div class="foo">foo</div>`;
    const testContainer = document.createElement("div");
    testContainer.replaceChildren(fragment);
    assert.equal(
      testContainer.innerHTML,
      '<style>.foo{color: red}</style><div class="foo">foo</div>',
    );
    signal.value = "blue";
    assert.equal(
      testContainer.innerHTML,
      '<style>.foo{color: blue}</style><div class="foo">foo</div>',
    );
  });

  it("should remove falsy attributes to omit boolean attribute passing", () => {
    const signal = new Signal(false);
    const fragment = _html`<div disabled="${signal}">foo</div>`;
    const testContainer = document.createElement("div");
    testContainer.replaceChildren(fragment);
    assert.equal(testContainer.innerHTML, "<div>foo</div>");
    signal.value = true;
    assert.equal(testContainer.innerHTML, '<div disabled="true">foo</div>');
  });
});
