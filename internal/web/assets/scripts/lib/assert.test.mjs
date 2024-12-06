import { describe, test } from "node:test";

import * as assert from "assert";

import {
  AssertArray,
  AssertBoolean,
  AssertInstance,
  AssertNumber,
  AssertObject,
  AssertOptional,
  AssertString,
  AssertTruthy,
} from "./assert.mjs";

const primitives = {
  string: "string",
  number: 1,
  undefined: undefined,
  truthy: true,
  falsy: true,
  symbol: Symbol,
  NaN: NaN,
  null: null,
};

describe("AssertNumber", () => {
  const nubmerAssertion = AssertNumber;
  test("returns the same number when checked", () => {
    /** @type {number} */
    const res = nubmerAssertion.check(2);
    assert.strictEqual(res, 2, "failed to return same value");
  });
  test("throws an error if data is not a number", () => {
    const { number, ...withoutNumber } = primitives;
    for (const testData of Object.values(withoutNumber)) {
      assert.throws(
        () => nubmerAssertion.check(testData, "custom error reason"),
        /custom error reason/,
      );
    }
  });
});

describe("AssertBoolean", () => {
  const booleanAssertion = AssertBoolean;
  test("returns true when checked with true", () => {
    /** @type {boolean} */
    const truthy = booleanAssertion.check(true);
    assert.strictEqual(truthy, true, "failed to return 'true'");
  });
  test("returns false when checked with false", () => {
    /** @type {boolean} */
    const falsy = booleanAssertion.check(false);
    assert.strictEqual(falsy, false, "failed to return 'false'");
  });
  test("throws an error if data is not a boolean", () => {
    const { truthy, falsy, ...withoutString } = primitives;
    for (const testData of Object.values(withoutString)) {
      assert.throws(
        () => booleanAssertion.check(testData, "custom error reason"),
        /custom error reason/,
      );
    }
  });
});

describe("AssertString", () => {
  const stringAssertion = AssertString;
  test("returns the same string when checked", () => {
    /** @type {string} */
    const res = stringAssertion.check("test");
    assert.strictEqual(res, "test", "failed to return same value");
  });
  test("throws an error if data is not a string", () => {
    const { string, ...withoutString } = primitives;
    for (const testData of Object.values(withoutString)) {
      assert.throws(
        () => stringAssertion.check(testData, "custom error reason"),
        /custom error reason/,
      );
    }
  });
});

describe("AssertOptional", () => {
  test("allows undefined as a valid value", () => {
    const optionalStringAssertion = new AssertOptional(
      AssertString,
      "custom error reason",
    );
    /** @type {string | undefined} */
    let stringCheck = optionalStringAssertion.check(undefined);
    assert.equal(stringCheck, undefined, "should allow undefined");
    stringCheck = optionalStringAssertion.check("test");
    assert.equal(stringCheck, "test", "should assert proper type");
  });

  test("throws an error for wrong type when optional", () => {
    const optionalStringAssertion = new AssertOptional(
      AssertString,
      "custom error reason",
    );
    assert.throws(
      () => optionalStringAssertion.check({}, "inline error reason"),
      /(custom error reason).*(inline error reason)/,
      "should throw with wrong type",
    );
  });

  test("handles optional object properties correctly", () => {
    const optionalObjectAssertion = new AssertObject(
      {
        optional: new AssertOptional(
          AssertString,
          "nested custom error reason",
        ),
      },
      "custom error reason",
    );
    /** @type {{ optional: string | undefined }} */
    let data = optionalObjectAssertion.check({});
    assert.deepEqual(data, {}, "should skip optional key");
    data = optionalObjectAssertion.check({ optional: "test" });
    assert.deepEqual(
      data,
      {
        optional: "test",
      },
      "should assert optional key if present",
    );
    assert.throws(
      () =>
        optionalObjectAssertion.check(
          { optional: 1234 },
          "inline error reason",
        ),
      /(custom error reason).*(inline error reason)/,
      "should throw if optional value has wrong type",
    );
  });
});

describe("AssertObject", () => {
  class Test {
    constructor() {}
  }
  const objectAssertion = new AssertObject(
    {
      string: AssertString,
      number: AssertNumber,
      boolean: AssertBoolean,
      object: new AssertObject(
        {
          nestedString: AssertString,
          nestedBoolean: AssertBoolean,
        },
        "nested error reason",
      ),
      array: new AssertArray(AssertString, "nested error reason"),
      instance: new AssertInstance(Test, "nested error reason"),
    },
    "custom error reason",
  );
  test("returns the same object when checked", () => {
    const testData = {
      string: "string",
      number: 123,
      boolean: true,
      object: { nestedString: "1234", nestedBoolean: true },
      array: ["foo"],
      instance: new Test(),
    };
    /** @type {typeof testData} */
    const res = objectAssertion.check(testData);
    assert.deepEqual(res, testData);
  });
  test("throws an error if data does not comply with schema", () => {
    const testDataArray = [
      {},
      {
        string: 1234,
        number: 123,
        boolean: false,
        object: { nestedString: "1234", nestedBoolean: true },
        array: ["foo"],
        instance: new Test(),
      },
      {
        string: "string",
        number: "123",
        boolean: false,
        object: { nestedString: "1234", nestedBoolean: true },
        array: ["foo"],
        instance: new Test(),
      },
      {
        string: "string",
        number: 123,
        boolean: false,
        object: null,
        array: ["foo"],
        instance: new Test(),
      },
      {
        string: "string",
        number: 123,
        boolean: false,
        object: { nestedString: 1234, nestedBoolean: true },
        array: ["foo"],
        instance: new Test(),
      },
      {
        string: "string",
        number: 123,
        boolean: "false",
        object: { nestedString: "1234", nestedBoolean: true },
        array: ["foo"],
        instance: new Test(),
      },
      {
        string: "string",
        number: 123,
        boolean: true,
        object: { nestedString: "1234", nestedBoolean: "true" },
        array: ["foo"],
        instance: new Test(),
      },
      {
        string: "string",
        number: 123,
        boolean: false,
        object: { nestedString: "1234", nestedBoolean: true },
        array: [1324],
        instance: new Test(),
      },
      {
        string: "string",
        number: 123,
        boolean: false,
        object: { nestedString: "1234", nestedBoolean: true },
        array: ["foo"],
        instance: "",
      },
    ];
    for (const testData of testDataArray) {
      assert.throws(
        () => objectAssertion.check(testData, "inline error reason"),
        (error) => {
          if (!(error instanceof Error)) {
            return false;
          }
          if (/(object)|(optional)|(array)/.test(error.message)) {
            return /(nested error reason).*(custom error reason).*(inline error reason)/.test(
              error.message,
            );
          }
          return /(custom error reason).*(inline error reason)/.test(
            error.message,
          );
        },
      );
    }
  });
});

describe("AssertArray", () => {
  const arrayAssertion = new AssertArray(
    new AssertObject({ id: AssertString }),
    "custom error reason",
  );
  test("returns the same array when checked", () => {
    /** @type {{id: string}[]} */
    const data = arrayAssertion.check([{ id: "qwerty-1" }, { id: "qwerty-2" }]);
    assert.deepStrictEqual(data, [{ id: "qwerty-1" }, { id: "qwerty-2" }]);
  });
  test("throws an error if data does not comply with schema", () => {
    assert.throws(
      () => arrayAssertion.check(["1234", 12340], "inline error reason"),
      /(custom error reason).*(inline error reason)/,
    );
  });
});

describe("AssertInstance", () => {
  test("asserts instance correctly", () => {
    const initial = new ArrayBuffer();
    /** @type {ArrayBuffer} */
    const checked = AssertInstance.once(initial, ArrayBuffer);
    assert.equal(
      initial,
      checked,
      "should return same array buffer after check",
    );
  });
  test("throws an error if data does not follow instance", () => {
    assert.throws(
      () => AssertInstance.once(42069, ArrayBuffer, "custom error reason"),
      /custom error reason/,
    );
  });

  class Test {
    foo = 2;
    constructor() {}
  }
  test("asserts instance when initiated with constructor", () => {
    const assertion = new AssertInstance(Test, "custom error reason");
    const test = new Test();
    /** @type {InstanceType<typeof Test>} foo */
    const checked = assertion.check(test);
    assert.equal(checked, test, "should return same instance after check");
  });
  test("throws an error if data does not follow instance", () => {
    const assertion = new AssertInstance(Test, "custom error reason");
    assert.throws(
      () => assertion.check("", "inline error reason"),
      /(custom error reason).*(inline error reason)/,
      "should throw if data does not follow instance",
    );
  });
});

describe("AssertTruthy", () => {
  test("does not throw an error for true assertion", () => {
    assert.doesNotThrow(() => AssertTruthy.check(true));
  });

  test("throws an error for false assertion", () => {
    assert.throws(
      () => AssertTruthy.check(false, "custom error reason"),
      /custom error reason/,
    );
  });
});
