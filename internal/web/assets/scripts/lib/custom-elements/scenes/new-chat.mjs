import { AssertInstance } from "../../assert.mjs";
import { html } from "../../html.mjs";
import { controlPalanelVisibility } from "../control-panel.mjs";

const HERMES_SHORTCUTS_TAG_NAME = "hermes-new-chat-shortcuts";
customElements.define(
  HERMES_SHORTCUTS_TAG_NAME,
  class Shortcuts extends HTMLElement {
    /** @type {{name: string, key: string, onclick?: () => void}[]} */
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
      this.shadow.innerHTML = html`
        <style>
          ul {
            display: grid;
            gap: 0.25rem;
            color: rgb(from var(--text-0) r g b / 0.33);
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
            background: transparent;
            outline: none;
            border: none;
            cursor: pointer;
            code:hover {
              --color: rgb(from var(--primary) r g b / 1);
            }
          }

          code {
            --color: rgb(from var(--text-0) r g b / 0.33);
            --transition: var(--color-transition-duration);
            transition:
              color var(--transition),
              border-color var(--transition);
            margin-right: auto;
            border-radius: 0.2rem;
            padding: 0.1rem 0.2rem;
            border: 1px solid var(--color);
            color: var(--color);
          }
        </style>

        <ul>
          ${Shortcuts.#list
            .map(
              (shortcut) => html`
                <li>
                  <span>${shortcut.name}</span>
                  <button><code>${shortcut.key}</code></button>
                </li>
              `,
            )
            .join("")}
        </ul>
      `;
    }

    connectedCallback() {
      const buttons = this.shadow.querySelectorAll("button");
      buttons.forEach((button, index) =>
        button.addEventListener("click", () =>
          Shortcuts.#list[index]?.onclick(),
        ),
      );
    }
  },
);

export class CreateChatScene extends HTMLElement {
  constructor() {
    super();
    this.attachShadow({ mode: "closed" }).innerHTML = html`
      <style>
        main {
          height: 100%;
          display: flex;
          justify-content: center;
          align-items: center;
          flex-direction: column;
          padding: 1rem;
          gap: 1rem;
        }

        hermes-message-form {
          width: 100%;
          max-width: var(--container);
        }
      </style>

      <main>
        <hermes-message-form></hermes-message-form>
        <${HERMES_SHORTCUTS_TAG_NAME}></${HERMES_SHORTCUTS_TAG_NAME}>
      </main>
    `;
  }
}

customElements.define("hermes-new-chat-scene", CreateChatScene);
