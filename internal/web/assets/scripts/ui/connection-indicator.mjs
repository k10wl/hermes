import { ServerEvents } from "../events/server-events.mjs";
import { Publisher } from "../utils/publisher.mjs";

class OnlineObserver {
  #offlineSufix = " - offline";

  /** @param {boolean} connected */
  notify(connected) {
    if (connected) {
      window.document.title = window.document.title.replaceAll(
        this.#offlineSufix,
        "",
      );
      return;
    }
    if (window.document.title.endsWith(this.#offlineSufix)) {
      return;
    }
    window.document.title += this.#offlineSufix;
  }
}

export function initConnectionIndicator() {
  const statusPublisher = new Publisher(ServerEvents.connected);
  ServerEvents.onClose(() => statusPublisher.update(ServerEvents.connected));
  ServerEvents.onOpen(() => statusPublisher.update(ServerEvents.connected));
  const onlineObserver = new OnlineObserver();
  statusPublisher.attach(onlineObserver);
  statusPublisher.notify();
}
