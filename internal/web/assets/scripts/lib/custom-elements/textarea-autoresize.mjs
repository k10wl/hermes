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
    if (typeof this.getAttribute("focus-on-input") !== "undefined") {
      this.#textAreaAutoresizeCleanup.push(
        ShortcutManager.keydown("<*>", this.focusOnKeydown),
      );
    }
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
}
