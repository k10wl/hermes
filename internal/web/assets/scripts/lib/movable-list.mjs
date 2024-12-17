import { Publisher } from "./publisher.mjs";

export class MovableList {
  #cursor = new Publisher(0);

  #elementsParent;

  /**
   * @param {HTMLElement} elementsParent
   * @param {(current: number, previous: number) => void} notify
   */
  constructor(elementsParent, notify) {
    this.#elementsParent = elementsParent;
    this.#cursor.attach({ notify: notify });
  }

  /**
   * @param {-1 | 1} direction
   */
  move(direction) {
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
