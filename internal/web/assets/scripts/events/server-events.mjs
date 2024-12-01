import { config } from "/assets/scripts/config.mjs";
import { assertInstance } from "/assets/scripts/utils/assert-instance.mjs";
import { currentUrl } from "/assets/scripts/utils/current-url.mjs";
import { ValidateString } from "/assets/scripts/utils/validate.mjs";

import { CallbackTracker } from "./callback-tracker.mjs";
import * as clientEventsList from "./client-events-list.mjs";
import * as serverEventsList from "./server-events-list.mjs";

const __isolatedServiceEvents = {
  [serverEventsList.ConnectionStatusChangeEvent.canonicalType]:
    serverEventsList.ConnectionStatusChangeEvent,
};

const serverEvents = {
  [serverEventsList.ChatCreatedEvent.canonicalType]:
    serverEventsList.ChatCreatedEvent,
  [serverEventsList.ServerEvent.canonicalType]: serverEventsList.ServerEvent,
  [serverEventsList.ReadChatEvent.canonicalType]:
    serverEventsList.ReadChatEvent,
  [serverEventsList.ServerErrorEvent.canonicalType]:
    serverEventsList.ServerErrorEvent,
  [serverEventsList.MessageCreatedEvent.canonicalType]:
    serverEventsList.MessageCreatedEvent,
  [serverEventsList.ReloadEvent.canonicalType]: serverEventsList.ReloadEvent,
};

const _clientEvents = {
  [clientEventsList.RequestReadChatEvent.canonicalType]:
    clientEventsList.RequestReadChatEvent,
  [clientEventsList.CreateCompletionMessageEvent.canonicalType]:
    clientEventsList.CreateCompletionMessageEvent,
};

const registeredEvents = {
  ...serverEvents,
  ...__isolatedServiceEvents,
};

/** just to be sure that all described events are known and handled */
(function () {
  /** @type {Record<string, boolean>} */
  const usedEvents = {};
  Object.values(registeredEvents).forEach((event) => {
    if (usedEvents[event.canonicalType]) {
      throw new Error(
        `redeclaration of registered event ${event.canonicalType}`,
      );
    }
    usedEvents[event.canonicalType] = true;
  });
  Object.values(serverEventsList).forEach((event) => {
    if (
      serverEventsList.ServerEvent !== Object.getPrototypeOf(event) &&
      serverEventsList.ServerEvent !== event
    ) {
      throw new Error(
        `exported non server event from server events list -  ${event}`,
      );
    }
    if (!usedEvents[event.canonicalType]) {
      throw new Error(`did not registered event ${event.canonicalType}`);
    }
  });
})();

export class ServerEvents {
  /** @type WebSocket | null */
  static #connection = null;
  static #reconnectTimeout = 1000;
  static #allowReconnect = true;
  static #callbackTracker = new CallbackTracker(registeredEvents);

  /** @param {string} addr */
  static __init(addr) {
    if (ServerEvents.#connection !== null) {
      throw new Error("do not reinitialize server events");
    }
    window.addEventListener("beforeunload", () => {
      ServerEvents.#allowReconnect = false;
    });
    const url = new URL(addr);
    ServerEvents.#connection = new WebSocket(url.toString());
    ServerEvents.#addListeners(ServerEvents.#connection);
  }

  /** @type {CallbackTracker<typeof registeredEvents>["on"]} */
  static on(type, callback) {
    return ServerEvents.#callbackTracker.on(type, callback);
  }

  /** @type {CallbackTracker<typeof registeredEvents>["off"]} */
  static off(type, callback) {
    return ServerEvents.#callbackTracker.off(type, callback);
  }

  /**
   * @param {InstanceType<typeof _clientEvents[keyof typeof _clientEvents]>} event
   */
  static send(event) {
    assertInstance(ServerEvents.#connection, WebSocket).send(
      JSON.stringify(event),
    );
  }

  /** @param {WebSocket} webSocket */
  static #addListeners(webSocket) {
    webSocket.addEventListener("open", ServerEvents.#onOpen, { once: true });
    webSocket.addEventListener("close", ServerEvents.#onClose, { once: true });
    webSocket.addEventListener("message", ServerEvents.#onMessage);
    webSocket.addEventListener("error", ServerEvents.#onError, { once: true });
  }

  static #onOpen() {
    ServerEvents.#log("connected");
    ServerEvents.#notifySubscribers(
      new serverEventsList.ConnectionStatusChangeEvent(true),
    );
  }

  /** @param {MessageEvent<unknown>} event */
  static #onMessage(event) {
    try {
      ServerEvents.#notifySubscribers(
        EmittedServerEventFactory.parse(event.data),
      );
    } catch (error) {
      ServerEvents.#error("failed to handle message", error, event);
    }
  }

  /** @param {Event} event */
  static #onClose(event) {
    ServerEvents.#warn("connection closed", event);
    ServerEvents.#notifySubscribers(
      new serverEventsList.ConnectionStatusChangeEvent(false),
    );
    ServerEvents.#reconnect();
  }

  /** @param {Event} event  */
  static #onError(event) {
    ServerEvents.#error("connection error", event);
    ServerEvents.#notifySubscribers(
      new serverEventsList.ConnectionStatusChangeEvent(false),
    );
    ServerEvents.#reconnect();
  }

  /** @param {InstanceType<registeredEvents[keyof typeof registeredEvents]>} event */
  static #notifySubscribers(event) {
    const callbacks = ServerEvents.#callbackTracker.getCallbacks(event.type);
    if (!callbacks) {
      return;
    }
    for (const callback of callbacks) {
      callback(/** @type {any} */ (event));
    }
  }

  /** @param {any[]} data  */
  static #log(...data) {
    console.log(`[ServerEvents]`, ...data);
  }

  /** @param {any[]} data  */
  static #warn(...data) {
    console.warn(`[ServerEvents]`, ...data);
  }

  /** @param {any[]} data  */
  static #error(...data) {
    console.error(`[ServerEvents]`, ...data);
  }

  // TODO reconnect only if tab is opened and active
  /** @throws if connection is null */
  static async #reconnect() {
    if (!ServerEvents.#allowReconnect) {
      return;
    }
    if (ServerEvents.#connection === null) {
      throw new Error("attempt to reconnect to non existing socket");
    }
    if (ServerEvents.#connection.readyState !== WebSocket.CLOSED) {
      return;
    }
    try {
      ServerEvents.#warn("connection lost, reconnecting...");
      const res = await fetch(currentUrl(config.server.pathnames.healthCheck));
      if (res.status == 200) {
        const url = new URL(ServerEvents.#connection.url);
        url.searchParams.set("reconnect", "true");
        ServerEvents.#connection = new WebSocket(url.toString());
        ServerEvents.#addListeners(ServerEvents.#connection);
        return;
      }
      ServerEvents.#reconnect();
    } catch {
      await new Promise((resolve) =>
        setTimeout(resolve, ServerEvents.#reconnectTimeout),
      );
      ServerEvents.#reconnect();
    }
  }

  static get connected() {
    return ServerEvents.#connection?.readyState === WebSocket.OPEN;
  }
}

class EmittedServerEventFactory {
  static #typeRegex = /"type":(\s?)+"(?<type>.*?)"/;
  /**
   * @param {unknown} data
   * @returns {InstanceType<typeof registeredEvents[keyof typeof registeredEvents]>}
   */
  static parse(data) {
    const res = ValidateString.parse(
      EmittedServerEventFactory.#typeRegex.exec(ValidateString.parse(data))
        ?.groups?.type,
    );
    if (!(res in registeredEvents)) {
      throw new Error("receivd unhandled event");
    }
    // @ts-expect-error literally has check for "in" three lines above
    return registeredEvents[res].parse(data);
  }
}

ServerEvents.__init(currentUrl(config.server.pathnames.webSocket));
Reflect.set(window, "ServerEvents", ServerEvents);
