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
    /** @typedef {Map<Callback<any>, ({callback: Callback<any>, teardown: Teardown})>} CallbackWithTeardown */
    /** @type {Map<Keys, CallbackWithTeardown>} */
    this.handlers = new Map();
  }

  /**
   * @template {Keys | Keys[]} T
   * @param {T} key
   * @param {Callback<T extends Keys[] ? T[number] : T>} callback
   * @returns {Teardown}
   */
  on(key, callback) {
    if (!Array.isArray(key)) {
      return this.#onSingle(
        /** @type {any} shut up, JSDoc, this type is definitely narrowed */ (
          key
        ),
        callback,
      );
    }
    return this.#onMultiple(key, callback);
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
   * @param {string} key
   * @return {((Callback<T>)[]) | undefined}
   */
  getCallbacks(key) {
    const handlers = this.handlers.get(key);
    if (!handlers) {
      return;
    }
    return Array.from(handlers.values()).map(({ callback }) => callback);
  }

  /**
   * @template {Keys} T
   * @param {T} key
   * @param {Callback<T>} callback
   * @returns {Teardown}
   */
  #onSingle(key, callback) {
    let handlers = this.handlers.get(key);
    if (!handlers) {
      handlers = new Map();
      this.handlers.set(key, handlers);
    }
    const handler = {
      callback,
      teardown: this.#teardown(key, callback),
    };
    handlers.set(callback, handler);
    return handler.teardown;
  }

  /**
   * @template {Keys[]} T
   * @param {T} keys
   * @param {Callback<T[number]>} callback
   * @returns {Teardown}
   */
  #onMultiple(keys, callback) {
    const teardown = keys.map((key) => this.#onSingle(key, callback));
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
