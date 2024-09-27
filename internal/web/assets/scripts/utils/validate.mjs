/**
 * @class
 * @template T
 * @typedef Assertion
 * @property {(data: unknown) => T} parse
 */

/** @implements Assertion<number> */
export class ValidateNumber {
  /** @param {unknown} data */
  parse(data) {
    if (typeof data !== "number" || isNaN(data)) {
      throw new Error(`expected number but got ${typeof data}`);
    }
    return data;
  }
}

/** @implements Assertion<string> */
export class ValidateString {
  /**
   * @param {unknown} data
   */
  parse(data) {
    if (typeof data !== "string") {
      throw new Error(`expected string but got ${typeof data}`);
    }
    return data;
  }
}

/** @template T */
export class ValidateOptional {
  #validator;
  /** @param {Assertion<T>} validator  */
  constructor(validator) {
    this.#validator = validator;
  }

  /** @param {unknown} data  */
  parse(data) {
    if (data === undefined) {
      return undefined;
    }
    return this.#validator.parse(data);
  }
}

/** @template K, T=Record<string, Assertion<K>> */
export class ValidateObject {
  #shape;
  /** @param {T} shape */
  constructor(shape) {
    this.#shape = shape;
  }

  /**
   * @param {unknown} data
   * @returns { { [K in keyof T]: ReturnType<T[K]['parse']> } }
   */
  parse(data) {
    if (typeof data !== "object" || Array.isArray(data) || data === null) {
      throw new Error(`expected object but got ${typeof data}`);
    }
    for (const [key, assertion] of Object.entries(
      /** @type {any} to reset type  */ (this.#shape),
    )) {
      assertion.parse(/** @type Record<string, unknown> */ (data)[key]);
    }
    /** @type {any} to reset type */
    const _data = data;
    return _data;
  }
}
