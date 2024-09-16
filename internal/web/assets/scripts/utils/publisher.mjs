/**
 * @template T
 * @typedef Observer
 * @property {(value: T) => void} notify
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
   * @param {ConcreteObserver} callback
   * @returns {boolean}
   */
  detach(callback) {
    const i = this.#observers.indexOf(callback);
    if (i === -1) {
      return false;
    }
    this.#observers.splice(i, 1);
    return true;
  }

  /** @param {T | ((currentValue: T) => T)} value */
  update(value) {
    this.value = value instanceof Function ? value(this.value) : value;
    this.notify();
  }

  notify() {
    for (let i = 0; i < this.#observers.length; i++) {
      this.#observers[i]?.notify(this.value);
    }
  }
}
