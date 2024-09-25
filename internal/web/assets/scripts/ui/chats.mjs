import { config } from "/assets/scripts/config.mjs";

import { PaginatedList } from "./custom-elements/paginated-list.mjs";

export class Chat {
  /**
   * @param {number} id
   * @param {string} name
   */
  constructor(id, name) {
    this.id = id;
    this.name = name;
  }
}

class ChatsRenderer {
  /** @param {Chat} chat  */
  createElement(chat) {
    const a = document.createElement("a");
    a.href = "/chats/" + chat.id;
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

export function initChats() {
  const list = /** @type {PaginatedList<Chat>} */ (
    assertInstance(document.getElementById("chats"), PaginatedList)
  );
  list.setIterator(new ChatsIterator());
  list.setRenderer(new ChatsRenderer());
  list.init();
}
