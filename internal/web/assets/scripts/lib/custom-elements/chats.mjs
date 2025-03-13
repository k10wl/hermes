import { config } from "/assets/scripts/config.mjs";
import { AssertInstance, AssertString } from "/assets/scripts/lib/assert.mjs";
import { currentUrl } from "/assets/scripts/lib/current-url.mjs";
import { escapeMarkup } from "/assets/scripts/lib/escape-markup.mjs";
import { ServerEvents } from "/assets/scripts/lib/events/server-events.mjs";
import { Bind, html } from "/assets/scripts/lib/libdim.mjs";
import { LocationControll } from "/assets/scripts/lib/location-control.mjs";
import { Chat } from "/assets/scripts/models.mjs";

import { ShortcutManager } from "../shortcut-manager.mjs";
import { PaginatedList } from "./paginated-list.mjs";

export class Chats extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanup = [];

  constructor() {
    super();
  }

  connectedCallback() {
    const list = new Bind(
      (el) =>
        /** @type {PaginatedList<Chat>} */ (
          AssertInstance.once(el, PaginatedList)
        ),
    );
    this.attachShadow({ mode: "open" }).append(html`
      <style>
        * {
          box-sizing: border-box;
        }
        a {
          display: block;
          margin: 0.1rem 0;
          text-decoration: none;
        }
        a:first-child {
          margin-bottom: 0;
        }
      </style>

      <h-paginated-list bind="${list}"></h-paginated-list>
    `);
    this.findNextChat = this.navigateInDir.bind(this);

    const iterator = new ChatsIterator();
    const rendrer = new ChatsRenderer();

    this.#cleanup.push(
      ServerEvents.on("chat-created", (data) => {
        list.current.prepandNodes(rendrer.createElement(data.payload.chat));
      }),
      ShortcutManager.keydown("<M-ArrowUp>", this.navigateInDir("prev")),
      ShortcutManager.keydown("<M-ArrowDown>", this.navigateInDir("next")),
    );

    list.current.setIterator(iterator);
    list.current.setRenderer(rendrer);
    list.current.init();
  }

  // XXX maybe into separate file?
  /**
   * @param {KeyboardEvent} event
   */
  navigatheHome(event) {
    /** @type {HTMLElement | undefined} */
    let el = undefined;
    try {
      el = AssertInstance.once(event.target, HTMLElement);
      el.blur();
    } catch {
      // whatever
    }
    if ((el === document.body || !el) && LocationControll.chatId) {
      LocationControll.navigate("/");
    }
  }

  /**
   * @param {"prev" | "next"} dir
   * @returns {(event: KeyboardEvent) => void}
   */
  navigateInDir(dir) {
    return (event) => {
      event.stopPropagation();
      event.preventDefault();
      const target = this.getSibling()[dir];
      if (target === null) {
        return;
      }
      target.scrollIntoView({ block: "nearest" });
      LocationControll.navigate(AssertString.check(target.href));
    };
  }

  disconnectedCallback() {
    for (const cleanup of this.#cleanup) {
      cleanup();
    }
  }

  /** @returns {{prev: null | HTMLAnchorElement, next: null | HTMLAnchorElement}} */
  getSibling() {
    // XXX not so elegant, but works
    const all = this.querySelectorAll("a");
    const path = currentUrl(LocationControll.pathname);
    /** @type {ReturnType<Chats['getSibling']>} */
    const res = { prev: null, next: null };
    for (let i = 0; i < all.length; i++) {
      const el = AssertInstance.once(all[i], HTMLAnchorElement);
      if (el.href !== path) {
        continue;
      }
      try {
        res.prev = AssertInstance.once(all[i - 1], HTMLAnchorElement);
      } catch {
        res.prev = null;
      }
      try {
        res.next = AssertInstance.once(all[i + 1], HTMLAnchorElement);
      } catch {
        res.next = null;
      }
      break;
    }
    return res;
  }
}

class ChatsRenderer {
  /**
   * @param {Chat} chat
   * @returns {DocumentFragment}
   */
  createElement(chat) {
    const name =
      escapeMarkup(chat.name.replaceAll(/(\n|\s)+/gi, " ")) || "**empty**";
    return html`
      <a href="/chats/${chat.id}" id="chat-${chat.id}">
        <h-card data-interactive>${name}</h-card>
      </a>
    `;
  }
}

class ChatsIterator {
  hasMore = true;
  #limit = config.chats.paginationLimit;
  #startBeforeID = -1;

  async #fetchChats() {
    const res = await fetch(
      `/api/v1/chats?limit=${this.#limit}&start-before-id=${this.#startBeforeID}`,
    );
    const data = await res.json();
    /** @type {Chat[]} */
    const chats = [];
    for (const chat of data) {
      chats.push(new Chat(chat.id, chat.name));
    }
    return chats;
  }

  /** @param {Chat[]} chats */
  #updateState(chats) {
    this.hasMore = chats.length === this.#limit;
    this.#startBeforeID = chats.at(-1)?.id ?? 0;
  }

  async next() {
    if (!this.hasMore) {
      return [];
    }
    const chats = await this.#fetchChats();
    this.#updateState(chats);
    return chats;
  }
}

customElements.define("h-chats", Chats);
