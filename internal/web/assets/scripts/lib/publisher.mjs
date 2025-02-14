/**
 * @template T
 * @typedef Subscriber
 * @type {{notify: (value: T) => void}}
 */

/**
 * @template T
 */
export class Publisher {
  /** @typedef {Subscriber<T>} ConcreteSubscriber */
  /** @typedef {() => boolean} detach */

  /** @type {ConcreteSubscriber[]} */
  subscribers = [];

  /** @type {T} */
  value;

  /** @param {T} initialValue */
  constructor(initialValue) {
    this.value = initialValue;
  }

  /**
   * @param {ConcreteSubscriber} subscriber
   * @returns {detach}
   */
  subscribe(subscriber) {
    this.subscribers.push(subscriber);
    return () => this.unsubscribe(subscriber);
  }

  /**
   * @param {ConcreteSubscriber} subscriber
   * @returns {boolean}
   */
  unsubscribe(subscriber) {
    const i = this.subscribers.indexOf(subscriber);
    if (i === -1) {
      return false;
    }
    this.subscribers.splice(i, 1);
    return true;
  }

  /** @param {T | ((currentValue: T) => T)} value */
  update(value) {
    const update =
      typeof value === "function"
        ? /** @type {Function} */ (value)(this.value)
        : value;
    if (update === this.value) {
      return;
    }
    this.value = update;
    this.notify();
  }

  notify() {
    let delta = 0;
    for (let i = 0; i < this.subscribers.length; i++) {
      const subscriber = this.subscribers[i];
      if (!subscriber) {
        this.subscribers.splice(i - delta, 1);
        delta += 1;
        i -= 1;
        continue;
      }
      subscriber.notify(this.value);
    }
  }
}
