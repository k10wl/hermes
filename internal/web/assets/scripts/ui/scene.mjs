import { LocationControll } from "/assets/scripts/lib/location-control.mjs";

import { html } from "../lib/libdim.mjs";

// TODO sooo without creating new fragments everytime app stucks
const scenes = {
  "/": () => html`<hermes-new-chat-scene></hermes-new-chat-scene>`,
  "/chats": () => html`<hermes-chats-list-scene></hermes-chats-list-scene>`,
  "/chats/{id}": () =>
    html`<hermes-existing-chat-scene></hermes-existing-chat-scene>`,
  "/templates": () =>
    html`<hermes-templates-list-scene></hermes-templates-list-scene>`,
  "/templates/new": () =>
    html`<hermes-view-template-scene></hermes-view-template-scene>`,
  "/templates/{id}": () =>
    html`<hermes-view-template-scene></hermes-view-template-scene>`,
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
    const { name, fragment } = this.#scenePicker(pathname);
    if (this.#activeSceneName === name) {
      return;
    }
    this.#activeSceneName = name;
    this.shadow.replaceChildren();
    this.shadow.append(fragment);
  }

  /**
   * @param {string} pathname
   * @returns {{name: keyof typeof scenes, fragment: DocumentFragment}}
   */
  #scenePicker(pathname) {
    if (pathname.startsWith("/templates/new")) {
      return {
        name: "/templates/new",
        fragment: scenes["/templates/new"](),
      };
    }
    if (pathname.startsWith("/templates")) {
      if (/\d+$/.test(pathname)) {
        return {
          name: "/templates/{id}",
          fragment: scenes["/templates/{id}"](),
        };
      }
      return {
        name: "/templates",
        fragment: scenes["/templates"](),
      };
    }
    const isChats = pathname.startsWith("/chats");
    if (isChats && LocationControll.chatId) {
      return {
        name: "/chats/{id}",
        fragment: scenes["/chats/{id}"](),
      };
    }
    if (isChats) {
      return {
        name: "/chats",
        fragment: scenes["/chats"](),
      };
    }
    return {
      name: "/",
      fragment: scenes["/"](),
    };
  }
}

customElements.define("hermes-scene", Scene);
