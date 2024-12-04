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

describe("number assertion", () => {
  const nubmerAssertion = AssertNumber;
  test("should return number", () => {
    /** @type {number} */
    const res = nubmerAssertion.check(2);
    assert.strictEqual(res, 2, "failed to return same value");
  });
  test("should throw if data is not number", () => {
    const { number, ...withoutNumber } = primitives;
    for (const testData of Object.values(withoutNumber)) {
      assert.throws(() => nubmerAssertion.check(testData));
    }
  });
});

describe("boolean assertion", () => {
  const booleanAssertion = AssertBoolean;
  test("should return boolean", () => {
    /** @type {boolean} */
    const truthy = booleanAssertion.check(true);
    assert.strictEqual(truthy, true, "failed to return 'true'");
    /** @type {boolean} */
    const falsy = booleanAssertion.check(false);
    assert.strictEqual(falsy, false, "failed to return 'false'");
  });
  test("should throw if data is not boolean", () => {
    const { truthy, falsy, ...withoutString } = primitives;
    for (const testData of Object.values(withoutString)) {
      assert.throws(() => booleanAssertion.check(testData));
    }
  });
});

describe("string assertion", () => {
  const stringAssertion = AssertString;
  test("should return string", () => {
    /** @type {string} */
    const res = stringAssertion.check("test");
    assert.strictEqual(res, "test", "failed to return same value");
  });
  test("should throw if data is not string", () => {
    const { string, ...withoutString } = primitives;
    for (const testData of Object.values(withoutString)) {
      assert.throws(() => stringAssertion.check(testData));
    }
  });
});

describe("optional", () => {
  test("primitive", () => {
    const optionalStringAssertion = new AssertOptional(AssertString);
    /** @type {string | undefined} */
    let stringCheck = optionalStringAssertion.check(undefined);
    assert.equal(stringCheck, undefined, "should allow undefined");
    stringCheck = optionalStringAssertion.check("test");
    assert.equal(stringCheck, "test", "should assert propper type");
    assert.throws(
      () => optionalStringAssertion.check({}),
      "should throw with wrong type",
    );
  });

  test("object property", () => {
    const optionalObjectAssertion = new AssertObject({
      optional: new AssertOptional(AssertString),
    });
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
      () => optionalObjectAssertion.check({ optional: 1234 }),
      "should throw if optional value has wrong type",
    );
  });
});

describe("object assertion", () => {
  class Test {
    constructor() {}
  }
  const objectAssertion = new AssertObject({
    string: AssertString,
    number: AssertNumber,
    boolean: AssertBoolean,
    object: new AssertObject({
      nestedString: AssertString,
      nestedBoolean: AssertBoolean,
    }),
    array: new AssertArray(AssertString),
    instance: new AssertInstance(Test),
  });
  test("should return object", () => {
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
  test("should throw if data does not comply schema", () => {
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
      assert.throws(() => objectAssertion.check(testData));
    }
  });
});

describe("array assertion", () => {
  const arrayAssertion = new AssertArray(
    new AssertObject({ id: AssertString }),
  );
  test("should return array", () => {
    /** @type {{id: string}[]} */
    const data = arrayAssertion.check([{ id: "qwerty-1" }, { id: "qwerty-2" }]);
    assert.deepStrictEqual(data, [{ id: "qwerty-1" }, { id: "qwerty-2" }]);
  });
  test("should throw if data does not comply schema", () => {
    assert.throws(() => arrayAssertion.check(["1234", 12340]));
  });
});

describe("instance assertion", () => {
  test("should assert instance", () => {
    const initial = new ArrayBuffer();
    /** @type {ArrayBuffer} */
    const checked = AssertInstance.once(initial, ArrayBuffer);
    assert.equal(
      initial,
      checked,
      "should return same array buffer after check",
    );
  });
  test("should throw if data does not follow instance", () => {
    assert.throws(() => AssertInstance.once(42069, ArrayBuffer));
  });

  class Test {
    foo = 2;
    constructor() {}
  }
  test("should assert if initiated with constructor", () => {
    const assertion = new AssertInstance(Test);
    const test = new Test();
    /** @type {InstanceType<typeof Test>} foo */
    const checked = assertion.check(test);
    assert.equal(checked, test, "should reutrn same array buffer after check");
  });
  test("should throw if data does not follow instance", () => {
    const assertion = new AssertInstance(Test);
    assert.throws(
      () => assertion.check(""),
      "should reutrn same array buffer after check",
    );
  });
});
