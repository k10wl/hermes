/** @template T */
export class Queue {
  /** @type {Node<T> | null} */
  #head = null;
  /** @type {Node<T> | null} */
  #tail = null;
  size = 0;

  constructor() {}

  /** @param {T} value  */
  enqueue(value) {
    this.size++;
    const node = new Node(value);
    if (!this.#tail) {
      this.#tail = this.#head = node;
      return;
    }
    this.#tail.next = node;
    node.prev = this.#tail;
    this.#tail = node;
  }

  /** @returns {T | null} */
  peek() {
    return this.#head?.value ?? null;
  }

  /** @returns {T | null} */
  dequeue() {
    this.size = Math.max(0, this.size - 1);
    const value = this.#head?.value;
    if (this.#head?.next) {
      this.#head = this.#head.next;
    } else {
      this.#tail = this.#head = null;
    }
    return value ?? null;
  }
}

/** @template T */
class Node {
  /** @type {Node<T> | null} */
  prev = null;
  /** @type {Node<T> | null} */
  next = null;
  /** @param {T} value  */
  constructor(value) {
    this.value = value;
  }
}
