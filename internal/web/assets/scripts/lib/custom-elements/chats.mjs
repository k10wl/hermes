import { config } from "/assets/scripts/config.mjs";
import { AssertInstance, AssertString } from "/assets/scripts/lib/assert.mjs";
import { currentUrl } from "/assets/scripts/lib/current-url.mjs";
import { ServerEvents } from "/assets/scripts/lib/events/server-events.mjs";
import { LocationControll } from "/assets/scripts/lib/location-control.mjs";
import { Chat } from "/assets/scripts/models.mjs";

import { ShortcutManager } from "../shortcut-manager.mjs";
import { PaginatedList } from "./paginated-list.mjs";

export class Chats extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanup = [];

  constructor() {
    super();
    this.innerHTML = `
<h-paginated-list>
    <!--<a is="hermes-link" href="/" class="chat-link">New chat</a>-->
</h-paginated-list>`;
    this.findNextChat = this.navigateInDir.bind(this);
  }

  connectedCallback() {
    const query = this.getElementsByTagName("h-paginated-list");
    const list = /** @type {PaginatedList<Chat>} */ (
      AssertInstance.once(query[0], PaginatedList)
    );

    const iterator = new ChatsIterator();
    const rendrer = new ChatsRenderer();

    this.#cleanup.push(
      ServerEvents.on("chat-created", (data) => {
        list.prepandNodes(rendrer.createElement(data.payload.chat));
      }),
      ShortcutManager.keydown("<M-ArrowUp>", this.navigateInDir("prev")),
      ShortcutManager.keydown("<M-ArrowDown>", this.navigateInDir("next")),
    );

    list.setIterator(iterator);
    list.setRenderer(rendrer);
    list.init();
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
  /** @param {Chat} chat  */
  createElement(chat) {
    const a = document.createElement("a", { is: "hermes-link" });
    const href = "/chats/" + chat.id;
    a.href = href;
    a.id = "chat-" + chat.id;
    a.classList.add("chat-link");
    a.innerText = chat.name.replaceAll(/(\n|\s)+/gi, " ");
    return a;
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
