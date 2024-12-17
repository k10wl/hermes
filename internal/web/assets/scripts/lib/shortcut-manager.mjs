// TODO chaining keys will be useful for closing stuff "<Escape> <Escape>"

import { AssertInstance } from "./assert.mjs";
import { CallbackTracker } from "./callback-tracker.mjs";

export class ShortcutManager {
  /**
   * @typedef {string} ShortcutNotation
   *
   * Keys must follow specific notation <M-C-S-KeyA>, where:
   * - M is optional and indicates Meta/Alt
   * - C is optional and indicates Ctrl
   * - S is optional and indicates Shift
   * - KeyA is mandatory and indicates key code - {@link https://developer.mozilla.org/en-US/docs/Web/API/KeyboardEvent/code}
   * - * indicates any key
   *
   * Examples:
   * - <M-KeyP>    - meta/alt + p
   * - <C-ArrowUp> - ctrl + arrow up
   * - <C-S-KeyA>  - ctrl + shift + a
   * - <KeyT>      - t
   * - <*>         - on any key
   */

  /** @typedef {KeyboardEvent & {notation: string}} KeyboardEventWithNotation */
  /** @typedef {Record<ShortcutNotation, KeyboardEventWithNotation>} Structure */

  static #keydownTracker = new CallbackTracker(/** @type {Structure} */ ({}));

  /**
   * @type {CallbackTracker<Structure>["on"]}
   */
  static keydown(...args) {
    return ShortcutManager.#keydownTracker.on(...args);
  }

  /**
   * @param {KeyboardEvent} event
   * @returns {ShortcutNotation}
   */
  static eventToNotation(event) {
    let modifiers = "";
    if (event.metaKey || event.altKey) {
      modifiers += "M-";
    }
    if (event.ctrlKey) {
      modifiers += "C-";
    }
    if (event.shiftKey) {
      modifiers += "S-";
    }
    return `<${modifiers}${event.code}>`;
  }

  /**
   * @param {KeyboardEventWithNotation} event
   * @returns {((event: KeyboardEventWithNotation) => void)[]}
   */
  static #getKeydownCallbacks(event) {
    const keySpecificCallbacks = ShortcutManager.#keydownTracker.getCallbacks(
      event.notation,
      "<*>",
    );
    return keySpecificCallbacks;
  }

  static #addKeydownListener() {
    window.addEventListener("keydown", (event) => {
      const eventWithNotation = Object.assign(event, {
        notation: ShortcutManager.eventToNotation(event),
      });
      const callbacks = ShortcutManager.#getKeydownCallbacks(eventWithNotation);
      if (callbacks.length === 0) {
        return;
      }
      let shouldBreak = false;
      for (const callback of callbacks) {
        const originalStopPropagation = event.stopPropagation.bind(event);
        eventWithNotation.stopPropagation = () => {
          shouldBreak = true;
          originalStopPropagation();
        };
        callback(eventWithNotation);
        if (shouldBreak) {
          break;
        }
      }
    });
  }

  /**
   * @param {KeyboardEvent} event
   * @returns {HTMLElement | null}
   */
  static getTarget(event) {
    try {
      return AssertInstance.once(
        // @ts-expect-error - https://developer.mozilla.org/en-US/docs/Web/API/Event/explicitOriginalTarget
        event.explicitOriginalTarget ?? event.target,
        HTMLElement,
      );
    } catch {
      return null;
    }
  }

  static __init__ = (() => {
    if (typeof window === "undefined") {
      return; // for node tests
    }
    ShortcutManager.#addKeydownListener();
  })();
}
