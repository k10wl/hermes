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
   *
   * Examples:
   * - <M-KeyP>    - meta/alt + p
   * - <C-ArrowUp> - ctrl + arrow up
   * - <C-S-KeyA>  - ctrl + shift + a
   * - <KeyT>      - t
   */

  /** @typedef {Record<ShortcutNotation, KeyboardEvent>} Structure */

  static #keydownTracker = new CallbackTracker(/** @type {Structure} */ ({}));

  /**
   * @type {CallbackTracker<Structure>["on"]}
   */
  static keydown(key, callback) {
    return ShortcutManager.#keydownTracker.on(key, callback);
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

  static __init__ = (() => {
    if (typeof window === "undefined") {
      return; // for node tests
    }
    window.addEventListener("keydown", (event) => {
      const callbacks = ShortcutManager.#keydownTracker.getCallbacks(
        ShortcutManager.eventToNotation(event),
      );
      if (!callbacks) {
        return;
      }
      let shouldBreak = false;
      for (const callback of callbacks) {
        const originalStopPropagation = event.stopPropagation.bind(event);
        event.stopPropagation = () => {
          shouldBreak = true;
          originalStopPropagation();
        };
        callback(event);
        if (shouldBreak) {
          break;
        }
      }
    });
  })();
}
