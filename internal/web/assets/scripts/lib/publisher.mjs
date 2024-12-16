/**
 * @template T
 * @typedef Observer
 * @property {(current: T, previous: T) => void} notify
 * @class
 */

/** @template T */
export class Publisher {
  /** @typedef {Observer<T>} ConcreteObserver */
  /** @typedef {() => boolean} detach */

  /** @type ConcreteObserver[] */
  #observers = [];

  /** @param {T} initialValue */
  constructor(initialValue) {
    this.value = initialValue;
  }

  /**
   * @param {ConcreteObserver} observer
   * @returns {detach}
   */
  attach(observer) {
    this.#observers.push(observer);
    return () => this.detach(observer);
  }

  /**
   * @param {ConcreteObserver} observer
   * @returns {boolean}
   */
  detach(observer) {
    const i = this.#observers.indexOf(observer);
    if (i === -1) {
      return false;
    }
    this.#observers.splice(i, 1);
    return true;
  }

  /** @param {T | ((currentValue: T) => T)} value */
  update(value) {
    const update = value instanceof Function ? value(this.value) : value;
    if (update === this.value) {
      return;
    }
    const previous = this.value;
    this.value = update;
    this.#notifyWithPrevious(previous);
  }

  /**
   * @param {T} previous
   */
  #notifyWithPrevious(previous) {
    for (let i = 0; i < this.#observers.length; i++) {
      this.#observers[i]?.notify(this.value, previous);
    }
  }

  notify() {
    this.#notifyWithPrevious(this.value);
  }
}
