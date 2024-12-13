import { LocationControll } from "/assets/scripts/lib/location-control.mjs";

import { html } from "../lib/html.mjs";

const messageContentForm = html`
  <style>
    #message-form {
      display: flex;
      justify-content: center;
      align-items: flex-end;
      gap: 8px;
    }

    #message-content-input {
      max-height: 50vh;
      max-width: min(calc(100% - 16px), 90ch);
      width: 100%;
      background: var(--bg-2);
      color: var(--text-0);
      padding: 8px 16px;
      margin: 0px;
      border-radius: 20px;
      resize: none;
      outline: none;
      border: none;
    }

    #message-content-input:invalid + #submit-message {
      background: var(--bg-2);
      color: rgb(from var(--text-0) r g b / 0.25);
      cursor: auto;
    }

    #submit-message {
      --_size: 32px;
      transition: all var(--color-transition-duration);
      flex-shrink: 0;
      background: var(--primary);
      font-size: calc(var(--_size) * 0.66);
      color: var(--text);
      outline: none;
      border: none;
      border-radius: var(--_size);
      width: var(--_size);
      height: var(--_size);
      cursor: pointer;
    }
  </style>

  <form id="message-form" is="hermes-message-content-form">
    <textarea
      id="message-content-input"
      is="hermes-textarea-autoresize"
      focus-on-input="true"
      max-rows="12"
      name="content"
      placeholder="content..."
      autofocus
      required
    ></textarea>
    <button id="submit-message" type="submit">↑</button>
  </form>
`;

const scenes = {
  templates: html`
    <style>
      main {
        height: 100%;
        display: grid;
        place-items: center;
        text-align: center;
      }
    </style>

    <main>
      <div>
        <h1>under construction</h1>
        <a is="hermes-link" href="/" id="new-chat" class="chat-link">
          back to chats
        </a>
      </div>
    </main>
  `,

  home: html`
    <style>
      main {
        height: 100%;
        display: grid;
        place-items: center;
        padding: 16px;
      }
      form {
        max-width: 60ch;
        width: 100%;
      }
    </style>

    <main>${messageContentForm}</main>
  `,

  chats: html`
    <style>
      #chats-list {
        padding: 5px 0;
        max-height: 100vh;
        overflow: auto;
        border-right: 1px solid var(--bg-1);
        grid-auto-rows: max-content;
      }

      .chat-link {
        transition: all 50ms;
        color: var(--text-0);
        border: 1px solid transparent;
        text-decoration: none;
        padding: 4px 8px;
        margin: 1px 6px;
        display: block;
        border-radius: 8px;
        white-space: nowrap;
        overflow: hidden;
        text-overflow: ellipsis;
        animation: fade-in 250ms forwards;
      }

      .chat-link:hover {
        scale: 1.05;
        border-color: var(--primary);
      }

      .input-form-wrapper {
        padding: 4px 16px 16px;
      }
    </style>

    <style>
      main {
        display: grid;
        grid-template-columns: 250px 1fr;
      }

      #templates {
        position: fixed;
        right: 0;
        top: 0;
      }

      #chat-content {
        height: 100vh;
        max-height: 100vh;
        overflow: auto;
        display: grid;
        grid-template-rows: 1fr auto;
      }

      #messages-list-wrapper {
        max-height: 100%;
        overflow: auto;
        display: flex;
        flex-direction: column-reverse;
      }

      #messages-list {
        margin: 0 auto;
        width: 100%;
        max-width: 120ch;
      }

      .message {
        border: 1px solid var(--bg-2);
        padding: 4px 8px;
        margin: 12px;
        width: fit-content;
        max-width: 80%;
        border-radius: 10px;
        background: var(--bg-1);

        pre {
          margin: 0;
          text-wrap: wrap;
        }
      }

      .role-assistant {
        color: var(--text-0);
        border-bottom-left-radius: 0;
      }

      .role-user {
        border-bottom-right-radius: 0;
        margin-left: auto;
        color: var(--light-bg-0);
      }
    </style>

    <main>
      <hermes-chats id="chats-list">
        <a is="hermes-link" href="/" id="new-chat" class="chat-link">
          New chat
        </a>
      </hermes-chats>
      <div id="chat-content">
        <div id="messages-list-wrapper">
          <hermes-messages></hermes-messages>
          <a tabindex="1" id="templates" is="hermes-link" href="/templates">
            go to templates
          </a>
        </div>

        <div class="input-form-wrapper">${messageContentForm}</div>
      </div>
    </main>
  `,
};

class Scene extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanup = [];

  #activeSceneName = "__unset__";

  constructor() {
    super();
  }

  connectedCallback() {
    this.#cleanup.push(LocationControll.attach(this));
    this.notify(LocationControll.pathname);
  }

  disconnectedCallback() {
    this.#cleanup.forEach((cb) => cb());
  }

  /**
   * @param {string} pathname
   */
  notify(pathname) {
    const { name, html } = this.#scenePicker(pathname);
    if (this.#activeSceneName === name) {
      return;
    }
    this.#activeSceneName = name;
    this.innerHTML = html;
  }

  /**
   * @param {string} pathname
   * @returns {{name: keyof typeof scenes, html: string}} html
   */
  #scenePicker(pathname) {
    if (pathname.startsWith("/templates")) {
      return {
        name: "templates",
        html: scenes.templates,
      };
    }
    if (LocationControll.chatId) {
      return {
        name: "chats",
        html: scenes.chats,
      };
    }
    return {
      name: "home",
      html: scenes.home,
    };
  }
}

customElements.define("hermes-scene", Scene);
