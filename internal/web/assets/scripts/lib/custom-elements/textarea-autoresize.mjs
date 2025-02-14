import { AssertString } from "/assets/scripts/lib/assert.mjs";

import { ShortcutManager } from "../shortcut-manager.mjs";

export class TextAreaAutoresize extends HTMLTextAreaElement {
  /** @type {(()=>void)[]}*/
  #textAreaAutoresizeCleanup = [];

  /** @type {number} */
  #selectionStart = this.selectionStart;
  /** @type {number} */
  #selectionEnd = this.selectionEnd;

  constructor() {
    super();
  }

  connectedCallback() {
    this.autoresize = this.autoresize.bind(this);
    this.focusOnKeydown = this.focusOnKeydown.bind(this);
    this.#limitHeight();
    this.#focusOnInput();
    this.autoresize();
    const resizeObserver = new ResizeObserver(this.autoresize);
    this.addEventListener("input", this.autoresize);
    this.addEventListener("change", this.autoresize);
    window.addEventListener("resize", this.autoresize);
    resizeObserver.observe(this);
    this.#textAreaAutoresizeCleanup.push(
      () => resizeObserver.disconnect(),
      () => {
        this.removeEventListener("input", this.autoresize);
        this.removeEventListener("change", this.autoresize);
        window.removeEventListener("resize", this.autoresize);
      },
    );
  }

  disconnectedCallback() {
    this.#textAreaAutoresizeCleanup.forEach((cb) => cb());
  }

  autoresize() {
    this.style.height = "0px";
    this.style.height = this.scrollHeight + "px";
  }

  #limitHeight() {
    const maxRowsAttribute = this.getAttribute("max-rows");
    if (!maxRowsAttribute) {
      return;
    }
    const maxRows = Number.parseInt(maxRowsAttribute, 10);
    if (Number.isNaN(maxRows)) {
      return;
    }
  }

  #focusOnInput() {
    const focusOnInput = this.getAttribute("focus-on-input");
    if (typeof focusOnInput === "undefined") {
      return;
    }
    this.addEventListener("focus", () => {
      window.removeEventListener("paste", this.focusOnPaste);
    });
    this.addEventListener("blur", () => {
      this.#selectionStart = this.selectionStart;
      this.#selectionEnd = this.selectionEnd;
      window.addEventListener("paste", this.focusOnPaste);
    });
    this.#textAreaAutoresizeCleanup.push(
      ShortcutManager.keydown("<*>", this.focusOnKeydown),
      () => {
        window.removeEventListener("paste", this.focusOnPaste);
      },
    );
  }

  /** @param {KeyboardEvent} e  */
  focusOnKeydown(e) {
    const target = ShortcutManager.getTarget(e);
    if (
      target === this ||
      target === null ||
      document.activeElement?.tagName === "TEXTAREA" ||
      document.activeElement?.tagName === "INPUT" ||
      e.shiftKey ||
      e.altKey ||
      e.metaKey ||
      e.ctrlKey ||
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
    this.focus();
  }

  /** @param {ClipboardEvent} e  */
  focusOnPaste = (e) => {
    try {
      e.stopPropagation();
      e.preventDefault();
      const text = AssertString.check(e.clipboardData?.getData("text"));
      this.value =
        this.value.substring(0, this.#selectionStart) +
        text +
        this.value.substring(this.#selectionEnd);
      this.focus();
      const adjustedEnd = this.#selectionEnd + text.length;
      this.setSelectionRange(adjustedEnd, adjustedEnd);
    } catch {
      // whatever, just don't explode
    }
  };
}
