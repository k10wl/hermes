/**
 * @class
 * @template T
 * @typedef Assertion
 * @property {(data: unknown) => T} check
 */

export class AssertNumber {
  /** @param {unknown} data */
  static check(data) {
    if (typeof data !== "number" || isNaN(data)) {
      throw new Error(`expected number but got ${typeof data}`);
    }
    return data;
  }
}

export class AssertString {
  /**
   * @param {unknown} data
   */
  static check(data) {
    if (typeof data !== "string") {
      throw new Error(`expected string but got ${typeof data}`);
    }
    return data;
  }
}

/** @template T */
export class AssertOptional {
  #validator;
  /** @param {Assertion<T>} validator  */
  constructor(validator) {
    this.#validator = validator;
  }

  /** @param {unknown} data  */
  check(data) {
    if (data === undefined) {
      return undefined;
    }
    return this.#validator.check(data);
  }
}

export class AssertBoolean {
  /**
   * @param {unknown} data
   */
  static check(data) {
    if (typeof data !== "boolean") {
      throw new Error(`expected boolean but got ${typeof data}`);
    }
    return data;
  }
}

/** @template K, T=Record<string, Assertion<K>> */
export class AssertObject {
  #shape;
  /** @param {T} shape */
  constructor(shape) {
    this.#shape = shape;
  }

  /**
   * @param {unknown} data
   * @returns { { [K in keyof T]: ReturnType<T[K]['check']> } }
   */
  check(data) {
    if (typeof data !== "object" || Array.isArray(data) || data === null) {
      throw new Error(`expected object but got ${typeof data}`);
    }
    for (const [key, assertion] of Object.entries(
      /** @type {any} to reset type  */ (this.#shape),
    )) {
      assertion.check(/** @type Record<string, unknown> */ (data)[key]);
    }
    /** @type {any} to reset type */
    const _data = data;
    return _data;
  }
}

/** @template K, T=Assertion<K> */
export class AssertArray {
  #shape;
  /** @param {T} shape */
  constructor(shape) {
    this.#shape = shape;
  }

  /**
   * @param {unknown} data
   * @returns { ReturnType<T['check']>[] } }
   */
  check(data) {
    if (!Array.isArray(data)) {
      throw new Error(`expected array but got ${typeof data}`);
    }
    for (const value of data) {
      /** @type Assertion<K> */ (this.#shape).check(value);
    }
    /** @type {any} to reset type */
    const _data = data;
    return _data;
  }
}

/** @template T */
export class AssertInstance {
  #instance;
  /** @param {T} instance */
  constructor(instance) {
    this.#instance = instance;
  }
  instance() {
    return this.#instance;
  }

  /**
   * @param {unknown} data
   * @returns {InstanceType<T>}
   */
  check(data) {
    return AssertInstance.once(data, this.#instance);
  }

  /**
   * @template K
   * @param {unknown} data
   * @param {K} type
   * @returns {InstanceType<K>}
   */
  static once(data, type) {
    // @ts-expect-error expected to throw on bad data
    if (data instanceof type) {
      /** @type {any} */
      const any = data;
      return any;
    }
    throw new Error(`Data ${data} does not have the right type '${type}'!`);
  }
}
