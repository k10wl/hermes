import { Publisher } from "./publisher.mjs";

export class MovableList {
  #cursor = new Publisher(0);
  /** @type {number | undefined} */
  previous = undefined;

  #elementsParent;
  /** @type {number | null} */

  /**
   * @param {HTMLElement} elementsParent
   * @param {(current: number, previous?: number) => void} notify
   */
  constructor(elementsParent, notify) {
    this.#elementsParent = elementsParent;
    this.#cursor.subscribe({
      notify: (current) => notify(current, this.previous),
    });
  }

  /**
   * @param {-1 | 1} direction
   */
  move(direction) {
    this.previous = this.cursor;
    this.#cursor.update((value) => {
      const newValue = value + direction;
      const bounds = this.#elementsParent.children.length;
      if (newValue < 0) {
        return bounds - 1;
      }
      if (newValue >= bounds) {
        return 0;
      }
      return newValue;
    });
  }

  /** @param {number} update */
  set cursor(update) {
    this.previous = this.cursor;
    this.#cursor.update(update);
    this.#cursor.notify();
  }

  /** @returns {number} */
  get cursor() {
    return this.#cursor.value;
  }

  notify() {
    this.#cursor.notify();
  }
}
