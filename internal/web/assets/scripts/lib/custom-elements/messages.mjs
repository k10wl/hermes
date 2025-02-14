import { RequestReadChatEvent } from "/assets/scripts/lib/events/client-events-list.mjs";
import { ServerEvents } from "/assets/scripts/lib/events/server-events.mjs";
import { LocationControll } from "/assets/scripts/lib/location-control.mjs";

export class Messages extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanupOnDisconnect = [];

  constructor() {
    super();
  }

  connectedCallback() {
    const container = document.createElement("div");
    this.appendChild(container);
    const messagesViewObserver = new MessagesViewObserver(container);
    const routeObserver = new RouteObserver(container);
    routeObserver.notify();
    const messageCreatedObserver = new MessageCreatedObserver(container);
    this.#cleanupOnDisconnect.push(
      LocationControll.attach(routeObserver),
      ServerEvents.on("message-created", (data) => {
        messageCreatedObserver.notify(data);
      }),
      ServerEvents.on("read-chat", (data) => messagesViewObserver.notify(data)),
    );
  }

  disconnectedCallback() {
    for (const cleanup of this.#cleanupOnDisconnect) {
      cleanup();
    }
  }
}

class RouteObserver {
  container;

  /** @type {number | null} */
  #last = null;

  /** @param {HTMLElement} container  */
  constructor(container) {
    this.container = container;
  }

  notify() {
    const chatId = LocationControll.chatId;
    const prev = this.#last;
    this.#last = chatId;
    if (prev === chatId) {
      return;
    }
    if (chatId) {
      ServerEvents.send(new RequestReadChatEvent(chatId));
      return;
    }
    this.container.innerHTML = "";
  }
}

class MessagesViewObserver {
  container;
  /** @param {HTMLElement} container  */
  constructor(container) {
    this.container = container;
  }

  /** @param {import( "../events/server-events-list.mjs").ReadChatEvent} readChatEvent  */
  notify(readChatEvent) {
    this.container.innerHTML = "";
    this.container.append(
      ...readChatEvent.payload.messages.map((message) =>
        MessageCreator.createElement(message),
      ),
    );
  }
}

class MessageCreatedObserver {
  #container;
  /** @param {HTMLElement} container  */
  constructor(container) {
    this.#container = container;
  }
  /** @param {import("../events/server-events-list.mjs").MessageCreatedEvent } event  */
  notify(event) {
    if (event.payload.chat_id !== LocationControll.chatId) {
      return;
    }
    this.#container.append(MessageCreator.createElement(event.payload.message));
  }
}

class MessageCreator {
  /** @param {import("/assets/scripts/models.mjs").Message} message  */
  static createElement(message) {
    const div = document.createElement("div");
    div.classList.add("message", `role-${message.role}`, "selectable");
    div.id = `message-${message.id}`;
    const pre = document.createElement("pre");
    pre.innerText = message.content.trim();
    div.appendChild(pre);
    return div;
  }
}
