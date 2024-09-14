import { config } from "/assets/scripts/config.mjs";
import { currentUrl } from "/assets/scripts/utils/current-url.mjs";

export class ServerEvents {
  /** @typedef {{type: string, payload: unknown}} ServerEvent */
  /** @typedef {{once?: boolean}} options */
  /** @typedef {(data: {type: string, payload: unknown}) => void} callback */
  /** @typedef {() => void} unsubscribe */
  /** @type WebSocket | null */
  static #connection = null;
  /** @type Map<string, {callback: callback, options: options}[]> */
  static #listeners = new Map();
  static #reconnectTimeout = 1000;
  /** @type (() => void)[] */
  static #onClose = [];
  /** @type (() => void)[] */
  static #onOpen = [];

  /** @param {string} addr */
  static __init(addr) {
    if (ServerEvents.#connection !== null) {
      throw new Error("do not reinitialize server events");
    }
    let clientId = sessionStorage.getItem("ws") ?? crypto.randomUUID();
    window.addEventListener("beforeunload", () => {
      sessionStorage.setItem("ws", clientId);
    });
    sessionStorage.removeItem("ws");

    const url = new URL(addr);
    url.searchParams.set("id", clientId);
    ServerEvents.#connection = new WebSocket(url.toString());
    ServerEvents.#addListeners(ServerEvents.#connection);
  }

  /**
   * @param {string} type
   * @param {(data: ServerEvent) => void} callback
   * @param {options} options
   * @returns {unsubscribe} unsubscribe
   */
  static on(type, callback, options = {}) {
    let events = ServerEvents.#listeners.get(type);
    if (!events) {
      events = [];
      ServerEvents.#listeners.set(type, events);
    }
    events.push({ callback, options });
    return () => ServerEvents.off(type, callback);
  }

  /**
   * @param {string} type
   * @param {(data: ServerEvent) => void} callback
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

  /** @param {() => void} callback  */
  static onOpen(callback) {
    ServerEvents.#onOpen.push(callback);
  }

  /** @param {() => void} callback  */
  static onClose(callback) {
    ServerEvents.#onClose.push(callback);
  }

  /** @param {WebSocket} webSocket */
  static #addListeners(webSocket) {
    webSocket.addEventListener(
      "open",
      () => {
        ServerEvents.#log("connected");
        for (let i = 0; i < ServerEvents.#onOpen.length; i++) {
          ServerEvents.#onOpen[i]?.();
        }
        webSocket.addEventListener("close", (closeEvent) => {
          ServerEvents.#log("connection closed", closeEvent);
          ServerEvents.#reconenct();
          for (let i = 0; i < ServerEvents.#onClose.length; i++) {
            ServerEvents.#onClose[i]?.();
          }
        });
        webSocket.addEventListener("err", (errorEvent) => {
          ServerEvents.#log("connection error", errorEvent);
          ServerEvents.#reconenct();
          for (let i = 0; i < ServerEvents.#onClose.length; i++) {
            ServerEvents.#onClose[i]?.();
          }
        });
        webSocket.addEventListener("message", (messageEvent) => {
          try {
            ServerEvents.#receiveEvent(messageEvent.data);
          } catch (error) {
            ServerEvents.#log("failed to hanlde message", error, messageEvent);
          }
        });
      },
      { once: true },
    );
  }

  /** @param {unknown} data  */
  static #receiveEvent(data) {
    if (typeof data !== "string") {
      return;
    }
    /** @type ServerEvent */
    const obj = JSON.parse(data);
    if (!("type" in obj)) {
      throw new Error("bad structure");
    }
    const listeners = ServerEvents.#listeners.get(obj.type);
    if (!listeners) {
      return;
    }
    for (let i = 0; i < listeners.length; i++) {
      const handler = listeners[i];
      if (!handler) {
        continue;
      }
      handler.callback(obj);
      if (handler.options.once) {
        ServerEvents.off(obj.type, handler.callback);
      }
    }
  }

  /** @param {any[]} data  */
  static #log(...data) {
    console.log(`[ServerEvents]`, ...data);
  }

  /** @throws if connection is null */
  static async #reconenct() {
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
        url.searchParams.set("reconnect", "1");
        ServerEvents.#connection = new WebSocket(url.toString());
        ServerEvents.#addListeners(ServerEvents.#connection);
        return;
      }
      ServerEvents.#reconenct();
    } catch {
      await new Promise((resolve) =>
        setTimeout(resolve, ServerEvents.#reconnectTimeout),
      );
      ServerEvents.#reconenct();
    }
  }
}

ServerEvents.__init(currentUrl(config.server.pathnames.webSocket));