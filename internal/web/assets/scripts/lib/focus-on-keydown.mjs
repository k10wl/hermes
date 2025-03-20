import { ShortcutManager } from "./shortcut-manager.mjs";

export class FocusOnKeydown {
  /** @type {HTMLElement | null} */
  #el = null;

  constructor() {
    this.detach = ShortcutManager.keydown("<*>", this.#focus);
  }

  /** @param {HTMLElement} el */
  attach(el) {
    this.#el = el;
  }

  /** @param {KeyboardEvent} e  */
  #focus = (e) => {
    const target = ShortcutManager.getTarget(e);
    if (
      target === this.#el ||
      target === null ||
      document.activeElement?.tagName === "TEXTAREA" ||
      document.activeElement?.tagName === "INPUT" ||
      e.altKey ||
      e.metaKey ||
      (e.shiftKey && e.code.startsWith("Shift")) ||
      (e.ctrlKey && "KeyV" !== e.code) ||
      "Escape" === e.code ||
      "Enter" === e.code ||
      "Tab" === e.code ||
      "ArrowLeft" === e.code ||
      "ArrowRight" === e.code ||
      "ArrowUp" === e.code ||
      "ArrowDown" === e.code ||
      "PageUp" === e.code ||
      "PageDown" === e.code ||
      "Home" === e.code ||
      "End" === e.code ||
      "Space" === e.code
    ) {
      return;
    }
    this.#el?.focus();
  };
}
