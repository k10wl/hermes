import { AssertInstance } from "../assert.mjs";
import { Bind, html } from "../libdim.mjs";

export class ResizableTextInput extends HTMLElement {
  #value = "";
  #content = new Bind((el) => AssertInstance.once(el, HTMLDivElement));

  constructor() {
    super();
  }

  connectedCallback() {
    this.attachShadow({ mode: "open" }).append(html`
      <style>
        * {
          box-sizing: border-box;
        }
        #wrapper {
          position: relative;
        }
        #content {
          overflow: auto;
        }
        #placeholder {
          pointer-events: none;
          display: none;
          margin: 0;
          position: absolute;
          left: 0;
          top: 0;
          opacity: 0.5;
        }
        :host([data-empty]) #placeholder {
            display: block;
          }
        }
      </style>

      <div id="wrapper" part="wrapper">
        <div
          id="content"
          part="content"
          role="textbox"
          contenteditable="plaintext-only"
          bind="${this.#content}"
          oninput="${() => {
            this.#value = this.#content.current.innerText;
            this.#updateEmptyState();
            this.#dispatchChangeEvent();
          }}"
        ></div>
        <div id="placeholder" part="placeholder">
          ${this.getAttribute("placeholder")}
        </div>
      </div>
    `);
    // formatter creates unwanted newline
    this.#content.current.innerHTML = "<br>";
    this.#updateEmptyState();
  }

  #updateEmptyState() {
    const empty = this.#value === "" || this.#value === "\n";
    if (empty) {
      this.setAttribute("data-empty", "true");
      return;
    }
    this.removeAttribute("data-empty");
  }

  #dispatchChangeEvent() {
    this.dispatchEvent(
      new CustomEvent("change", {
        detail: {
          value: this.#value,
        },
      }),
    );
  }

  get value() {
    return this.#value;
  }

  set value(value) {
    this.#value = value;
    this.#content.current.innerText = value;
    this.#updateEmptyState();
    this.#dispatchChangeEvent();
  }

  get content() {
    return this.#content.current;
  }
}

customElements.define("h-resizable-text-input", ResizableTextInput);
