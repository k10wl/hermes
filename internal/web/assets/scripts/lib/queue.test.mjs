import * as assert from "node:assert";
import { describe } from "node:test";

import { Queue } from "./queue.mjs";

describe("queue test", () => {
  /** @type {Queue<{data: string}>} */
  const queue = new Queue();
  const obj1 = { data: "first" };
  const obj2 = { data: "second" };
  const obj3 = { data: "three" };
  for (let i = 0; i < 3; i++) {
    queue.enqueue(obj1);
    assert.equal(queue.size, 1, "expected length of 1");
    queue.enqueue(obj2);
    assert.equal(queue.size, 2, "expected length of 2");
    queue.enqueue(obj3);
    assert.equal(queue.size, 3, "expected length of 3");
    assert.strictEqual(queue.peek(), obj1, "peek expected obj1 to be in head");
    assert.strictEqual(
      queue.dequeue(),
      obj1,
      "dequeue expected obj1 to be in head",
    );
    assert.equal(queue.size, 2, "expected length of 2");
    assert.strictEqual(queue.peek(), obj2, "peek expected obj2 to be in head");
    assert.strictEqual(
      queue.dequeue(),
      obj2,
      "dequeue expected obj2 to be in head",
    );
    assert.equal(queue.size, 1, "expected length of 1");
    assert.strictEqual(queue.peek(), obj3, "peek expected obj3 to be in head");
    assert.strictEqual(
      queue.dequeue(),
      obj3,
      "dequeue expected obj3 to be in head",
    );
    assert.strictEqual(
      queue.peek(),
      null,
      "peek expected queue head to be null",
    );
    assert.strictEqual(
      queue.dequeue(),
      null,
      "dequeue expected queue head to be null",
    );
    assert.equal(queue.size, 0, "expected length of 0");
  }
});
