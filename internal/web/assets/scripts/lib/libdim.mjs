/**
 * libdim - Library for DOM Interactions Management
 * Utilizes template literals to bind event handlers and DOM elements
 *
 * @example
 * // Create a button
 * const name = "my-button";
 * const button = new Bind();
 * fucntion onClick() {
 *   alert("clicked!");
 * }
 * html`<button bind="${button}" onclick=${onClick}>${name}</button>`;
 * console.log("> this is my button", button.target);
 */

const eventAttributePrefix = "data-dim-event";
const nestedTypeAttribute = "data-dim-nested-type";
const indexAttribute = "data-dim-nested-index";

/**
 * @param {string} type
 * @param {number} index
 */
function createNestedElementHolder(type, index) {
  return `<div ${nestedTypeAttribute}="${type}" ${indexAttribute}="${index}"></div>`;
}

/**
 * @param {Element} element
 * @returns {string}
 * @throws {Error} if element has no holding data attribute
 */
function getTypeAttribute(element) {
  const holding = element.getAttribute(nestedTypeAttribute);
  if (typeof holding !== "string") {
    console.error(element);
    throw new Error(
      `expected element to have holding attribute: ${JSON.stringify(holding, null, 2)}`,
    );
  }
  return holding;
}

/**
 * @param {Element} element
 * @param {string} attribute
 * @returns {number}
 * @throws {Error} if element has no at data attribute
 */
function getAttributeIndex(element, attribute) {
  const index = +(element.getAttribute(attribute) ?? NaN);
  if (Number.isNaN(index)) {
    console.error(element);
    throw new Error(
      `expected "${attribute}" to be a number: ${JSON.stringify(index, null, 2)}`,
    );
  }
  return index;
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
 * // Adding event listeners
 * const fragment = html`<button onclick="${() => alert('clicked!')}">Click me</button>`;
 *
 * @example
 * // Binding elements to variables
 * const input = new Bind();
 * const fragment = html`
 *   <form>
 *     <input bind="${input}" type="text">
 *     <button onclick="${() => console.log(input.value)}">Submit</button>
 *   </form>
 * `;
 * input.current.value = "Hello, world!";
 */
export function html(...params) {
  const str = Array.from(params[0].raw);

  /** @type {Set<string>} */
  const listeners = new Set();

  const withOffset = new WithOffset();

  let hasBindings = false;
  let hasEventsListeners = false;
  let hasNestedFragments = false;

  params.forEach((param, i) => {
    if (i === 0) {
      return; // first element is always string
    }

    if (typeof param === "function") {
      const index = withOffset.adjust(i);
      const prev = str[index - 1];
      const name = prev?.match(/on(?<name>\w+)="$/)?.groups?.name;
      if (name && prev && str[index]?.startsWith('"')) {
        str[index - 1] = prev.replace(
          /on\w+="$/,
          `${eventAttributePrefix}-${name}="${i}"`,
        );
        hasEventsListeners = true;
        listeners.add(name);
        str[index] = str[index].substring(1);
        return;
      }
    }

    if (param instanceof Bind) {
      const index = withOffset.adjust(i);
      if (str[index - 1]?.endsWith('bind="') && str[index]?.startsWith('"')) {
        str[index - 1] += `${i}`;
        hasBindings = true;
        return;
      }
    }

    if (param instanceof DocumentFragment) {
      str.splice(
        withOffset.increment(i),
        0,
        createNestedElementHolder("fragment", i),
      );
      hasNestedFragments = true;
      return;
    }

    if (Array.isArray(param)) {
      str.splice(
        withOffset.increment(i),
        0,
        createNestedElementHolder("array", i),
      );
      hasNestedFragments = true;
      return;
    }

    str.splice(withOffset.increment(i), 0, param);
  });

  let raw = str.join("");

  const fragment = document.createDocumentFragment();
  let dummy = document.createElement("div");
  dummy.innerHTML = raw;
  fragment.append(...dummy.childNodes); // chrome does not work without dummy

  let selectors = "";
  if (hasEventsListeners) {
    listeners.forEach((name) => {
      selectors += `[${eventAttributePrefix}-${name}],`;
    });
  }
  if (hasNestedFragments) {
    selectors += `[${nestedTypeAttribute}],`;
  }
  if (hasBindings) {
    selectors += "[bind],";
  }

  if (selectors) {
    selectors = selectors.slice(0, -1);
    fragment.querySelectorAll(selectors).forEach((element) => {
      if (hasEventsListeners) {
        const eventListenerNames = Array.from(listeners).filter((name) =>
          element.matches(`[${eventAttributePrefix}-${name}]`),
        );

        for (const name of eventListenerNames) {
          const attribute = `${eventAttributePrefix}-${name}`;
          element.addEventListener(
            name,
            params[getAttributeIndex(element, attribute)],
          );
          element.removeAttribute(attribute);
        }
      }

      if (hasNestedFragments && element.matches(`[${nestedTypeAttribute}]`)) {
        const typeAttribute = getTypeAttribute(element);

        switch (typeAttribute) {
          case "fragment": {
            if (!(element.parentElement instanceof HTMLElement)) {
              throw new Error(
                "fragments must have parent element to be rendered in",
              );
            }
            const fragment = params[getAttributeIndex(element, indexAttribute)];
            if (!(fragment instanceof DocumentFragment)) {
              throw new Error("expected fragment replacement");
            }
            element.parentElement.append(fragment);
            element.remove();
            break;
          }

          case "array": {
            const arrayFragment = new DocumentFragment();
            const arr = params[getAttributeIndex(element, indexAttribute)];
            arr.forEach((/** @type {unknown} */ el) => {
              if (el instanceof DocumentFragment) {
                arrayFragment.append(el);
              } else {
                arrayFragment.append(`${el}`);
              }
            });
            element.replaceWith(arrayFragment);
            break;
          }

          default:
            throw new Error(
              `encountered unhandled holding value - '${typeAttribute}'`,
            );
        }
      }

      if (hasBindings && element.hasAttribute("bind")) {
        const binding = params[getAttributeIndex(element, "bind")];
        if (!(binding instanceof Bind)) {
          throw new Error("binding is expected to be a function");
        }
        binding.current = element;
        element.removeAttribute("bind");
      }
    });
  }

  return fragment;
}

/**
 * Allows to bind created elements to a value
 * Accepts assertion function to check if the value is valid on read and write
 * Can be omitted, but types will fallback to `unknown`
 *
 * @example
 * const dropdown = new Bind(assertButton)
 * const fragment = html`<button bind="${dropdown}">Click me</button>`;
 * console.log(dropdown.current.textContent) // "Click me"
 *
 * @template [T=unknown]
 */
export class Bind {
  /** @typedef {(value: unknown) => T} Assertion */

  /** @type {T | null} */
  #current = null;

  /** @type {Assertion | undefined} */
  #assertion;

  /** @param {Assertion} [assertion] */
  constructor(assertion) {
    this.#assertion = assertion;
  }

  /** @returns {this['#assertion'] extends undefined ? (T | null) : T} */
  get current() {
    if (this.#assertion) {
      return this.#assertion(this.#current);
    }
    return /** @type {any} consumer will receive `unknown` */ (this.#current);
  }

  set current(current) {
    let _current = current;
    if (this.#assertion) {
      _current = this.#assertion(current);
    }
    this.#current = _current;
  }
}
