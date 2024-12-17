import { AssertInstance } from "/assets/scripts/lib/assert.mjs";

import { ShortcutManager } from "../shortcut-manager.mjs";
import { TextAreaAutoresize } from "./textarea-autoresize.mjs";

export class Form extends HTMLFormElement {
  /** @type (() => void)[] */
  #messageContentFormCleanup = [];

  constructor() {
    super();
    this.detectKeyboardSubmit = this.detectKeyboardSubmit.bind(this);
  }

  connectedCallback() {
    this.#messageContentFormCleanup.push(
      ShortcutManager.keydown("<Enter>", this.detectKeyboardSubmit),
    );
  }

  reset() {
    super.reset();
    this.querySelectorAll("textarea").forEach((el) => {
      try {
        AssertInstance.once(el, TextAreaAutoresize).autoresize();
      } catch {
        // just don't explode
      }
    });
  }

  disconnectedCallback() {
    this.#messageContentFormCleanup.forEach((cb) => cb());
  }

  /** @param {KeyboardEvent} e */
  detectKeyboardSubmit(e) {
    if (e.shiftKey || e.metaKey || e.ctrlKey) {
      return;
    }
    e.preventDefault();
    e.stopPropagation();
    this.requestSubmit();
  }
}

customElements.define("hermes-form", Form, { extends: "form" });
