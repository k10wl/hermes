import { html } from "../../html.mjs";

export class ChatsListScene extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "closed" }).innerHTML = html`
      <style>
        main {
          height: 100%;
          display: grid;
          place-items: center;
          padding: 16px;
          overflow: auto;
        }

        hermes-chats {
          max-width: var(--container);
          width: 100%;
        }

        a {
          margin: 0.25rem 0;
          color: var(--text-0);
          display: block;
          padding: 0.5rem 1rem;
          border-radius: 0.5rem;
          border: 1px solid rgb(from var(--text-0) r g b / 0.25);
          text-decoration: none;
          transition: border-color var(--color-transition-duration);
        }

        a:hover {
          border-color: var(--primary);
        }
      </style>

      <main>
        <hermes-chats></hermes-chats>
      </main>
    `;
  }
}

customElements.define("hermes-chats-list-scene", ChatsListScene);
