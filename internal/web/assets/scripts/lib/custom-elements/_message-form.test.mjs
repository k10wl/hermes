import * as assert from "node:assert";
import { after, before, beforeEach, describe, it, mock } from "node:test";

import { creator } from "./message-form.mjs";

describe("hermes-message-form", () => {
  const provider = {
    submit: mock.fn(() => new Promise(() => {})),
  };

  const subjectName = "hermes-message-form-test";
  /** @type {HTMLElement | null} */
  let submit = null;
  /** @type {HTMLElement & {value: string} | null} */
  let content = null;
  /** @type {HTMLElement | null} */
  let form = null;
  /** @type {HTMLInputElement | null} */
  let max_tokens = null;
  /** @type {HTMLInputElement | null} */
  let model = null;
  /** @type {HTMLInputElement | null} */
  let temperature = null;

  before(() => {
    customElements.define(subjectName, creator(provider));
  });
  beforeEach(() => {
    provider.submit.mock.resetCalls();
    document.body.innerHTML = `<${subjectName}></${subjectName}>`;
    const subject = document.querySelector(subjectName);
    if (!subject?.shadowRoot) {
      throw new Error(
        "test subject cannot be found, shadowRoot is null, please check element contents",
      );
    }
    submit = subject.shadowRoot.querySelector('button[type="submit"]');
    content = subject.shadowRoot.querySelector("#content");
    form = subject.shadowRoot.querySelector("form");
    max_tokens = subject.shadowRoot.querySelector("#max_tokens");
    model = subject.shadowRoot.querySelector("#model");
    temperature = subject.shadowRoot.querySelector("#temperature");
  });
  after(() => {
    document.body.innerHTML = "";
  });

  const expired = new Error("setup submit promise, but it never resolved");

  it("should submit message with arguments", async () => {
    const args = {
      content: "foo bar baz",
      params: {
        model: "openai/o3-mini",
        max_tokens: 1000,
        temperature: 0.5,
      },
    };
    if (content === null) {
      throw new Error("content element is null, please check structure");
    }
    content.value = args.content;
    if (max_tokens === null) {
      throw new Error("max_tokens element is null, please check structure");
    }
    max_tokens.value = args.params.max_tokens.toString();
    if (model === null) {
      throw new Error("model element is null, please check structure");
    }
    model.value = args.params.model;
    if (temperature === null) {
      throw new Error("temperature element is null, please check structure");
    }
    temperature.value = args.params.temperature.toString();
    const { promise, resolve, reject } = Promise.withResolvers();
    provider.submit.mock.mockImplementationOnce(async () => resolve(undefined));
    const id = setTimeout(() => reject(expired), 100);
    submit?.click();
    await promise;
    clearTimeout(id);
    assert.equal(
      provider.submit.mock.callCount(),
      1,
      "should call submit logic provider",
    );
    assert.deepEqual(provider.submit.mock.calls.at(0)?.arguments.at(0), args);
    assert.equal(content?.value, "");
  });

  it("should keep message content on error", async () => {
    const args = { content: "foo bar baz" };
    if (content === null) {
      throw new Error("content element is null, please check structure");
    }
    content.value = args.content;
    const { promise, resolve } = Promise.withResolvers();
    provider.submit.mock.mockImplementationOnce(() => {
      setTimeout(resolve);
      return Promise.reject();
    });
    submit?.click();
    await promise;
    assert.equal(
      content?.value,
      args.content,
      "content differs after rejection",
    );
  });

  it("should not submit empty messages", async () => {
    const { promise, resolve } = Promise.withResolvers();
    provider.submit.mock.mockImplementationOnce(Promise.resolve);
    const fn = mock.fn();
    form?.addEventListener("submit", fn);
    submit?.click();
    setTimeout(resolve, 100);
    await promise;
    assert.equal(fn.mock.callCount(), 1, "should still call submit event");
    assert.equal(
      provider.submit.mock.callCount(),
      0,
      "submit must not be called if content is empty",
    );
  });

  it("should submit on enter", async () => {
    const args = { content: "foo bar baz" };
    if (content === null) {
      throw new Error("content element is null, please check structure");
    }
    content.value = args.content;
    const { promise, resolve, reject } = Promise.withResolvers();
    provider.submit.mock.mockImplementationOnce(() => {
      setTimeout(resolve);
      return Promise.resolve();
    });
    form?.dispatchEvent(
      new KeyboardEvent("keydown", {
        key: "Enter",
        code: "Enter",
        bubbles: true,
        cancelable: true,
      }),
    );
    const id = setTimeout(() => reject(expired), 1000);
    await promise;
    clearTimeout(id);
    assert.equal(
      provider.submit.mock.callCount(),
      1,
      "should call submit logic provider",
    );
  });

  it("should not submit on shift-enter", async () => {
    const args = { content: "foo bar baz" };
    if (content === null) {
      throw new Error("content element is null, please check structure");
    }
    content.value = args.content;
    const { promise, resolve } = Promise.withResolvers();
    provider.submit.mock.mockImplementationOnce(Promise.resolve);
    const fn = mock.fn();
    form?.addEventListener("submit", fn);
    form?.dispatchEvent(
      new KeyboardEvent("keydown", {
        key: "Enter",
        code: "Enter",
        bubbles: true,
        cancelable: true,
        shiftKey: true,
      }),
    );
    setTimeout(resolve, 100);
    await promise;
    assert.equal(fn.mock.callCount(), 0, "shift-enter should not submit");
    assert.equal(
      provider.submit.mock.callCount(),
      0,
      "shift-enter should not submit",
    );
  });
});
