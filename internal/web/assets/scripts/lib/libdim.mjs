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
const PARAMETER_INDEX_ATTRIBUTE = "data-dim-processing-parameter-index";
const PARAMETER_TYPE_ATTRIBUTE = "data-dim-processing-parameter-type";
// attribute holds index of binding class
const BIND_ATTRIBUTE = "bind";

const NESTED_TYPE = {
  FRAGMENT: "fragment",
  SIGNAL_NODE: "signal-node",
  SIGNAL_ATTRIBUTE: "signal-attribute",
};
const SIGNAL_ATTRIBUTE_DELIMITER = "::";

/**
 * At preparation step, creates a placeholder for nested element
 * Placeholder will be replaced with actual element during fragment creation
 * @param {number} index
 * @param {string} type
 */
function createNestedElementHolder(type, index) {
  return `<slot ${PARAMETER_INDEX_ATTRIBUTE}="${index}" ${PARAMETER_TYPE_ATTRIBUTE}="${type}"></slot>`;
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
    console.trace();
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
 *   hasProcessingParameters: boolean,
 * }}
 */
function prepareData(params) {
  const str = Array.from(params[0].raw);

  /** @type {Set<string>} */
  const listeners = new Set();

  const withOffset = new WithOffset();

  let hasBindings = false;
  let hasEventsListeners = false;
  let hasProcessingParameters = false;

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

    if (param instanceof Signal) {
      const prev = str[i - 1];
      if (prev) {
        const attributeRegex = /(?<name>\w+(-\w+)?)="$/;
        const exec = attributeRegex.exec(prev);
        if (exec?.groups?.name) {
          str[i - 1] +=
            `${i}" ${PARAMETER_TYPE_ATTRIBUTE}="${NESTED_TYPE.SIGNAL_ATTRIBUTE}${SIGNAL_ATTRIBUTE_DELIMITER}${exec.groups.name}`;
          hasProcessingParameters = true;
          return;
        }
      }

      str.splice(
        withOffset.increment(i),
        0,
        createNestedElementHolder(NESTED_TYPE.SIGNAL_NODE, i),
      );
      hasProcessingParameters = true;
      return;
    }

    if (param instanceof DocumentFragment) {
      str.splice(
        withOffset.increment(i),
        0,
        createNestedElementHolder(NESTED_TYPE.FRAGMENT, i),
      );
      hasProcessingParameters = true;
      return;
    }

    // could be buggy with 'interesting' types, but it's not worth the effort
    // expect that consumer will use 'todoArr.map((el) => html`<div>${el}</div>`)'
    if (Array.isArray(param)) {
      str.splice(
        withOffset.increment(i),
        0,
        createNestedElementHolder(NESTED_TYPE.FRAGMENT, i),
      );
      hasProcessingParameters = true;
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
    hasProcessingParameters,
  };
}

/**
 * @param {Object} args
 * @param {Signal<unknown>} args.signal
 * @param {WeakRef<Element>} args.element
 * @param {string} args.attribute
 */
function subscribeSignalAttribute(args) {
  const { signal, element, attribute } = args;
  signal.subscribe(function attributeChange(value) {
    const elRef = element.deref();
    if (!elRef) {
      return;
    }
    elRef.setAttribute(attribute, `${value}`);
  });
}

/**
 * @param {Object} args
 * @param {Signal<unknown>} args.signal
 * @param {WeakRef<Text>} args.mark
 * @param {WeakRef<Text>} args.bound
 */
function subscribeSignalNode(args) {
  const { mark, bound, signal } = args;
  signal.subscribe(function nodeChange(value) {
    const markRef = mark.deref();
    const boundRef = bound.deref();
    if (!markRef || !boundRef) {
      signal.unsubscribe(nodeChange);
      return;
    }
    while (markRef.nextSibling !== boundRef && markRef.nextSibling) {
      markRef.nextSibling.remove();
    }
    if (!markRef.nextSibling) {
      markRef.after(boundRef);
    }
    if (Array.isArray(value)) {
      markRef.after(...value.map((el) => (el instanceof Node ? el : `${el}`)));
    }
    markRef.after(value instanceof Node ? value : `${value}`);
  });
}

/**
 * @param {{
 *   raw: string,
 *   listeners: Set<string>,
 *   hasBindings: boolean,
 *   hasEventsListeners: boolean,
 *   hasProcessingParameters: boolean,
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
    hasProcessingParameters,
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
  if (hasProcessingParameters) {
    selectors += `[${PARAMETER_TYPE_ATTRIBUTE}],`;
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
        hasProcessingParameters &&
        element.matches(
          `[${PARAMETER_TYPE_ATTRIBUTE}^="${NESTED_TYPE.SIGNAL_ATTRIBUTE}"]`,
        )
      ) {
        const attr = element.getAttribute(PARAMETER_TYPE_ATTRIBUTE);
        if (!attr) {
          throw new Error(`unexpected attribute: ${PARAMETER_TYPE_ATTRIBUTE}`);
        }
        const [, targetName] = attr.split(SIGNAL_ATTRIBUTE_DELIMITER);
        if (!targetName) {
          throw new Error(
            `target name missuse, expected: ${PARAMETER_TYPE_ATTRIBUTE}${SIGNAL_ATTRIBUTE_DELIMITER}attribute-name`,
          );
        }
        const binding = params[getAttributeIndex(element, targetName)];
        if (!(binding instanceof Signal)) {
          throw new Error(`expected signal node, got ${binding}`);
        }
        element.removeAttribute(PARAMETER_TYPE_ATTRIBUTE);
        subscribeSignalAttribute({
          signal: binding,
          element: new WeakRef(element),
          attribute: targetName,
        });
        return;
      }

      if (
        hasProcessingParameters &&
        element.matches(`[${PARAMETER_TYPE_ATTRIBUTE}]`) &&
        element instanceof HTMLSlotElement
      ) {
        const type = element.getAttribute(PARAMETER_TYPE_ATTRIBUTE);
        const nested =
          params[getAttributeIndex(element, PARAMETER_INDEX_ATTRIBUTE)];

        switch (type) {
          case NESTED_TYPE.FRAGMENT: {
            /** @type {DocumentFragment | unknown[]} */
            if (Array.isArray(nested)) {
              element.replaceWith(
                ...nested.map((element) =>
                  element instanceof Node ? element : `${element}`,
                ),
              );
              return;
            }
            element.replaceWith(nested);
            return;
          }
          case NESTED_TYPE.SIGNAL_NODE: {
            if (!(nested instanceof Signal)) {
              throw new Error(`expected signal node, got ${nested}`);
            }
            const mark = document.createTextNode("");
            const bound = document.createTextNode("");
            subscribeSignalNode({
              signal: nested,
              mark: new WeakRef(mark),
              bound: new WeakRef(bound),
            });
            element.replaceWith(mark, nested.value, bound);
            return;
          }
          default:
            throw new Error(`unexpected nested type: ${type}`);
        }
      }

      if (hasBindings && element.hasAttribute(BIND_ATTRIBUTE)) {
        const binding = params[getAttributeIndex(element, BIND_ATTRIBUTE)];
        if (!(binding instanceof Bind)) {
          console.trace();
          console.error(binding);
          throw new Error("binding is expected to be a function");
        }
        binding.current = element;
        element.removeAttribute(BIND_ATTRIBUTE);
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
export function html(...params) {
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

/**
 * Notifies subscribers when value changes
 * To work as signal DOM Node instance must be used as a parameter of `html`
 * Works best with primitives
 *
 * @example
 * const signal = new Signal(10);
 * html`
 *   <button onclick="${() => signal.value += 1}">
 *     increment
 *   </button>
 *   <button onclick="${() => signal.value -= 1}">
 *     decrement
 *   </button>
 *   <p data-signal-value="${signal}">${signal}</p>
 * `;
 *
 * @template {unknown} T
 */
export class Signal {
  #value;
  /** @type {(((value: T) => void) | null)[]} */
  #subscribers = [];

  /** @param {T} value */
  constructor(value) {
    this.#value = value;
  }

  set value(value) {
    this.#value = value;
    this.#publish();
  }

  get value() {
    return this.#value;
  }

  /** @param {(value: T) => void} trigger */
  subscribe(trigger) {
    this.#subscribers.push(trigger);
  }

  /** @param {(value: T) => void} trigger */
  unsubscribe(trigger) {
    const index = this.#subscribers.indexOf(trigger);
    if (index === -1) {
      return;
    }
    this.#subscribers[index] = null;
  }

  #publish() {
    for (let i = 0; i < this.#subscribers.length; i++) {
      const subscriber = this.#subscribers[i];
      if (subscriber) {
        subscriber(this.#value);
        continue;
      }
      this.#subscribers.splice(i, 1);
      i--;
    }
  }
}
