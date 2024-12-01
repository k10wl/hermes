import * as assert from "node:assert";
import test, { describe } from "node:test";

import { CallbackTracker } from "./callback-tracker.mjs";

const sampleDatasetForCompletionCheck = {
  number: 2134,
  string: "string",
  array: [],
};

describe("should track callbacks", () => {
  test("using returned teardrop to unsub", () => {
    const callbackTracker = new CallbackTracker(
      sampleDatasetForCompletionCheck,
    );
    const cb =
      /** @param {number} _ */
      (_) => {};
    const teardown = callbackTracker.on("number", cb);
    assert.equal(
      cb,
      callbackTracker.handlers.get("number")?.get(cb)?.callback,
      "callback was not added to handlers",
    );
    teardown();
    assert.equal(
      callbackTracker.handlers.get("number")?.get(cb)?.callback,
      undefined,
      "callback was not removed from handlers",
    );
  });

  test("using 'off' to unsub", () => {
    const callbackTracker = new CallbackTracker(
      sampleDatasetForCompletionCheck,
    );
    const cb =
      /** @param {number} _ */
      (_) => {};
    callbackTracker.on("number", cb);
    assert.equal(
      cb,
      callbackTracker.handlers.get("number")?.get(cb)?.callback,
      "callback was not added to handlers",
    );
    callbackTracker.off("number", cb);
    assert.equal(
      callbackTracker.handlers.get("number")?.get(cb)?.callback,
      undefined,
      "callback was not removed from handlers",
    );
  });

  test("having multiple listeners", () => {
    const callbackTracker = new CallbackTracker(
      sampleDatasetForCompletionCheck,
    );
    const cb1 =
      /** @param {number} _ */
      (_) => {};
    const cb2 =
      /** @param {number} _ */
      (_) => {};
    callbackTracker.on("number", cb1);
    callbackTracker.on("number", cb2);
    assert.equal(
      cb1,
      callbackTracker.handlers.get("number")?.get(cb1)?.callback,
      "callback was not added to handlers",
    );
    callbackTracker.off("number", cb1);
    assert.notEqual(
      cb1,
      callbackTracker.handlers.get("number")?.get(cb1)?.callback,
      "callback was not removed from handlers",
    );
    assert.equal(
      cb2,
      callbackTracker.handlers.get("number")?.get(cb2)?.callback,
      "removed callback that should be remained",
    );
  });

  test("array keys tracking", () => {
    const callbackTracker = new CallbackTracker(
      sampleDatasetForCompletionCheck,
    );
    const cb = () => {};
    callbackTracker.on(["string", "number"], cb);
    assert.equal(
      cb,
      callbackTracker.handlers.get("string")?.get(cb)?.callback,
      "string callback not tracked",
    );
    assert.equal(
      cb,
      callbackTracker.handlers.get("number")?.get(cb)?.callback,
      "number callback not tracked",
    );
  });

  test("array keys teardown", () => {
    const callbackTracker = new CallbackTracker(
      sampleDatasetForCompletionCheck,
    );
    const cb = () => {};
    const teardown = callbackTracker.on(["string", "number"], cb);
    teardown();
    assert.notEqual(
      cb,
      callbackTracker.handlers.get("string")?.get(cb)?.callback,
      "string callback not removed after teardown",
    );
    assert.notEqual(
      cb,
      callbackTracker.handlers.get("number")?.get(cb)?.callback,
      "number callback not removed after teardown",
    );
  });

  test("array keys removal using 'off'", () => {
    const callbackTracker = new CallbackTracker(
      sampleDatasetForCompletionCheck,
    );
    const cb = () => {};
    callbackTracker.on(["string", "number"], cb);
    callbackTracker.off("number", cb);
    assert.equal(
      cb,
      callbackTracker.handlers.get("string")?.get(cb)?.callback,
      "string callback not removed after off",
    );
    assert.notEqual(
      cb,
      callbackTracker.handlers.get("number")?.get(cb)?.callback,
      "number callback not removed after 'off'",
    );
  });

  test("remove handlers record if not used", () => {
    const callbackTracker = new CallbackTracker(
      sampleDatasetForCompletionCheck,
    );
    const cb = () => {};
    callbackTracker.on(["string", "number"], cb)();
    assert.strictEqual(
      undefined,
      callbackTracker.handlers.get("string"),
      "did not remove string handlers record",
    );
    assert.strictEqual(
      undefined,
      callbackTracker.handlers.get("number"),
      "did not remove number handlers record",
    );
    callbackTracker.on("array", cb);
    callbackTracker.off("array", cb);
    assert.strictEqual(
      undefined,
      callbackTracker.handlers.get("array"),
      "did not remove array handlers record",
    );
  });

  test("get all handlers for key", () => {
    const callbackTracker = new CallbackTracker(
      sampleDatasetForCompletionCheck,
    );
    const cb1 = () => {};
    const cb2 = () => {};
    callbackTracker.on("string", cb1);
    callbackTracker.on("string", cb2);
    assert.deepEqual(
      callbackTracker.getCallbacks("string"),
      [cb1, cb2],
      "failed to retrive all tracked callbacks",
    );
    assert.strictEqual(
      callbackTracker.getCallbacks("number"),
      null,
      "number tracking should have no callbacks",
    );
  });
});
