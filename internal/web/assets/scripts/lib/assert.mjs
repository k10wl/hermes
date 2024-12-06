/**
 * @class
 * @template T
 * @typedef Assertion
 * @property {(data: unknown, reason?: string) => T} check
 */

export class AssertNumber {
  /**
   * @param {unknown} data
   * @param {string} [reason]
   */
  static check(data, reason) {
    if (typeof data !== "number" || isNaN(data)) {
      throw new Error(
        combineString(`expected number but got ${typeof data}`, reason),
      );
    }
    return data;
  }
}

export class AssertString {
  /**
   * @param {unknown} data
   * @param {string} [reason]
   */
  static check(data, reason) {
    if (typeof data !== "string") {
      throw new Error(
        combineString(`expected string but got ${typeof data}`, reason),
      );
    }
    return data;
  }
}

/** @template T */
export class AssertOptional {
  #validator;
  #reason;
  /**
   * @param {Assertion<T>} validator
   * @param {string} [reason]
   */
  constructor(validator, reason) {
    this.#validator = validator;
    this.#reason = reason;
  }

  /**
   * @param {unknown} data
   * @param {string} [reason]
   */
  check(data, reason) {
    if (data === undefined) {
      return undefined;
    }
    return this.#validator.check(data, combineString(this.#reason, reason));
  }
}

export class AssertBoolean {
  /**
   * @param {unknown} data
   * @param {string} [reason]
   */
  static check(data, reason) {
    if (typeof data !== "boolean") {
      throw new Error(
        combineString(`expected boolean but got ${typeof data}`, reason),
      );
    }
    return data;
  }
}

/** @template K, T=Record<string, Assertion<K>> */
export class AssertObject {
  #shape;
  #reason;
  /**
   * @param {T} shape
   * @param {string} [reason]
   */
  constructor(shape, reason) {
    this.#shape = shape;
    this.#reason = reason;
  }

  /**
   * @param {unknown} data
   * @param {string} [reason]
   * @returns { { [K in keyof T]: ReturnType<T[K]['check']> } }
   */
  check(data, reason) {
    if (typeof data !== "object" || Array.isArray(data) || data === null) {
      throw new Error(
        combineString(
          this.#reason,
          reason,
          `expected object but got ${typeof data}`,
        ),
      );
    }
    for (const [key, assertion] of Object.entries(
      /** @type {any} to reset type  */ (this.#shape),
    )) {
      try {
        assertion.check(/** @type Record<string, unknown> */ (data)[key]);
      } catch (e) {
        const error = AssertInstance.once(
          e,
          Error,
          "did not throw error class",
        );
        error;
        throw new Error(
          combineString(`${key}`, error.message, this.#reason, reason),
        );
      }
    }
    /** @type {any} to reset type */
    const _data = data;
    return _data;
  }
}

/** @template K, T=Assertion<K> */
export class AssertArray {
  #shape;
  #reason;
  /**
   * @param {T} shape
   * @param {string} [reason]
   */
  constructor(shape, reason) {
    this.#shape = shape;
    this.#reason = reason;
  }

  /**
   * @param {unknown} data
   * @param {string} [reason]
   * @returns { ReturnType<T['check']>[] } }
   */
  check(data, reason) {
    if (!Array.isArray(data)) {
      throw new Error(
        combineString(
          `expected array but got ${typeof data}`,
          this.#reason,
          reason,
        ),
      );
    }
    for (const value of data) {
      /** @type Assertion<K> */ (this.#shape).check(
        value,
        combineString(this.#reason, reason),
      );
    }
    /** @type {any} to reset type */
    const _data = data;
    return _data;
  }
}

/** @template T */
export class AssertInstance {
  #instance;
  #reason;
  /**
   * @param {T} instance
   * @param {string} [reason]
   */
  constructor(instance, reason) {
    this.#instance = instance;
    this.#reason = reason;
  }

  /**
   * @param {unknown} data
   * @param {string} [reason]
   * @returns {InstanceType<T>}
   */
  check(data, reason) {
    return AssertInstance.once(
      data,
      this.#instance,
      combineString(this.#reason, reason),
    );
  }

  /**
   * @template K
   * @param {unknown} data
   * @param {string} [reason]
   * @param {K} type
   * @returns {InstanceType<K>}
   */
  static once(data, type, reason) {
    // @ts-expect-error expected to throw on bad data
    if (data instanceof type) {
      /** @type {any} */
      const any = data;
      return any;
    }
    throw new Error(
      combineString(
        `Data ${data} does not have the right type '${type}'!`,
        reason,
      ),
    );
  }
}

export class AssertTruthy {
  /**
   * @param {boolean} assertion
   * @param {string} [reason]
   * @returns {void}
   */
  static check(assertion, reason) {
    if (!assertion) {
      throw new Error(combineString(`Truthy assertion failed`, reason));
    }
  }
}

/**
 * @param {(string|undefined)[]} strings
 * @returns {string}
 */
function combineString(...strings) {
  return strings.filter(Boolean).join(" - ");
}
