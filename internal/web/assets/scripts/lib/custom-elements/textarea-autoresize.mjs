import { AssertString } from "/assets/scripts/lib/assert.mjs";

import { ShortcutManager } from "../shortcut-manager.mjs";

export class TextAreaAutoresize extends HTMLTextAreaElement {
  /** @type {(()=>void)[]}*/
  #textAreaAutoresizeCleanup = [];

  constructor() {
    super();
  }

  connectedCallback() {
    this.autoresize = this.autoresize.bind(this);
    this.focusOnKeydown = this.focusOnKeydown.bind(this);
    this.focusOnPaste = this.focusOnPaste.bind(this);
    this.#limitHeight();
    this.#focusOnInput();
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
    this.style.setProperty("max-height", `${maxRows}lh`);
  }

  #focusOnInput() {
    const focusOnInput = this.getAttribute("focus-on-input");
    if (typeof focusOnInput === "undefined") {
      return;
    }
    window.addEventListener("paste", this.focusOnPaste);
    this.#textAreaAutoresizeCleanup.push(
      ShortcutManager.keydown("<*>", this.focusOnKeydown),
      () => {
        window.removeEventListener("paste", this.focusOnPaste);
      },
    );
  }

  /** @param {KeyboardEvent} e  */
  focusOnKeydown(e) {
    if (
      e.target === this ||
      e.target === null ||
      document.activeElement?.tagName === "TEXTAREA" ||
      document.activeElement?.tagName === "INPUT" ||
      e.shiftKey ||
      e.altKey ||
      e.metaKey ||
      e.ctrlKey ||
      "Escape" === e.key ||
      "Enter" === e.key ||
      "Tab" === e.key ||
      "ArrowLeft" === e.key ||
      "ArrowRight" === e.key ||
      "ArrowTop" === e.key ||
      "ArrowBottom" === e.key ||
      " " === e.key
    ) {
      return;
    }
    this.focus();
  }

  /** @param {ClipboardEvent} e  */
  focusOnPaste(e) {
    if (window.document.activeElement === this) {
      return;
    }
    try {
      this.value += AssertString.check(e.clipboardData?.getData("text"));
      this.focus();
    } catch {
      // whatever, just don't explode
    }
  }
}
