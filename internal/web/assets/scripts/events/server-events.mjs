import { config } from "/assets/scripts/config.mjs";
import { currentUrl } from "/assets/scripts/utils/current-url.mjs";

export class ServerEvents {
  /** @typedef {{type: string, payload: unknown}} ServerEvent */
  /** @typedef {{once?: boolean}} options */
  /** @typedef {(data: {type: string, payload: unknown}) => void} callback */
  /** @typedef {() => void} unsubscribe */
  /** @type WebSocket | null */
  #connection = null;
  /** @type Map<string, {callback: callback, options: options}[]> */
  #listeners = new Map();
  #reconnectTimeout = 1000;

  /**
   * @param {string} addr
   */
  constructor(addr) {
    this.#connection = new WebSocket(addr);
    this.#attach(this.#connection);
  }

  /**
   * @param {string} type
   * @param {(data: ServerEvent) => void} callback
   * @param {options} options
   * @returns {unsubscribe} unsubscribe
   */
  on(type, callback, options = {}) {
    let events = this.#listeners.get(type);
    if (!events) {
      events = [];
      this.#listeners.set(type, events);
    }
    events.push({ callback, options });
    return () => this.off(type, callback);
  }

  /**
   * @param {string} type
   * @param {(data: ServerEvent) => void} callback
   */
  off(type, callback) {
    let events = this.#listeners.get(type);
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
  #attach(webSocket) {
    webSocket.addEventListener(
      "open",
      () => {
        this.#log("connected");

        webSocket.addEventListener("close", (closeEvent) => {
          this.#log("connection closed", closeEvent);
          this.#reconenct();
        });
        webSocket.addEventListener("err", (errorEvent) => {
          this.#log("connection error", errorEvent);
          this.#reconenct();
        });
        webSocket.addEventListener("message", (messageEvent) => {
          try {
            this.#receiveEvent(messageEvent.data);
          } catch (error) {
            this.#log("failed to hanlde message", error, messageEvent);
          }
        });
      },
      { once: true },
    );
  }

  /** @param {unknown} data  */
  #receiveEvent(data) {
    if (typeof data !== "string") {
      return;
    }
    /** @type ServerEvent */
    const obj = JSON.parse(data);
    if (!("type" in obj)) {
      throw new Error("bad structure");
    }
    const listeners = this.#listeners.get(obj.type);
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
        this.off(obj.type, handler.callback);
      }
    }
  }

  /** @param {any[]} data  */
  #log(...data) {
    console.log(`[ServerEvents]`, ...data);
  }

  /** @throws if connection is null */
  async #reconenct() {
    if (this.#connection === null) {
      throw new Error("attempt to reconnect to non existing socket");
    }
    if (this.#connection.readyState !== WebSocket.CLOSED) {
      return;
    }
    try {
      this.#log("connection lost, reconnecting...");
      const res = await fetch(currentUrl(config.server.pathnames.healthCheck));
      if (res.status == 200) {
        this.#connection = new WebSocket(this.#connection?.url);
        this.#attach(this.#connection);
        return;
      }
      this.#reconenct();
    } catch {
      await new Promise((r) => setTimeout(r, this.#reconnectTimeout));
      this.#reconenct();
    }
  }
}
