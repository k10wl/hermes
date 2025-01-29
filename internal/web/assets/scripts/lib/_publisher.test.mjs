import { describe, test } from "node:test";

import * as assert from "assert";

import { Publisher } from "./publisher.mjs";

describe("Publisher", () => {
  test("should subscribe elements", () => {
    const initial = 0;
    const publisher = new Publisher(initial);
    assert.equal(
      publisher.value,
      initial,
      "expected to be initialized with given value",
    );

    const subscriber = new Subscriber();
    publisher.subscribe(subscriber);
    assert.equal(
      publisher.subscribers[0],
      subscriber,
      "expected subscriber to become part of publisher",
    );

    publisher.update(1);
    assert.equal(publisher.value, 1, "expected publisher value update");
    assert.deepEqual(
      subscriber.args,
      [1],
      "expected observer to receive new value upon update",
    );
    publisher.unsubscribe(subscriber);
    assert.equal(
      publisher.subscribers.length,
      0,
      "expected empty subscribers after unsub",
    );

    const unsubscribe = publisher.subscribe(subscriber);
    assert.equal(
      publisher.subscribers[0],
      subscriber,
      "expected subscriber to become part of publisher",
    );
    unsubscribe();
    assert.equal(
      publisher.subscribers.length,
      0,
      "expected empty subscribers after unsub",
    );
  });
});

class Subscriber {
  /**
   * @param {unknown[]} args
   */
  notify(...args) {
    this.args = args;
  }
}
