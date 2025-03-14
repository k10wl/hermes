/** @template {Record<string, unknown>} Data */
export class CallbackTracker {
  /** @typedef {keyof Data} Keys */
  /**
   * @template {Keys} T
   * @typedef {(data: Data[T] extends abstract new (...args: any) => any ? InstanceType<Data[T]> : Data[T]) => void} Callback
   */
  /** @typedef {() => void} Teardown */

  /** @param {Data} _data used only in JSDoc, but is passed as value to allow using static types without JSDoc comments */
  constructor(_data) {
    /** @typedef {Map<Callback<any>, ({callback: Callback<any>, teardown: Teardown, priority: number})>} CallbackWithTeardown */
    /** @type {Map<Keys, CallbackWithTeardown>} */
    this.handlers = new Map();
  }

  /**
   * @template {Keys | Keys[]} T
   * @param {T} key
   * @param {Callback<T extends Keys[] ? T[number] : T>} callback
   * @param {{priority: number}} [options]
   * @returns {Teardown}
   */
  on(key, callback, options = { priority: 1 }) {
    if (!Array.isArray(key)) {
      return this.#onSingle(
        /** @type {any} shut up, JSDoc, this type is definitely narrowed */ (
          key
        ),
        callback,
        options,
      );
    }
    return this.#onMultiple(key, callback, options);
  }

  /**
   * @template {Keys} T
   * @param {T} key
   * @param {(...args: any) => void} callback
   */
  off(key, callback) {
    const listeners = this.handlers.get(key);
    if (!listeners) {
      return;
    }
    listeners.get(callback)?.teardown();
  }

  /**
   * @template {Keys} T
   * @param {string[]} keys
   * @return {((Callback<T>)[]) | undefined}
   */
  getCallbacks(...keys) {
    const handlers = [];
    for (const key of keys) {
      handlers.push(...(this.handlers.get(key)?.values().toArray() ?? []));
    }
    if (handlers.length === 0) {
      return undefined;
    }
    return handlers
      .sort((a, b) => b.priority - a.priority)
      .map(({ callback }) => callback);
  }

  /**
   * @template {Keys} T
   * @param {T} key
   * @param {Callback<T>} callback
   * @param {{priority: number}} options
   * @returns {Teardown}
   */
  #onSingle(key, callback, options) {
    const handler = {
      callback,
      teardown: this.#teardown(key, callback),
      priority: options.priority,
    };
    const existingHandlers = this.handlers.get(key)?.entries().toArray() ?? [];
    existingHandlers.push([callback, handler]);
    this.handlers.set(
      key,
      new Map(existingHandlers.sort(([, a], [, b]) => b.priority - a.priority)),
    );
    return handler.teardown;
  }

  /**
   * @template {Keys[]} T
   * @param {T} keys
   * @param {Callback<T[number]>} callback
   * @param {{priority: number}} options
   * @returns {Teardown}
   */
  #onMultiple(keys, callback, options) {
    const teardown = keys.map((key) => this.#onSingle(key, callback, options));
    return () => teardown.forEach((cb) => cb());
  }

  /**
   * @template {Keys} T
   * @param {T} key
   * @param {Callback<T>} callback
   * @return {Teardown}
   */
  #teardown(key, callback) {
    return () => {
      const handlers = this.handlers.get(key);
      if (!handlers) {
        return;
      }
      handlers.delete(callback);
      if (handlers.size < 1) {
        this.handlers.delete(key);
      }
    };
  }
}
