import { config } from "/assets/scripts/config.mjs";
import { AssertInstance, AssertString } from "/assets/scripts/lib/assert.mjs";
import { backoff, exponent } from "/assets/scripts/lib/backoff.mjs";
import { CallbackTracker } from "/assets/scripts/lib/callback-tracker.mjs";
import { Queue } from "/assets/scripts/lib/queue.mjs";
import { sleep } from "/assets/scripts/lib/sleep.mjs";

import * as clientEventsList from "./client-events-list.mjs";
import * as serverEventsList from "./server-events-list.mjs";

const _isolatedServiceEvents = {
  [serverEventsList.ConnectionStatusChangeEvent.canonicalType]:
    serverEventsList.ConnectionStatusChangeEvent,
};

const _serverEvents = {
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
  [serverEventsList.ReadTemplatesEvent.canonicalType]:
    serverEventsList.ReadTemplatesEvent,
  [serverEventsList.ReadTemplateEvent.canonicalType]:
    serverEventsList.ReadTemplateEvent,
  [serverEventsList.TemplateChangedEvent.canonicalType]:
    serverEventsList.TemplateChangedEvent,
  [serverEventsList.TemplateCreatedEvent.canonicalType]:
    serverEventsList.TemplateCreatedEvent,
  [serverEventsList.TemplateDeletedEvent.canonicalType]:
    serverEventsList.TemplateDeletedEvent,
};

const _clientEvents = {
  [clientEventsList.RequestReadChatEvent.canonicalType]:
    clientEventsList.RequestReadChatEvent,
  [clientEventsList.CreateCompletionMessageEvent.canonicalType]:
    clientEventsList.CreateCompletionMessageEvent,
  [clientEventsList.RequestReadTemplatesEvent.canonicalType]:
    clientEventsList.RequestReadTemplatesEvent,
  [clientEventsList.RequestReadTemplateEvent.canonicalType]:
    clientEventsList.RequestReadTemplateEvent,
  [clientEventsList.DeleteTemplateEvent.canonicalType]:
    clientEventsList.DeleteTemplateEvent,
};

const _registeredEvents = {
  ..._serverEvents,
  ..._isolatedServiceEvents,
};

/** just to be sure that all described events are known and handled */
(function () {
  /** @type {Record<string, boolean>} */
  const usedEvents = {};
  Object.values(_registeredEvents).forEach((event) => {
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
  static #callbackTracker = new CallbackTracker(_registeredEvents);
  /** @type {Queue<InstanceType<typeof _clientEvents[keyof typeof _clientEvents]>>} */
  static #queue = new Queue();

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

  /** @type {CallbackTracker<typeof _registeredEvents>["on"]} */
  static on(...args) {
    return ServerEvents.#callbackTracker.on(...args);
  }

  /** @type {CallbackTracker<typeof _registeredEvents>["off"]} */
  static off(type, callback) {
    return ServerEvents.#callbackTracker.off(type, callback);
  }

  /**
   * @param {InstanceType<typeof _clientEvents[keyof typeof _clientEvents]>} event
   */
  static send(event) {
    const socket = AssertInstance.once(ServerEvents.#connection, WebSocket);
    if (socket.readyState !== WebSocket.OPEN) {
      ServerEvents.#queue.enqueue(event);
      ServerEvents.#flush();
      return;
    }
    try {
      ServerEvents.#unsafeSend(event);
    } catch (error) {
      ServerEvents.#error("error upon sending", event, error);
    }
  }

  static #flushing = false;
  static async #flush() {
    if (ServerEvents.#flushing) {
      return;
    }
    ServerEvents.#warn("entered recovering flush");
    ServerEvents.#flushing = true;
    const time = backoff(10, exponent);
    let event;
    while ((event = ServerEvents.#queue.peek())) {
      try {
        if (!ServerEvents.connected) {
          const promise = Promise.withResolvers();
          const off = ServerEvents.on("connection-status-change", (data) => {
            if (!data.payload.connected) {
              return;
            }
            promise.resolve(null);
            off();
          });
          await promise.promise;
        }
        ServerEvents.#unsafeSend(event);
        event = ServerEvents.#queue.dequeue();
      } catch (error) {
        ServerEvents.#error("caught in flush loop", error);
        await sleep(time());
      }
    }
    ServerEvents.#flushing = false;
    ServerEvents.#log("recovered in flush");
  }

  /**
   * @param {InstanceType<typeof _clientEvents[keyof typeof _clientEvents]>} event
   */
  static #unsafeSend(event) {
    const socket = AssertInstance.once(ServerEvents.#connection, WebSocket);
    socket.send(JSON.stringify(event));
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

  /** @param {InstanceType<_registeredEvents[keyof typeof _registeredEvents]>} event */
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
      const res = await fetch(config.server.pathnames.healthCheck);
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
   * @returns {InstanceType<typeof _registeredEvents[keyof typeof _registeredEvents]>}
   */
  static parse(data) {
    const res = AssertString.check(
      EmittedServerEventFactory.#typeRegex.exec(AssertString.check(data))
        ?.groups?.type,
    );
    if (!(res in _registeredEvents)) {
      throw new Error("receivd unhandled event");
    }
    // @ts-expect-error literally has check for "in" three lines above
    return _registeredEvents[res].parse(data);
  }
}

ServerEvents.__init(config.server.pathnames.webSocket);
Reflect.set(window, "ServerEvents", ServerEvents);
