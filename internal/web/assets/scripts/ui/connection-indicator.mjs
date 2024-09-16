import { ServerEvents } from "../events/server-events.mjs";
import { Publisher } from "../utils/publisher.mjs";

class OnlineObserver {
  #offlineSufix = " - offline";

  /** @param {boolean} isOnline */
  notify(isOnline) {
    if (isOnline) {
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
  const statusPublisher = new Publisher(ServerEvents.isOpen);
  ServerEvents.onClose(() => statusPublisher.update(ServerEvents.isOpen));
  ServerEvents.onOpen(() => statusPublisher.update(ServerEvents.isOpen));
  const onlineObserver = new OnlineObserver();
  statusPublisher.attach(onlineObserver);
  statusPublisher.notify();
}
