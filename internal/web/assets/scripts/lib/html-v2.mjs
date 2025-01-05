import { AssertInstance, AssertNumber, AssertString } from "./assert.mjs";
import { Publisher } from "./publisher.mjs";

const tmpHolderElementTagName =
  "tmp-holder-element-0f00ba0c-4d56-4f59-854e-e0bef122f22d";
const tmpListenerPositionAttributeName =
  "tmp-param-listener-position-0f00ba0c-4d56-4f59-854e-e0bef122f22d";

/**
 * @param {string} holding
 * @param {number} at
 */
function createTmpHolderElementMarkdown(holding, at) {
  return `<${tmpHolderElementTagName} holding="${holding}" at="${at}"></${tmpHolderElementTagName}>`;
}

/**
 * @param {Element} element
 * @returns {string}
 * @throws {Error} if element has no holding data attribute
 */
function tmpHolding(element) {
  return AssertString.check(element.getAttribute("holding"));
}

class WithOffset {
  offset = 0;

  /**
   * @param {number} i
   * @returns {number}
   */
  adjust(i) {
    return this.offset + i;
  }

  /**
   * @param {number} i
   * @returns {number}
   */
  increment(i) {
    const update = this.offset + i;
    this.offset += 1;
    return update;
  }
}

/**
 * Transform template string into document structure on intuitive rules
 * This function is useful for creating dynamic HTML content with embedded expressions.
 *
 * @param {Parameters<typeof String.raw>} params
 * @returns {DocumentFragment} DOM representation of provided string
 *
 * @example
 * // Basic usage
 * const element = html`<div>Hello, ${name}!</div>`;
 *
 * @example
 * // Rendering multiple elements
 * const elements = html`<div>${item1}</div><div>${item2}</div>`;
 *
 * @example
 * // Using with a publisher. Update will replace inner html of parent
 * const publisher = new Publisher("Initial Value");
 * const elementWithPublisher = html`<div>${publisher}</div>`;
 *
 * @example
 * // Handling events
 * const button = html`<button onclick="${() => publisher.update(0)}">Click me</button>`;
 *
 * @example
 * // Binding element to expose node in callback. Useful to omit parsing
 * let buttonElement
 * html`<button bind="${(element) => (buttonElement = element)}"></button>`;
 */
export function html(...params) {
  const str = Array.from(params[0].raw);

  /** @type {Set<string>} */
  const listeners = new Set();

  const withOffset = new WithOffset();

  let hasTmpElements = false;
  let hasTmpDataAttributes = false;
  let hasBindings = false;

  params.forEach((param, i) => {
    if (i === 0) {
      return; // raw string array
    }

    try {
      AssertInstance.once(param, Publisher);
      str.splice(
        withOffset.increment(i),
        0,
        createTmpHolderElementMarkdown("publisher", i),
      );
      hasTmpElements = true;
      return;
    } catch {
      // just don't blow up
    }

    try {
      // absolutely crucial to handle before listeners
      if (typeof param !== "function") {
        throw null;
      }
      const index = withOffset.adjust(i);
      if (!str[index - 1]?.endsWith('bind="') || !str[index]?.startsWith('"')) {
        throw null;
      }
      str[index - 1] += `${i}`;
      hasBindings = true;
      return;
    } catch {
      // just don't blow up
    }

    try {
      if (typeof param !== "function") {
        throw null;
      }
      const index = withOffset.adjust(i);
      const prev = str[index - 1];
      const name = prev?.match(/on(?<name>\w+)="$/)?.groups?.name;
      if (!name || !prev || !str[index]?.startsWith('"')) {
        throw null;
      }
      str[index - 1] = prev.replace(
        /on\w+="$/,
        `${tmpListenerPositionAttributeName}-${name}="${i}"`,
      );
      hasTmpDataAttributes = true;
      listeners.add(name);
      str[index] = str[index].substring(1);
      return;
    } catch {
      // just don't blow up
    }

    try {
      AssertInstance.once(param, DocumentFragment);
      str.splice(
        withOffset.increment(i),
        0,
        createTmpHolderElementMarkdown("fragment", i),
      );
      hasTmpElements = true;
      return;
    } catch {
      // just don't blow up
    }

    if (Array.isArray(param)) {
      str.splice(
        withOffset.increment(i),
        0,
        createTmpHolderElementMarkdown("array", i),
      );
      hasTmpElements = true;
      return;
    }

    str.splice(withOffset.increment(i), 0, param);
  });

  let raw = str.join("");

  const fragment = document.createDocumentFragment();
  let dummy = document.createElement("div");
  fragment.append(dummy);
  dummy.outerHTML = raw;

  if (hasTmpDataAttributes) {
    listeners.forEach((name) => {
      const attribute = `${tmpListenerPositionAttributeName}-${name}`;
      fragment.querySelectorAll(`[${attribute}]`).forEach((element) => {
        element.addEventListener(
          name,
          params[
            AssertNumber.check(
              +(element.getAttribute(attribute) ?? NaN),
              "expected element to hold index value of function",
            )
          ],
        );
        element.removeAttribute(attribute);
      });
    });
  }

  if (hasTmpElements) {
    fragment.querySelectorAll(tmpHolderElementTagName).forEach((element) => {
      const holding = tmpHolding(element);

      switch (holding) {
        case "publisher": {
          const publisher = AssertInstance.once(
            params[AssertNumber.check(+(element.getAttribute("at") ?? NaN))],
            Publisher,
            "expected param to be a publisher",
          );

          const parent = AssertInstance.once(
            element.parentElement,
            HTMLElement,
            "parent element must exist as it's inner html will be replaced",
          );

          // meeeeeeeeh, this is leaky...
          Reflect.set(parent, "notify", (/** @type {unknown} */ value) => {
            parent.innerHTML = `${value}`;
          });

          publisher.subscribe(/** @type {any} it's fine */ (parent));

          parent.innerHTML = `${publisher.value}`;
          break;
        }

        case "fragment": {
          AssertInstance.once(
            element.parentElement,
            HTMLElement,
            "collection should have parent element to be rendered in",
          ).append(
            AssertInstance.once(
              params[AssertNumber.check(+(element.getAttribute("at") ?? NaN))],
              DocumentFragment,
              "expected collection replacement",
            ),
          );
          element.remove();
          break;
        }

        case "array": {
          const arrayFragment = new DocumentFragment();
          const arr =
            params[
              AssertNumber.check(
                +(element.getAttribute("at") ?? NaN),
                "expected at to point at params index",
              )
            ];
          arr.forEach((/** @type {unknown} */ el) => {
            try {
              arrayFragment.append(AssertInstance.once(el, DocumentFragment));
            } catch {
              arrayFragment.append(`${el}`);
            }
          });
          element.replaceWith(arrayFragment);
          break;
        }

        default:
          throw new Error(`enountered unhandled holding value - '${holding}'`);
      }
    });
  }

  if (hasBindings) {
    fragment.querySelectorAll("[bind]").forEach((element) => {
      const at = AssertNumber.check(
        +(element.getAttribute("bind") ?? NaN),
        "bind attribute should handle param index",
      );
      if (typeof params[at] !== "function") {
        throw new Error("binding is expected to be a function");
      }
      params[at](element);
      element.removeAttribute("bind");
    });
  }

  return fragment;
}
