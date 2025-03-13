/**
 * @template T
 * @typedef Iterator
 * @property {() => Promise<T[]> | T[]} next
 * @property {boolean} hasMore
 */

/**
 * @template T
 * @typedef Renderer
 * @property {(data: T) => DocumentFragment} createElement
 */

/**
 * @template T
 */
export class PaginatedList extends HTMLElement {
  /** @type {Iterator<T> | null} */
  #iterator = null;
  /** @type {Renderer<T> | null} */
  #renderer = null;
  scrollTriggerThreshold = 200;
  #loading = false;
  #observer = new IntersectionObserver(async ([entry]) => {
    if (!this.#iterator) {
      throw new Error("observer triggered before iterator was set");
    }
    if (!entry?.isIntersecting || this.#loading) {
      return;
    }
    await this.#withLoading(this.#loadNext());
    this.#positionTriggerElement();
  });
  /** @type {Element | null} */
  #currentlyObserving = null;

  /** @type {HTMLElement} */
  #startContainer = document.createElement("div");

  /** @type {HTMLElement} */
  #content = document.createElement("div");

  /** @type {HTMLElement[]} */
  #pageContainers = [];

  constructor() {
    super();
    this.append(this.#startContainer, this.#content);
  }

  #positionTriggerElement() {
    const toObserve = this.#pageContainers.at(-1)?.firstElementChild;
    if (!toObserve) {
      return;
    }
    if (this.#currentlyObserving) {
      this.#observer.unobserve(this.#currentlyObserving);
    }
    this.#observer.observe(toObserve);
    this.#currentlyObserving = toObserve;
  }

  /**
   * @template K
   * @param {Promise<K>} promise
   * @returns {Promise<K>}
   */
  async #withLoading(promise) {
    try {
      this.#loading = true;
      return await promise;
    } finally {
      this.#loading = false;
    }
  }

  /** @returns {Promise<boolean>} boolean indicator of more records */
  async #loadNext() {
    if (this.#iterator === null || this.#renderer === null) {
      throw new Error("iterator or renderer not set before loading next");
    }
    const page = await this.#iterator.next();
    /** @type {DocumentFragment[]} */
    const el = [];
    for (const data of page) {
      el.push(this.#renderer.createElement(data));
    }
    if (el.length === 0) {
      return false;
    }
    const container = document.createElement("div");
    container.append(...el);
    this.#content.append(container);
    this.#pageContainers.push(container);
    return true;
  }

  async init() {
    await this.#withLoading(this.#loadNext());
    this.#positionTriggerElement();
  }

  /** @param {Renderer<T>} renderer  */
  setRenderer(renderer) {
    this.#renderer = renderer;
  }

  /** @param {Iterator<T>} iterator  */
  setIterator(iterator) {
    this.#iterator = iterator;
  }

  /** @param {Node[]} nodes  */
  prepandNodes(...nodes) {
    this.#startContainer.prepend(...nodes);
  }
}

customElements.define("h-paginated-list", PaginatedList);
