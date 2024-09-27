import { describe, test } from "node:test";

import * as assert from "assert";

import { ValidateNumber, ValidateObject, ValidateString } from "./validate.mjs";

const primitives = {
  string: "string",
  number: 1,
  undefined: undefined,
  boolean: true,
  symbol: Symbol,
  NaN: NaN,
  null: null,
};

describe("number validation", () => {
  const nubmerValidation = new ValidateNumber();
  test("should return number", () => {
    const res = nubmerValidation.parse(2);
    assert.strictEqual(res, 2, "failed to return same value");
  });
  test("should throw if data is not number", () => {
    const { number, ...withoutNumber } = primitives;
    for (const testData of Object.values(withoutNumber)) {
      assert.throws(() => nubmerValidation.parse(testData));
    }
  });
});

describe("string validation", () => {
  const stringValidation = new ValidateString();
  test("should return string", () => {
    const res = stringValidation.parse("test");
    assert.strictEqual(res, "test", "failed to return same value");
  });
  test("should throw if data is not string", () => {
    const { string, ...withoutString } = primitives;
    for (const testData of Object.values(withoutString)) {
      assert.throws(() => stringValidation.parse(testData));
    }
  });
});

describe("object validation", () => {
  const objectValidation = new ValidateObject({
    string: new ValidateString(),
    number: new ValidateNumber(),
    object: new ValidateObject({
      nestedString: new ValidateString(),
    }),
  });
  test("should return object", () => {
    const testData = {
      string: "string",
      number: 123,
      object: { nestedString: "1234" },
    };
    assert.deepEqual(objectValidation.parse(testData), testData);
  });
  test("should throw if data does not comply schema", () => {
    const testDataArray = [
      {},
      { string: 1234, number: 123, object: { nestedString: "1234" } },
      { string: "string", number: "123", object: { nestedString: "1234" } },
      { string: "string", number: 123, object: null },
      { string: "string", number: 123, object: { nestedString: 1234 } },
    ];
    for (const testData of testDataArray) {
      assert.throws(() => objectValidation.parse(testData));
    }
  });
});
