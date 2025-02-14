import { html } from "../html-v2.mjs";

customElements.define(
  "h-key",
  class extends HTMLElement {
    constructor() {
      super();
      this.attachShadow({ mode: "open" }).append(
        html`<style>
            code {
              user-select: none;
              color: rgb(from var(--text-0) r g b / 0.25);
              &:before {
                content: "[";
              }
              &:after {
                content: "]";
              }
            }</style
          ><code><slot></slot></code>`,
      );
    }
  },
);
