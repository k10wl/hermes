/**
 * libdim - Library for DOM Interactions Management
 * Utilizes template literals to bind event handlers and DOM elements
 *
 * @example
 * // Create a button
 * const name = "my-button";
 * const button = new Bind();
 * html`<button bind="${button}" onclick="${() => alert("clicked!")}">${name}</button>`;
 * console.log("> this is my button", button.target);
 */

// attribute holds index of event listener
const EVENT_ATTRIBUTE_PREFIX = "data-dim-processing-event";
// attribute holds index of nested element
const NESTED_INDEX_ATTRIBUTE = "data-dim-processing-nested-index";
// attribute holds index of binding class
const BIND_ATTRIBUTE = "bind";

/**
 * At preparation step, creates a placeholder for nested element
 * Placeholder will be replaced with actual element during fragment creation
 * @param {number} index
 */
function createNestedElementHolder(index) {
  return `<slot ${NESTED_INDEX_ATTRIBUTE}="${index}"></slot>`;
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
 * @param {Parameters<typeof String.raw>} params
 * @returns {{
 *   raw: string,
 *   listeners: Set<string>,
 *   hasBindings: boolean,
 *   hasEventsListeners: boolean,
 *   hasNestedFragments: boolean,
 * }}
 */
function prepareData(params) {
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
          `${EVENT_ATTRIBUTE_PREFIX}-${name}="${i}"`,
        );
        hasEventsListeners = true;
        listeners.add(name);
        str[index] = str[index].substring(1);
        return;
      }
    }

    if (param instanceof Bind) {
      const index = withOffset.adjust(i);
      if (
        str[index - 1]?.endsWith(`${BIND_ATTRIBUTE}="`) &&
        str[index]?.startsWith('"')
      ) {
        str[index - 1] += `${i}`;
        hasBindings = true;
        return;
      }
    }

    if (param instanceof DocumentFragment) {
      str.splice(withOffset.increment(i), 0, createNestedElementHolder(i));
      hasNestedFragments = true;
      return;
    }

    // could be buggy with 'interesting' types, but it's not worth the effort
    // expect that consumer will use 'todoArr.map((el) => html`<div>${el}</div>`)'
    if (Array.isArray(param)) {
      str.splice(withOffset.increment(i), 0, createNestedElementHolder(i));
      hasNestedFragments = true;
      return;
    }

    str.splice(withOffset.increment(i), 0, param);
  });

  let raw = str.join("");
  return {
    raw,
    listeners,
    hasBindings,
    hasEventsListeners,
    hasNestedFragments,
  };
}

/**
 * @param {{
 *   raw: string,
 *   listeners: Set<string>,
 *   hasBindings: boolean,
 *   hasEventsListeners: boolean,
 *   hasNestedFragments: boolean,
 * }} data
 * @param {Parameters<typeof String.raw>} params
 * @returns {DocumentFragment}
 */
function createFragment(data, params) {
  const {
    raw,
    listeners,
    hasBindings,
    hasEventsListeners,
    hasNestedFragments,
  } = data;

  const fragment = document.createDocumentFragment();
  let dummy = document.createElement("div");
  dummy.innerHTML = raw;
  fragment.append(...dummy.childNodes); // chrome does not work without dummy

  let selectors = "";
  if (hasEventsListeners) {
    listeners.forEach((name) => {
      selectors += `[${EVENT_ATTRIBUTE_PREFIX}-${name}],`;
    });
  }
  if (hasNestedFragments) {
    selectors += `[${NESTED_INDEX_ATTRIBUTE}],`;
  }
  if (hasBindings) {
    selectors += `[${BIND_ATTRIBUTE}],`;
  }

  if (selectors) {
    selectors = selectors.slice(0, -1);
    fragment.querySelectorAll(selectors).forEach((element) => {
      if (hasEventsListeners) {
        const eventListenerNames = Array.from(listeners).filter((name) =>
          element.matches(`[${EVENT_ATTRIBUTE_PREFIX}-${name}]`),
        );

        for (const name of eventListenerNames) {
          const attribute = `${EVENT_ATTRIBUTE_PREFIX}-${name}`;
          element.addEventListener(
            name,
            params[getAttributeIndex(element, attribute)],
          );
          element.removeAttribute(attribute);
        }
      }

      if (
        hasNestedFragments &&
        element.matches(`[${NESTED_INDEX_ATTRIBUTE}]`) &&
        element instanceof HTMLSlotElement
      ) {
        /** @type {DocumentFragment | unknown[]} */
        const fragment =
          params[getAttributeIndex(element, NESTED_INDEX_ATTRIBUTE)];
        if (Array.isArray(fragment)) {
          element.replaceWith(
            ...fragment.map((el) => (el instanceof Node ? el : `${el}`)),
          );
          return;
        }
        element.replaceWith(fragment);
      }

      if (hasBindings && element.hasAttribute(`${BIND_ATTRIBUTE}`)) {
        const binding = params[getAttributeIndex(element, `${BIND_ATTRIBUTE}`)];
        if (!(binding instanceof Bind)) {
          console.error(binding);
          throw new Error("binding is expected to be a function");
        }
        binding.current = element;
        element.removeAttribute(`${BIND_ATTRIBUTE}`);
      }
    });
  }
  return fragment;
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
 * const fragment = html`<button onclick="${() => alert("clicked!")}">Click me</button>`;
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
function html(...params) {
  const data = prepareData(params);
  return createFragment(data, params);
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
class Bind {
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

export { Bind, html };
