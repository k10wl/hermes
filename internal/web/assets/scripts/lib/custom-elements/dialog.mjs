import {
  AssertInstance,
  AssertObject,
  AssertString,
} from "/assets/scripts/lib/assert.mjs";
import { Bind, html } from "/assets/scripts/lib/libdim.mjs";
import { ShortcutManager } from "/assets/scripts/lib/shortcut-manager.mjs";

export class HermesDialog extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    const dialog = new Bind((el) => AssertInstance.once(el, HTMLDialogElement));
    this.shadow.append(html`
      <style>
        :host {
          --duration: 100ms;
        }

        dialog[open] {
          opacity: 1;
        }

        dialog {
          background: transparent;
          border: none;
          outline: none;
          opacity: 0;
          padding: 0;
          transition: all var(--duration) allow-discrete;
          ${this.getAttribute("dialog-style") ?? ""}
        }

        dialog::backdrop {
          background: rgb(0 0 0 / 0.5);
        }

        @starting-style {
          dialog[open] {
            opacity: 0;
          }
        }
      </style>

      <dialog
        bind="${dialog}"
        onclose="${() => {
          this.dispatchEvent(
            new Event("close", { composed: true, bubbles: true }),
          );
        }}"
      >
        <slot></slot>
      </dialog>
    `);
    this.element = dialog.current;

    const showModal = this.element.showModal.bind(this.element);
    this.element.showModal = () => {
      this.dispatchEvent(new Event("show", { composed: true, bubbles: true }));
      showModal();
    };
  }

  connectedCallback() {
    this.element.addEventListener("click", (event) => {
      const rect = this.element.getBoundingClientRect();
      const isInDialog =
        rect.top <= event.clientY &&
        event.clientY <= rect.top + rect.height &&
        rect.left <= event.clientX &&
        event.clientX <= rect.left + rect.width;
      if (!isInDialog) {
        this.element.close();
      }
    });
  }
}

class HermesDialogCard extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" }).append(this.#html);
  }

  #html = html`
    <style>
      :host {
        --dialog-section-margin: 1.5rem;
      }

      section {
        max-width: 25rem;
        border-radius: 1rem;
        overflow: hidden;
        background: var(--bg-2);
        border: 1px solid rgb(from var(--text-0) r g b / 0.25);
        ${this.getAttribute("section-style") ?? ""};
      }
    </style>

    <section><slot></slot></section>
  `;
}

class HermesDialogTitle extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" }).append(this.#html);
  }

  #html = html`
    <style>
      h3 {
        margin: var(--dialog-section-margin);
      }
    </style>

    <h3><slot></slot></h3>
  `;
}

class HermesDialogBlock extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "open" }).append(this.#html);
  }

  #html = html`
    <style>
      div {
        margin: var(--dialog-section-margin);
      }
    </style>

    <div><slot></slot></div>
  `;
}

customElements.define("h-dialog", HermesDialog);
customElements.define("h-dialog-card", HermesDialogCard);
customElements.define("h-dialog-title", HermesDialogTitle);
customElements.define("h-dialog-block", HermesDialogBlock);

export class AlertDialog extends HTMLElement {
  titleSlot = new Bind((el) => AssertInstance.once(el, HTMLSlotElement));
  descriptionSlot = new Bind((el) => AssertInstance.once(el, HTMLSlotElement));
  element = new Bind((el) => AssertInstance.once(el, HermesDialog));

  constructor() {
    super();
    this.attachShadow({ mode: "open" }).append(html`
      <style>
        * {
          color: var(--text-0);
        }
        h-button {
          display: block;
          &::part(button) {
            width: 100%;
            font-size: 1.1rem;
          }
        }
      </style>

      <h-dialog bind="${this.element}">
        <h-dialog-card>
          <h-dialog-title>
            <slot bind="${this.titleSlot}" name="title"
              >Something went wrong</slot
            >
          </h-dialog-title>
          <h-dialog-block>
            <slot bind="${this.descriptionSlot}" name="description"
              >Unexpected error occurred, but there are no details</slot
            >
          </h-dialog-block>
          <h-dialog-block id="actions">
            <h-button
              variant="primary"
              onclick="${() => this.element.current.element.close()}"
            >
              OK <h-key>enter / y</h-key>
            </h-button>
          </h-dialog-block>
        </h-dialog-card>
      </h-dialog>
    `);
  }

  /**
   * @param {{title?: string, description?: string}} [info]
   */
  alert(info) {
    const preupdate = {
      titleSlot: this.titleSlot.current.textContent,
      descriptionSlot: this.descriptionSlot.current.textContent,
    };
    if (info?.title) {
      this.titleSlot.current.textContent = info.title;
    }
    if (info?.description) {
      this.descriptionSlot.current.textContent = info.description;
    }
    this.element.current.element.showModal();
    const off = ShortcutManager.keydown(
      ["<Enter>", "<KeyY>"],
      (event) => {
        this.element.current.element.close();
        event.preventDefault();
        event.stopPropagation();
      },
      { priority: 99999 },
    );
    this.element.current.addEventListener(
      "close",
      () => {
        off();
        this.element.current.element.addEventListener(
          "animationend",
          () => {
            this.titleSlot.current.textContent = preupdate.titleSlot;
            this.descriptionSlot.current.textContent =
              preupdate.descriptionSlot;
          },
          { once: true },
        );
      },
      { once: true },
    );
  }

  /** @type {AlertDialog} */
  static instance;
  connectedCallback() {
    if (AlertDialog.instance) {
      throw new Error("Only one context menu is allowed to exist");
    }
    AlertDialog.instance = this;
  }
}

customElements.define("h-alert-dialog", AlertDialog);

export class ConfirmDialog extends HTMLElement {
  dialog = new Bind((el) => AssertInstance.once(el, HermesDialog));
  titleSlot = new Bind((el) => AssertInstance.once(el, HTMLSlotElement));
  descriptionSlot = new Bind((el) => AssertInstance.once(el, HTMLSlotElement));

  constructor() {
    super();
    this.attachShadow({ mode: "open" }).append(html`
      <style>
        * {
          color: var(--text-0);
        }
        h-button {
          display: block;
          &::part(button) {
            width: 100%;
          }
        }
        #actions {
          display: flex;
          flex-wrap: wrap;
          justify-content: flex-end;
          gap: 0.5rem;
          h-button::part(button) {
            font-size: 1.1rem;
          }
        }
      </style>

      <h-dialog bind="${this.dialog}">
        <h-dialog-card>
          <h-dialog-title>
            <slot bind="${this.titleSlot}" name="title"></slot>
          </h-dialog-title>
          <h-dialog-block>
            <slot bind="${this.descriptionSlot}" name="description"></slot>
          </h-dialog-block>
          <h-dialog-block>
            <div id="actions">
              <h-button
                variant="error"
                onclick="${() => this.dialog.current.element.close()}"
              >
                Cancel <h-key>n</h-key>
              </h-button>
              <h-button
                variant="primary"
                onclick="${() => this.dialog.current.element.close()}"
              >
                Confirm <h-key>↵/y</h-key>
              </h-button>
            </div>
          </h-dialog-block>
        </h-dialog-card>
      </h-dialog>
    `);
  }

  #assertion = new AssertObject(
    {
      title: AssertString,
      description: AssertString,
    },
    "confirm modal is expected to render with both title and decription",
  );
  /**
   * @param {{title: string, description: string}} info
   * @returns {Promise<boolean>}
   */
  confirm = async (info) => {
    this.#assertion.check(info);
    /** @type {PromiseWithResolvers<boolean>} */
    const { promise, resolve } = Promise.withResolvers();
    const preupdate = {
      titleSlot: this.titleSlot.current.textContent,
      descriptionSlot: this.descriptionSlot.current.textContent,
    };
    if (info.title) {
      this.titleSlot.current.textContent = info.title;
    }
    if (info.description) {
      this.descriptionSlot.current.textContent = info.description;
    }
    this.dialog.current.element.showModal();
    const confirmKeys = ShortcutManager.keydown(
      ["<Enter>", "<KeyY>"],
      (event) => {
        this.dialog.current.element.close();
        event.preventDefault();
        event.stopPropagation();
        resolve(true);
      },
      { priority: 99999 },
    );
    const cancelKeys = ShortcutManager.keydown(
      ["<Escape>", "<KeyN>"],
      (event) => {
        this.dialog.current.element.close();
        event.preventDefault();
        event.stopPropagation();
        resolve(false);
      },
    );
    this.dialog.current.addEventListener(
      "close",
      () => {
        confirmKeys();
        cancelKeys();
        resolve(false);
        this.dialog.current.element.addEventListener(
          "animationend",
          () => {
            this.titleSlot.current.textContent = preupdate.titleSlot;
            this.descriptionSlot.current.textContent =
              preupdate.descriptionSlot;
          },
          { once: true },
        );
      },
      { once: true },
    );
    return promise;
  };

  /** @type {ConfirmDialog} */
  static instance;
  connectedCallback() {
    if (ConfirmDialog.instance) {
      throw new Error("Only one context menu is allowed to exist");
    }
    ConfirmDialog.instance = this;
  }
}

customElements.define("h-confirm-dialog", ConfirmDialog);
