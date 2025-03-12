import { html } from "/assets/scripts/lib/libdim.mjs";

export class HermesButton extends HTMLElement {
  constructor() {
    super();
    const variant = this.getAttribute("variant");
    const shadow = this.attachShadow({ mode: "open" });
    shadow.append(html`
      <style>
        :host {
          --h-button-color: var(
            ${(() => {
              switch (variant) {
                case "primary":
                  return "--primary";
                case "error":
                  return "--color-error";
                default:
                  return "--text-0";
              }
            })()}
          );
        }
        button {
          cursor: pointer;
          color: var(--h-button-color);
          border-color: transparent;
          outline-color: transparent;
          background-color: transparent;
          transition: all var(--color-transition-duration);
          border-radius: 0.25rem;
          user-select: none;
          &:hover {
            background-color: rgb(from var(--h-button-color) r g b / 0.08);
          }
        }
      </style>
      <button part="button">
        <slot></slot>
      </button>
    `);
  }
}

customElements.define("h-button", HermesButton);
