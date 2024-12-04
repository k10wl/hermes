import { describe, test } from "node:test";

import * as assert from "assert";

import {
  ValidateArray,
  ValidateBoolean,
  ValidateNumber,
  ValidateObject,
  ValidateOptional,
  ValidateString,
} from "./validate.mjs";

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

describe("number validation", () => {
  const nubmerValidation = ValidateNumber;
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

describe("boolean validation", () => {
  const booleanValidation = ValidateBoolean;
  test("should return boolean", () => {
    const truthy = booleanValidation.parse(true);
    assert.strictEqual(truthy, true, "failed to return 'true'");
    const falsy = booleanValidation.parse(false);
    assert.strictEqual(falsy, false, "failed to return 'false'");
  });
  test("should throw if data is not boolean", () => {
    const { truthy, falsy, ...withoutString } = primitives;
    for (const testData of Object.values(withoutString)) {
      assert.throws(() => booleanValidation.parse(testData));
    }
  });
});

describe("string validation", () => {
  const stringValidation = ValidateString;
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

describe("optional", () => {
  test("primitive", () => {
    const optionalStringValidator = new ValidateOptional(ValidateString);
    assert.equal(
      optionalStringValidator.parse(undefined),
      undefined,
      "should allow undefined",
    );
    assert.equal(
      optionalStringValidator.parse("test"),
      "test",
      "should validate propper type",
    );
    assert.throws(
      () => optionalStringValidator.parse({}),
      "should throw with wrong type",
    );
  });

  test("object property", () => {
    const optionalObjectValidator = new ValidateObject({
      optional: new ValidateOptional(ValidateString),
    });
    assert.deepEqual(
      optionalObjectValidator.parse({}),
      {},
      "should skip optional key",
    );
    assert.deepEqual(
      optionalObjectValidator.parse({ optional: "test" }),
      {
        optional: "test",
      },
      "should validate optional key if present",
    );
    assert.throws(
      () => optionalObjectValidator.parse({ optional: 1234 }),
      "should throw if optional value has wrong type",
    );
  });
});

describe("object validation", () => {
  const objectValidation = new ValidateObject({
    string: ValidateString,
    number: ValidateNumber,
    boolean: ValidateBoolean,
    object: new ValidateObject({
      nestedString: ValidateString,
      nestedBoolean: ValidateBoolean,
    }),
  });
  test("should return object", () => {
    const testData = {
      string: "string",
      number: 123,
      boolean: true,
      object: { nestedString: "1234", nestedBoolean: true },
    };
    assert.deepEqual(objectValidation.parse(testData), testData);
  });
  test("should throw if data does not comply schema", () => {
    const testDataArray = [
      {},
      {
        string: 1234,
        number: 123,
        boolean: false,
        object: { nestedString: "1234", nestedBoolean: true },
      },
      {
        string: "string",
        number: "123",
        boolean: false,
        object: { nestedString: "1234", nestedBoolean: true },
      },
      { string: "string", number: 123, boolean: false, object: null },
      {
        string: "string",
        number: 123,
        boolean: false,
        object: { nestedString: 1234, nestedBoolean: true },
      },
      {
        string: "string",
        number: 123,
        boolean: "false",
        object: { nestedString: "1234", nestedBoolean: true },
      },
      {
        string: "string",
        number: 123,
        boolean: true,
        object: { nestedString: "1234", nestedBoolean: "true" },
      },
    ];
    for (const testData of testDataArray) {
      assert.throws(() => objectValidation.parse(testData));
    }
  });
});

describe("array validation", () => {
  const arrayValidation = new ValidateArray(
    new ValidateObject({ id: ValidateString }),
  );
  test("should return array", () => {
    assert.deepStrictEqual(
      arrayValidation.parse([{ id: "qwerty-1" }, { id: "qwerty-2" }]),
      [{ id: "qwerty-1" }, { id: "qwerty-2" }],
    );
  });
  test("should throw if data does not comply schema", () => {
    assert.throws(() => arrayValidation.parse(["1234", 12340]));
  });
});
