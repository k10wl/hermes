import { config } from "/assets/scripts/config.mjs";
import { assertInstance } from "/assets/scripts/utils/assert-instance.mjs";
import { currentUrl } from "/assets/scripts/utils/current-url.mjs";
import { ValidateString } from "/assets/scripts/utils/validate.mjs";

import {
  ChatCreatedEvent,
  ConnectionStatusChangeEvent,
  MessageCreatedEvent,
  ReadChatEvent,
  ServerErrorEvent,
  ServerEvent,
} from "./server-events-list.mjs";

/**
 * @typedef {Object} IsolatedServiceEvents
 * @property {ConnectionStatusChangeEvent} connection-status-change
 */

/**
 * @typedef {Object} ClientEmittedEvents
 * @property {import("./client-events-list.mjs").RequestReadChatEvent} request-read-chat
 */

/**
 * @typedef {Object} ServerEmittedEvents
 * @property {ChatCreatedEvent} chat-created
 * @property {ServerEvent} reload
 * @property {ReadChatEvent} read-chat
 * @property {ServerErrorEvent} server-error
 * @property {MessageCreatedEvent} message-created
 */

/** @typedef { ServerEmittedEvents & IsolatedServiceEvents } RegisteredEvents */

export class ServerEvents {
  /**
   * @template {keyof RegisteredEvents} T
   * @typedef {(data: RegisteredEvents[T]) => void} callback
   */
  /** @typedef {{once?: boolean}} options */
  /** @typedef {() => void} unsubscribe */
  /** @type WebSocket | null */
  static #connection = null;
  /** @type Map<string, {callback: callback<any>, options: options}[]> */
  static #listeners = new Map();
  static #reconnectTimeout = 1000;
  static #allowReconnect = true;

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

  /**
   * @template {keyof RegisteredEvents} T
   * @param {T} type
   * @param {callback<T>} callback
   * @param {object} [options]
   * @returns {unsubscribe} unsubscribe
   */
  static on(type, callback, options = {}) {
    let events = ServerEvents.#listeners.get(type);
    ServerEvents.#listeners.get("reload");
    if (!events) {
      events = [];
      ServerEvents.#listeners.set(type, events);
    }
    events.push({ callback, options });
    return () => ServerEvents.off(type, callback);
  }

  /**
   * @template {keyof RegisteredEvents} T
   * @param {T} type
   * @param {callback<T>} callback
   */
  static off(type, callback) {
    let events = ServerEvents.#listeners.get(type);
    if (!events) {
      return;
    }
    const callbackIndex = events.findIndex(
      (handler) => handler.callback === callback,
    );
    if (callbackIndex === -1) {
      return;
    }
    events.splice(callbackIndex, 1);
  }

  /**
   * @template {keyof ClientEmittedEvents} T
   * @param {ClientEmittedEvents[T]} event
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
    ServerEvents.#notifySubscribers(new ConnectionStatusChangeEvent(true));
  }

  /** @param {MessageEvent<unknown>} event */
  static #onMessage(event) {
    try {
      ServerEvents.#notifySubscribers(
        EmittedServerEventFactory.parse(event.data),
      );
    } catch (error) {
      ServerEvents.#log("failed to handle message", error, event);
    }
  }

  /** @param {Event} event */
  static #onClose(event) {
    ServerEvents.#log("connection closed", event);
    ServerEvents.#notifySubscribers(new ConnectionStatusChangeEvent(false));
    ServerEvents.#reconnect();
  }

  /** @param {Event} event  */
  static #onError(event) {
    ServerEvents.#log("connection error", event);
    ServerEvents.#notifySubscribers(new ConnectionStatusChangeEvent(false));
    ServerEvents.#reconnect();
  }

  /** @param {RegisteredEvents[keyof RegisteredEvents]} event */
  static #notifySubscribers(event) {
    const listeners = ServerEvents.#listeners.get(event.type);
    if (!listeners) {
      return;
    }
    for (let i = 0; i < listeners.length; i++) {
      const handler = listeners[i];
      if (!handler) {
        continue;
      }
      handler.callback(event);
      if (handler.options.once) {
        ServerEvents.off(event.type, handler.callback);
      }
    }
  }

  /** @param {any[]} data  */
  static #log(...data) {
    console.log(`[ServerEvents]`, ...data);
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
      ServerEvents.#log("connection lost, reconnecting...");
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

// I HATE ADDING EVERY MESSAGE IN HERE HOLYFUCK CAN THIS BE AUTOMATED???
class EmittedServerEventFactory {
  static #typeRegex = /"type":(\s?)+"(?<type>.*?)"/;
  /**
   * @param {unknown} data
   * @returns {RegisteredEvents[keyof RegisteredEvents]}
   */
  static parse(data) {
    const res = EmittedServerEventFactory.#typeRegex.exec(
      ValidateString.parse(data),
    );
    switch (ValidateString.parse(res?.groups?.type)) {
      case "reload":
        return ServerEvent.parse(data);
      case "server-error":
        return ServerErrorEvent.parse(data);
      case "chat-created":
        return ChatCreatedEvent.parse(data);
      case "message-created":
        return MessageCreatedEvent.parse(data);
      case "read-chat":
        return ReadChatEvent.parse(data);
      default:
        throw new Error(`unhandled server event - ${data}`);
    }
  }
}

ServerEvents.__init(currentUrl(config.server.pathnames.webSocket));
Reflect.set(window, "ServerEvents", ServerEvents);
