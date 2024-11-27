import { RequestReadChatEvent } from "/assets/scripts/events/client-events-list.mjs";
import { ServerEvents } from "/assets/scripts/events/server-events.mjs";
import { LocationControll } from "/assets/scripts/lib/navigation/location.mjs";
import { SoundManager } from "/assets/scripts/lib/sound-manager.mjs";

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
    const messageCreatedObserver = new MessageCreatedObserver(container);
    const audioNotificationObserver = new AudioNotificaitonsObserver();
    this.#cleanupOnDisconnect.push(
      LocationControll.attach(routeObserver),
      ServerEvents.on("message-created", (data) => {
        messageCreatedObserver.notify(data);
        audioNotificationObserver.notify(data);
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
  #chatRegex = /\/chats\/(?<chatId>\d+)$/;

  /** @param {HTMLElement} container  */
  constructor(container) {
    this.container = container;
  }

  /** @param {string} pathname  */
  notify(pathname) {
    const match = this.#chatRegex.exec(pathname);
    if (match?.groups?.chatId) {
      ServerEvents.send(new RequestReadChatEvent(+match.groups.chatId));
      return;
    }
    this.container.innerHTML = "";
  }
}

class MessagesViewObserver {
  container;
  /** @type {null | number} */
  chatId;
  /** @param {HTMLElement} container  */
  constructor(container) {
    this.container = container;
    this.chatId = null;
  }

  /** @param {import( "/assets/scripts/events/server-events-list.mjs").ReadChatEvent} readChatEvent  */
  notify(readChatEvent) {
    if (this.chatId !== readChatEvent.payload.messages.at(0)?.chat_id) {
      this.container.innerHTML = "";
    }
    this.container.append(
      ...readChatEvent.payload.messages.map((message) =>
        MessageCreator.createElement(message),
      ),
    );
  }
}

class AudioNotificaitonsObserver {
  /**
   * @param {import("/assets/scripts/events/server-events-list.mjs").MessageCreatedEvent } event
   */
  notify(event) {
    if (event.payload.message.role === "user") {
      return;
    }
    if (event.payload.chat_id === LocationControll.chatId) {
      SoundManager.play("message-in-local");
      return;
    }
    SoundManager.play("message-in-global");
  }
}

class MessageCreatedObserver {
  #container;
  /** @param {HTMLElement} container  */
  constructor(container) {
    this.#container = container;
  }
  /** @param {import("/assets/scripts/events/server-events-list.mjs").MessageCreatedEvent } event  */
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
    div.classList.add("message", `role-${message.role}`);
    div.id = `message-${message.id}`;
    const pre = document.createElement("pre");
    pre.innerText = message.content;
    div.appendChild(pre);
    return div;
  }
}
