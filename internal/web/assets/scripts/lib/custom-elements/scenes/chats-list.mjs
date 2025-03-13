import { html } from "/assets/scripts/lib/libdim.mjs";

export class ChatsListScene extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "closed" }).replaceChildren(html`
      <style>
        main {
          height: 100%;
          display: grid;
          place-items: center;
          overflow: auto;
        }

        h-chats {
          max-width: var(--container-max-width);
          margin: var(--container-margin);
          width: 100%;
        }
      </style>

      <main>
        <h-chats></h-chats>
      </main>
    `);
  }
}

customElements.define("hermes-chats-list-scene", ChatsListScene);
