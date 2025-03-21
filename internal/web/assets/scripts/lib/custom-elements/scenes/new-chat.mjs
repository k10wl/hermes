import { html } from "/assets/scripts/lib/libdim.mjs";

import { controlPalanelVisibility } from "../control-panel.mjs";

const HERMES_SHORTCUTS_TAG_NAME = "hermes-new-chat-shortcuts";
customElements.define(
  HERMES_SHORTCUTS_TAG_NAME,
  class Shortcuts extends HTMLElement {
    /** @type {{name: string, key: string, onclick: () => void}[]} */
    static #list = [
      {
        name: "Control Panel",
        key: "Ctrl-P",
        onclick: () => {
          controlPalanelVisibility.update(true);
        },
      },
    ];

    constructor() {
      super();
      this.shadow = this.attachShadow({ mode: "closed" });
    }

    connectedCallback() {
      this.shadow.append(html`
        <style>
          :host {
            --default-color: rgb(from var(--text-0) r g b / 0.33);
            --hover-color: var(--text-0);
          }

          ul {
            display: grid;
            gap: 0.25rem;
            color: var(--default-color);
            font-size: 0.75rem;
            padding: 0;
            margin: 0;
          }

          li {
            list-style: none;
            display: grid;
            grid-template-columns: 1fr 1fr;
            align-items: center;
            gap: 0.5rem;
          }

          span {
            text-align: end;
            margin-left: auto;
          }

          button {
            --color: var(--default-color);
            background: transparent;
            outline-color: transparent;
            border-color: transparent;
            cursor: pointer;
            border: 1px solid var(--color);
            padding: 0.1rem 0.2rem;
            --transition: var(--color-transition-duration);
            transition:
              color var(--transition),
              border-color var(--transition);
            border-radius: 0.2rem;
            margin: 0;
            color: var(--color);
            &:hover {
              --color: var(--hover-color);
            }
          }
        </style>

        <ul>
          ${Shortcuts.#list.map(
            (shortcut) => html`
              <li>
                <span>${shortcut.name}</span>
                <div>
                  <button onclick="${() => shortcut.onclick()}">
                    ${shortcut.key}
                  </button>
                </div>
              </li>
            `,
          )}
        </ul>
      `);
    }
  },
);

export class CreateChatScene extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    this.attachShadow({ mode: "closed" }).append(html`
      <style>
        * {
          box-sizing: border-box;
        }

        main {
          height: 100%;
          display: grid;
          place-items: center;
        }

        div {
          width: 100%;
          max-width: var(--container-max-width);
          margin: var(--container-margin);
          display: flex;
          justify-content: center;
          align-items: center;
          flex-direction: column;
          gap: 1rem;
        }

        hermes-message-form {
          width: 100%;
        }
      </style>

      <main>
        <div>
          <hermes-message-form
            placeholder="What do you want to know?"
          ></hermes-message-form>
          <hermes-new-chat-shortcuts></hermes-new-chat-shortcuts>
        </div>
      </main>
    `);
  }
}

customElements.define("hermes-new-chat-scene", CreateChatScene);
