import { AssertInstance } from "../assert.mjs";
import { html } from "../html-v2.mjs";
import { ShortcutManager } from "../shortcut-manager.mjs";

export class HermesDialog extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    let dialog;
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

        @starting-style {
          dialog[open] {
            opacity: 0;
          }
        }
      </style>

      <dialog
        bind="${(/** @type {unknown} */ e) => (dialog = e)}"
        onclose="${() => {
          this.dispatchEvent(
            new Event("close", { composed: true, bubbles: true }),
          );
        }}"
      >
        <slot></slot>
      </dialog>
    `);
    this.element = AssertInstance.once(dialog, HTMLDialogElement);

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

export class HermesAlertDialog extends HTMLElement {
  constructor() {
    super();
    let dialog;
    this.attachShadow({ mode: "open" }).append(html`
      <style>
        * {
          color: var(--text-0);
          text-wrap: balance;
        }
        h-button {
          display: block;
          &::part(button) {
            width: 100%;
          }
        }
      </style>

      <h-dialog
        bind="${(/** @type {unknown} */ element) => (dialog = element)}"
      >
        <h-dialog-card>
          <h-dialog-title>
            <slot
              bind="${(/** @type {unknown} */ element) => {
                this.titleSlot = AssertInstance.once(element, HTMLSlotElement);
              }}"
              name="title"
              >Something went wrong</slot
            >
          </h-dialog-title>
          <h-dialog-block>
            <slot
              bind="${(/** @type {unknown} */ element) => {
                this.descriptionSlot = AssertInstance.once(
                  element,
                  HTMLSlotElement,
                );
              }}"
              name="description"
              >Unexpected error occurred, but there are no details</slot
            >
          </h-dialog-block>
          <h-dialog-block id="actions">
            <h-button
              variant="primary"
              onclick="${() => this.dialog.element.close()}"
            >
              OK <h-key>enter / y</h-key>
            </h-button>
          </h-dialog-block>
        </h-dialog-card>
      </h-dialog>
    `);

    this.dialog = AssertInstance.once(
      dialog,
      HermesDialog,
      "alert should be a dialog",
    );
  }

  /**
   * @param {{title?: string, description?: string}} [info]
   */
  alert(info) {
    const preupdate = {
      titleSlot: this.titleSlot.textContent,
      descriptionSlot: this.descriptionSlot.textContent,
    };
    if (info?.title) {
      this.titleSlot.textContent = info.title;
    }
    if (info?.description) {
      this.descriptionSlot.textContent = info.description;
    }
    this.dialog.element.showModal();
    const off = ShortcutManager.keydown(["<Enter>", "<KeyY>"], (event) => {
      this.dialog.element.close();
      event.preventDefault();
      event.stopPropagation();
    });
    this.dialog.addEventListener(
      "close",
      () => {
        off();
        this.titleSlot.textContent = preupdate.titleSlot;
        this.descriptionSlot.textContent = preupdate.descriptionSlot;
      },
      { once: true },
    );
  }
}
customElements.define("h-alert-dialog", HermesAlertDialog);
