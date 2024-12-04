import { ValidateString } from "/assets/scripts/lib/validate.mjs";

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
    window.addEventListener("resize", this.autoresize);
    resizeObserver.observe(this);
    this.#textAreaAutoresizeCleanup.push(resizeObserver.disconnect, () => {
      this.removeEventListener("input", this.autoresize);
      window.removeEventListener("resize", this.autoresize);
    });
  }

  disconnectedCallback() {
    this.#textAreaAutoresizeCleanup.forEach((cb) => cb());
  }

  autoresize() {
    this.style.height = "auto";
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
    window.addEventListener("keydown", this.focusOnKeydown);
    window.addEventListener("paste", this.focusOnPaste);
    this.#textAreaAutoresizeCleanup.push(() => {
      window.removeEventListener("keydown", this.focusOnKeydown);
      window.removeEventListener("paste", this.focusOnPaste);
    });
  }

  /** @param {KeyboardEvent} e  */
  focusOnKeydown(e) {
    if (
      document.activeElement === this ||
      document.activeElement === null ||
      document.activeElement.tagName === "TEXTAREA" ||
      document.activeElement.tagName === "INPUT" ||
      e.shiftKey ||
      e.altKey ||
      e.metaKey ||
      e.ctrlKey ||
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
      this.value += ValidateString.parse(e.clipboardData?.getData("text"));
      this.focus();
    } catch {
      // whatever, just don't explode
    }
  }
}
