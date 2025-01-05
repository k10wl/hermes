import { describe, test } from "node:test";

import * as assert from "assert";
import { readFileSync } from "fs";
import { JSDOM } from "jsdom";

import { AssertInstance, AssertString } from "./assert.mjs";

/**
 * @param {string} name
 * @returns {string}
 */
function compileModularScripts(name) {
  const path = import.meta.dirname + name.replace(/^\./g, "");
  const importMatcher = /import .*? from "(?<path>.*?)";/gms;
  let text = readFileSync(path, {
    encoding: "utf8",
    flag: "r",
  });
  text = text.replaceAll(/^export /gm, "");

  const matches = text.matchAll(importMatcher);
  matches.forEach((match) => {
    text = text.replace(
      match[0],

      compileModularScripts(AssertString.check(match.groups?.path)),
    );
  });

  const classMatcher = /^class (?<name>\w+)/gm;
  const classMatches = text.matchAll(classMatcher);
  classMatches.forEach((match) => {
    text = text.replace(
      match[0],
      `window.${AssertString.check(match.groups?.name)} = ${match[0]}`,
    );
  });

  return text;
}

const scriptContents = compileModularScripts("./html-v2.mjs");

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

  test("publisher as template param", () => {
    const publisher = new jsdom.window.Publisher("foo");

    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(jsdom.window.html`<div>${publisher}</div>`);

    assert.equal(
      el.innerHTML,
      "<div>foo</div>",
      "expected to insert string value",
    );

    publisher.update("bar");
    assert.equal(
      el.innerHTML,
      "<div>bar</div>",
      "expected to update value on publisher update",
    );
  });

  test("event as template param wiht publisher update as plain value", () => {
    const jsdom = _prepareJSDOM();
    let counter = new jsdom.window.Publisher(0);

    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(
      jsdom.window.html`<button onclick="${() =>
        counter.update(1)}">${counter}</button>`,
    );

    const button = AssertInstance.once(
      el.querySelector("button"),
      jsdom.window.HTMLButtonElement,
      "expected button as result of html",
    );
    button.dispatchEvent(new jsdom.window.Event("click"));

    assert.equal(
      button.innerHTML,
      "1",
      "expected counter to increment after click",
    );
  });

  test("event as template param wiht publisher update as callback", () => {
    const jsdom = _prepareJSDOM();
    let counter = new jsdom.window.Publisher(0);

    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(
      jsdom.window.html`<button onclick="${() =>
        counter.update(
          (/** @type {number} */ prev) => prev + 1,
        )}">${counter}</button>`,
    );

    const button = AssertInstance.once(
      el.querySelector("button"),
      jsdom.window.HTMLButtonElement,
      "expected button as result of html",
    );
    button.dispatchEvent(new jsdom.window.Event("click"));

    assert.equal(
      button.innerHTML,
      "1",
      "expected counter to increment after click",
    );
  });

  test("should support multiple events on single elemnt", () => {
    const jsdom = _prepareJSDOM();
    let clickCounter = new jsdom.window.Publisher(0);
    let pointerOverCounter = new jsdom.window.Publisher(0);

    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(
      jsdom.window.html`
        <button 
            onclick="${() => clickCounter.update(1)}"
            onpointerover="${() => pointerOverCounter.update(1)}"
        ></button>
        `,
    );

    const button = AssertInstance.once(
      el.querySelector("button"),
      jsdom.window.HTMLButtonElement,
      "expected button as result of html",
    );

    button.dispatchEvent(new jsdom.window.Event("click"));
    button.dispatchEvent(new jsdom.window.Event("pointerover"));

    assert.equal(
      clickCounter.value,
      1,
      "expected click counter to increment after click",
    );
    assert.equal(
      pointerOverCounter.value,
      1,
      "expected copies counter to increment after click",
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

  test("should render maps of publishers and track listeners", () => {
    const jsdom = _prepareJSDOM();
    const e0 = new jsdom.window.Publisher(0);
    const e1 = new jsdom.window.Publisher(0);
    const e2 = new jsdom.window.Publisher(0);

    const el = jsdom.window.document.createElement("div");

    el.append(
      jsdom.window.html`${[e0, e1, e2].map(
        (publisher) =>
          jsdom.window
            .html`<button onclick="${() => e0.update(1)}" onpointerover="${() => e1.update(1)}" onmousemove="${() => e2.update(1)}">${publisher}</button>`,
      )}`,
    );

    const buttons = el.querySelectorAll("button");

    assert.equal(
      el.innerHTML,
      `<button>0</button><button>0</button><button>0</button>`,
      "nested html calls should be correctly placed in DOM",
    );

    buttons.item(0).dispatchEvent(new jsdom.window.Event("click"));
    assert.equal(
      el.innerHTML,
      `<button>1</button><button>0</button><button>0</button>`,
      "nested html calls should be correctly placed in DOM",
    );

    buttons.item(1).dispatchEvent(new jsdom.window.Event("pointerover"));
    assert.equal(
      el.innerHTML,
      `<button>1</button><button>1</button><button>0</button>`,
      "nested html calls should be correctly placed in DOM",
    );

    buttons.item(2).dispatchEvent(new jsdom.window.Event("mousemove"));
    assert.equal(
      el.innerHTML,
      `<button>1</button><button>1</button><button>1</button>`,
      "nested html calls should be correctly placed in DOM",
    );
  });

  test("should be able to render events within attributes", () => {
    const jsdom = _prepareJSDOM();
    const el = jsdom.window.document.createElement("div");
    const id = 0;

    const publisher = new jsdom.window.Publisher(0);

    el.append(
      jsdom.window
        .html`<button class="${id}" onclick="${() => publisher.update((/** @type {number} */ prev) => prev + 1)}" id="${id}">${publisher}</button>`,
    );

    assert.equal(
      el.innerHTML,
      `<button class="0" id="0">0</button>`,
      "listeners should be correctly assigned between template params",
    );

    const button = el.querySelector("button");
    button?.dispatchEvent(new jsdom.window.Event("click"));

    assert.equal(
      el.innerHTML,
      `<button class="0" id="0">1</button>`,
      "listeners should still be attached",
    );
  });

  test("should allow element binding", () => {
    const jsdom = _prepareJSDOM();
    const el = jsdom.window.document.createElement("div");

    let bind;

    el.append(
      jsdom.window
        .html`<div id="target" bind="${(/** @type {unknown} */ el) => (bind = el)}"></div>`,
    );

    const target = el.querySelector("#target");

    assert.equal(target !== null, true, "target should be rendered");

    assert.equal(bind, target, "variable and document element should be same");
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

  test("should handle multiple updates to publishers", () => {
    const publisher = new jsdom.window.Publisher(0);
    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(jsdom.window.html`<div>${publisher}</div>`);

    publisher.update(1);
    assert.equal(
      el.innerHTML,
      "<div>1</div>",
      "expected to update value on first publisher update",
    );

    publisher.update(2);
    assert.equal(
      el.innerHTML,
      "<div>2</div>",
      "expected to update value on second publisher update",
    );
  });

  test("should handle multiple elements with the same publisher", () => {
    const publisher = new jsdom.window.Publisher(0);
    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(
      jsdom.window.html`<div>${publisher}</div><div>${publisher}</div>`,
    );

    publisher.update(1);
    assert.equal(
      el.innerHTML,
      "<div>1</div><div>1</div>",
      "expected both elements to update with the same publisher value",
    );
  });

  test("should handle event listeners with different contexts", () => {
    const jsdom = _prepareJSDOM();
    let counter1 = new jsdom.window.Publisher(0);
    let counter2 = new jsdom.window.Publisher(0);

    const el = jsdom.window.document.createElement("div");
    jsdom.window.document.body.append(el);
    el.append(
      jsdom.window.html`
<button onclick="${() => counter1.update(1)}">Counter 1: <span>${counter1}</span></button>
<button onclick="${() => counter2.update(1)}">Counter 2: <span>${counter2}</span></button>
        `,
    );

    const buttons = el.querySelectorAll("button");
    buttons.item(0).dispatchEvent(new jsdom.window.Event("click"));
    assert.equal(
      buttons.item(0).innerHTML,
      "Counter 1: <span>1</span>",
      "expected counter 1 to increment after click",
    );

    buttons.item(1).dispatchEvent(new jsdom.window.Event("click"));
    assert.equal(
      buttons.item(1).innerHTML,
      "Counter 2: <span>1</span>",
      "expected counter 2 to increment after click",
    );
  });

  test("should render html publisher as dom contents", () => {
    const jsdom = _prepareJSDOM();
    const html = jsdom.window.html;
    const publisher = new jsdom.window.Publisher(html`<span>foo</span>`);

    const el = jsdom.window.document.createElement("div");
    el.append(html`<div>${publisher}</div>`);

    assert.equal(
      el.innerHTML,
      "<div><span>foo</span></div>",
      "expected html publisher to be interprited as dom nodes",
    );

    const click = new jsdom.window.Publisher(0);
    publisher.update(
      html`<button onclick="${() => click.update(1)}">${click}</button>`,
    );

    assert.equal(
      el.innerHTML,
      "<div><button>0</button></div>",
      "expected inner html to be updated on publisher update",
    );

    el.querySelector("button")?.dispatchEvent(new jsdom.window.Event("click"));
    assert.equal(
      el.innerHTML,
      "<div><button>1</button></div>",
      "expected inner publisher to update its value",
    );
  });
});
