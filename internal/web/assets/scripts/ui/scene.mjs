import { LocationControll } from "/assets/scripts/lib/location-control.mjs";

import { html } from "../lib/html.mjs";

const scenes = {
  "/": html`<hermes-new-chat-scene></hermes-new-chat-scene>`,

  "/chats": html`<hermes-chats-list-scene></hermes-chats-list-scene>`,

  "/chats/{id}": html`<hermes-existing-chat-scene></hermes-existing-chat-scene>`,

  "/templates": html`
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
          back to new chat
        </a>
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
    this.shadow = this.attachShadow({ mode: "closed" });
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
    this.shadow.innerHTML = html;
  }

  /**
   * @param {string} pathname
   * @returns {{name: keyof typeof scenes, html: string}} html
   */
  #scenePicker(pathname) {
    if (pathname.startsWith("/templates")) {
      return {
        name: "/templates",
        html: scenes["/templates"],
      };
    }
    const isChats = pathname.startsWith("/chats");
    if (isChats && LocationControll.chatId) {
      return {
        name: "/chats/{id}",
        html: scenes["/chats/{id}"],
      };
    }
    if (isChats) {
      return {
        name: "/chats",
        html: scenes["/chats"],
      };
    }
    return {
      name: "/",
      html: scenes["/"],
    };
  }
}

customElements.define("hermes-scene", Scene);
