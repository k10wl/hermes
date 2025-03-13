import { html } from "/assets/scripts/lib/libdim.mjs";

class Card extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
  }

  connectedCallback() {
    this.shadow.append(this.#html);
  }

  #html = html`
    <style>
      div {
        box-sizing: border-box;
        color: var(--text-0);
        padding: 0.5rem 1rem;
        border-radius: 0.5rem;
        border: 1px solid rgb(from var(--text-0) r g b / 0.25);
        transition: border-color var(--color-transition-duration);
      }
      :host([data-interactive]:hover) div {
        border-color: var(--primary);
      }
    </style>

    <div part="card"><slot></slot></div>
  `;
}

customElements.define("h-card", Card);
