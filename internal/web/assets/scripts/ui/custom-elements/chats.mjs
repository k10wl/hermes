import { config } from "/assets/scripts/config.mjs";
import { ServerEvents } from "/assets/scripts/events/server-events.mjs";
import { LocationControll } from "/assets/scripts/lib/navigation/location.mjs";
import { Chat } from "/assets/scripts/models.mjs";
import { assertInstance } from "/assets/scripts/utils/assert-instance.mjs";
import { currentUrl } from "/assets/scripts/utils/current-url.mjs";
import { ValidateString } from "/assets/scripts/utils/validate.mjs";

import { PaginatedList } from "./paginated-list.mjs";

export class Chats extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanup = [];

  constructor() {
    super();
    this.innerHTML = `
<hermes-paginated-list>
    <a is="hermes-link" href="/" class="chat-link">New chat</a>
</hermes-paginated-list>`;
  }

  connectedCallback() {
    const query = this.getElementsByTagName("hermes-paginated-list");
    const list = /** @type {PaginatedList<Chat>} */ (
      assertInstance(query[0], PaginatedList)
    );

    const activeChatObserver = new ActiveChatObserver(list);
    const iterator = new ChatsIterator();
    const rendrer = new ChatsRenderer(activeChatObserver);

    window.addEventListener("keydown", (e) => {
      if (e.key === "Escape") {
        /** @type {HTMLElement | undefined} */
        let el = undefined;
        try {
          el = assertInstance(window.document.activeElement, HTMLElement);
        } catch {
          // whatever
        }
        if (el) {
          el.blur();
        } else {
          LocationControll.navigate("/");
        }
        return;
      }
      if (!e.altKey) {
        return;
      }
      const dir =
        e.key === "ArrowDown" ? "next" : e.key === "ArrowUp" ? "prev" : null;
      if (dir === null) {
        return;
      }
      e.stopPropagation();
      e.preventDefault();
      const target = this.getSibling()[dir];
      if (target === null) {
        return;
      }
      target.scrollIntoView({ block: "nearest" });
      LocationControll.navigate(ValidateString.parse(target.href));
    });

    this.#cleanup.push(
      LocationControll.attach(activeChatObserver),
      ServerEvents.on("chat-created", (data) => {
        list.prepandNodes(rendrer.createElement(data.payload.chat));
      }),
    );

    list.setIterator(iterator);
    list.setRenderer(rendrer);
    list.init();
  }

  disconenctedCallback() {
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
      const el = assertInstance(all[i], HTMLAnchorElement);
      if (el.href !== path) {
        continue;
      }
      try {
        res.prev = assertInstance(all[i - 1], HTMLAnchorElement);
      } catch {
        res.prev = null;
      }
      try {
        res.next = assertInstance(all[i + 1], HTMLAnchorElement);
      } catch {
        res.next = null;
      }
      break;
    }
    return res;
  }
}

class ChatsRenderer {
  #activeProvider;

  /**
   * @param {{
   *   activePathname: string,
   *   updateActive: (el: HTMLAnchorElement) => void,
   * }} provider
   */
  constructor(provider) {
    this.#activeProvider = provider;
  }

  /** @param {Chat} chat  */
  createElement(chat) {
    const a = document.createElement("a", { is: "hermes-link" });
    const href = "/chats/" + chat.id;
    a.href = href;
    a.id = "chat-" + chat.id;
    a.classList.add("chat-link");
    if (chat.id === LocationControll.chatId) {
      this.#activeProvider.updateActive(a);
    }
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

class ActiveChatObserver {
  /** @type {HTMLAnchorElement | null} */
  #active = null;
  #container;

  activePathname = "";

  /** @param {HTMLElement} container  */
  constructor(container) {
    this.#container = container;
  }

  /** @param {HTMLAnchorElement | null} element  */
  updateActive(element) {
    this.#removeActive(this.#active);
    this.#active = element;
    this.#setActive(this.#active);
  }

  /** @param {string} pathname */
  notify(pathname) {
    this.activePathname = pathname;
    const selected = this.#container.querySelector(`a[href="${pathname}"]`);
    if (!selected) {
      return;
    }
    this.updateActive(assertInstance(selected, HTMLAnchorElement));
  }

  /** @param {HTMLAnchorElement | null} element  */
  #setActive(element) {
    element?.classList.add("primary-bg");
  }

  /** @param {HTMLAnchorElement | null} element  */
  #removeActive(element) {
    element?.classList.remove("primary-bg");
  }
}
