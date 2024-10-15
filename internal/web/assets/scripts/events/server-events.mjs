import { config } from "/assets/scripts/config.mjs";
import { currentUrl } from "/assets/scripts/utils/current-url.mjs";
import { ValidateString } from "/assets/scripts/utils/validate.mjs";

import {
  ChatCreatedEvent,
  ConnectionStatusChangeEvent,
  ServerEvent,
} from "./server-events-list.mjs";

/**
 * @typedef {Object} expectedEvents
 * @property {ChatCreatedEvent} chat-created
 * @property {ConnectionStatusChangeEvent} connection-status-change
 * @property {ServerEvent} reload
 */

export class ServerEvents {
  /**
   * @template {keyof expectedEvents} T
   * @typedef {(data: expectedEvents[T]) => void} callback
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
   * @template {keyof expectedEvents} T
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
   * @template {keyof expectedEvents} T
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

  /** @param {WebSocket} webSocket */
  static #addListeners(webSocket) {
    webSocket.addEventListener(
      "open",
      () => {
        ServerEvents.#log("connected");
        ServerEvents.#notifySubscribers(new ConnectionStatusChangeEvent(true));
        webSocket.addEventListener("close", (closeEvent) => {
          ServerEvents.#log("connection closed", closeEvent);
          ServerEvents.#notifySubscribers(
            new ConnectionStatusChangeEvent(false),
          );
          ServerEvents.#reconnect();
        });
        webSocket.addEventListener("message", (messageEvent) => {
          try {
            ServerEvents.#notifySubscribers(
              ServerEvents.#parseEvent(messageEvent.data),
            );
          } catch (error) {
            ServerEvents.#log("failed to hanlde message", error, messageEvent);
          }
        });
      },
      { once: true },
    );
    webSocket.addEventListener(
      "error",
      (errorEvent) => {
        ServerEvents.#log("connection error", errorEvent);
        ServerEvents.#notifySubscribers(new ConnectionStatusChangeEvent(false));
        ServerEvents.#reconnect();
      },
      { once: true },
    );
  }

  /** @param {expectedEvents[keyof expectedEvents]} event */
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

  static #typeRegex = /"type":(\s?)+"(?<type>.*?)"/;
  /**
   * @param {string} data
   * @returns {expectedEvents[keyof expectedEvents]}
   */
  static #parseEvent(data) {
    const res = ServerEvents.#typeRegex.exec(ValidateString.parse(data));
    switch (ValidateString.parse(res?.groups?.type)) {
      case "reload":
        return ServerEvent.parse(data);
      case "chat-created":
        return ChatCreatedEvent.parse(data);
      default:
        throw new Error(`unhandled server event - ${data}`);
    }
  }
}

ServerEvents.__init(currentUrl(config.server.pathnames.webSocket));
