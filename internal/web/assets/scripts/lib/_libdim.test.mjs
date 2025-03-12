import { describe, test } from "node:test";

import * as assert from "assert";
import { readFileSync } from "fs";
import { JSDOM } from "jsdom";

import { Bind } from "./libdim.mjs";

/**
 * @param {string} name
 * @returns {string}
 */
function compileModularScripts(name) {
  const path = import.meta.dirname + name.replace(/^\./g, "");
  let text = readFileSync(path, {
    encoding: "utf8",
    flag: "r",
  });
  text = text.replaceAll(/^export /gm, "");

  const classMatcher = /^class (?<name>\w+)/gm;
  const classMatches = text.matchAll(classMatcher);
  classMatches.forEach((match) => {
    const name = match.groups?.name;
    if (typeof name !== "string") {
      throw new Error(
        `expected class member to have name ${JSON.stringify(match, null, 2)}`,
      );
    }
    text = text.replace(match[0], `window.${name} = ${match[0]}`);
  });

  return text;
}

const scriptContents = compileModularScripts("./libdim.mjs");

/**
 * @returns {JSDOM}
 */
function _prepareJSDOM() {
  const jsdom = new JSDOM(``, { runScripts: "dangerously" });
  const scriptElement = jsdom.window.document.createElement("script");
  scriptElement.textContent = scriptContents;
  jsdom.window.document.body.appendChild(scriptElement);
  return jsdom;
}

/**
 * @param {JSDOM} jsdom
 * @param {DocumentFragment} nodes
 */
function _asString(jsdom, nodes) {
  const el = jsdom.window.document.createElement("div");
  el.append(nodes);
  return el.innerHTML;
}

describe("html", () => {
  const jsdom = _prepareJSDOM();

  assert.equal(
    _asString(jsdom, jsdom.window.html`<div>test</div>`),
    "<div>test</div>",
    "expected to return string",
  );

  assert.equal(
    _asString(jsdom, jsdom.window.html`<div>first</div><div>second</div>`),
    "<div>first</div><div>second</div>",
    "expected to render multiple elements at same level",
  );

  assert.equal(
    _asString(
      jsdom,
      jsdom.window.html`<style>.div{color: red;}</style><div>second</div>`,
    ),
    "<style>.div{color: red;}</style><div>second</div>",
    "expected to return style elemetn",
  );

  test("common primitives with template input", () => {
    assert.equal(
      _asString(jsdom, jsdom.window.html`<div>${"foo"}</div>`),
      "<div>foo</div>",
      "expected to insert string value",
    );

    assert.equal(
      _asString(jsdom, jsdom.window.html`<div>${null}</div>`),
      "<div></div>",
      "expected empty on null value",
    );

    assert.equal(
      _asString(jsdom, jsdom.window.html`<div>${42}</div>`),
      "<div>42</div>",
      "expected string interpritation of number on number",
    );

    assert.equal(
      _asString(jsdom, jsdom.window.html`<div>${true}</div>`),
      "<div>true</div>",
      "expected true to be in final",
    );

    assert.equal(
      _asString(jsdom, jsdom.window.html`<div>${false}</div>`),
      "<div>false</div>",
      "expected false to be in final",
    );

    assert.equal(
      _asString(jsdom, jsdom.window.html`foo`),
      "foo",
      "expected to render plain string",
    );
  });

  test("should be able to render nested collecections", () => {
    const jsdom = _prepareJSDOM();
    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(
      jsdom.window
        .html`<div id="wrapper">${jsdom.window.html`<div>${jsdom.window.html`<span>nested</span>`}</div>`}</div>`,
    );

    assert.equal(
      el.innerHTML,
      `<div id="wrapper"><div><span>nested</span></div></div>`,
      "should render nested collections",
    );
  });

  test("should render maps", () => {
    const jsdom = _prepareJSDOM();
    const content = ["foo", "bar", "baz"];
    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(
      jsdom.window
        .html`<div>${content.map((data) => jsdom.window.html`<span>${data}</span>`)}</div>`,
    );

    assert.equal(
      el.innerHTML,
      `<div><span>foo</span><span>bar</span><span>baz</span></div>`,
      "nested html calls should be correctly placed in DOM",
    );
  });

  test("should allow element binding", () => {
    const jsdom = _prepareJSDOM();
    const el = jsdom.window.document.createElement("div");
    const div = new jsdom.window.Bind();
    el.append(jsdom.window.html`<div id="target" bind="${div}"></div>`);
    const target = el.querySelector("#target");
    assert.equal(target !== null, true, "target should be rendered");
    assert.equal(
      div.current,
      target,
      "variable and document element should be same",
    );
  });

  test("should handle empty template input", () => {
    assert.equal(
      _asString(jsdom, jsdom.window.html``),
      "",
      "expected to return empty string for empty template",
    );
  });

  test("should handle undefined values in template", () => {
    assert.equal(
      _asString(jsdom, jsdom.window.html`<div>${undefined}</div>`),
      "<div></div>",
      "expected empty on undefined value",
    );
  });

  test("should handle nested templates with different types", () => {
    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(
      jsdom.window
        .html`<div>${jsdom.window.html`<span>${"nested"}</span>`}</div>`,
    );

    assert.equal(
      el.innerHTML,
      `<div><span>nested</span></div>`,
      "should render nested templates correctly",
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

  test("should allow any bindings if no assertion is provided", () => {
    const binding = new Bind();
    tries.forEach((value) => {
      binding.current = value;
      assert.equal(binding.current, value, "should allow any bindings");
    });
  });

  test("should allow any bindings if assertion is provided", () => {
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

  test("should throw if binding does not match assertion", () => {
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
