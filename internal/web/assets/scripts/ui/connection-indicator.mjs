import { ServerEvents } from "/assets/scripts/lib/events/server-events.mjs";
import { Publisher } from "/assets/scripts/lib/publisher.mjs";

import { AssertInstance } from "../lib/assert.mjs";

const link = AssertInstance.once(
  document.querySelector("link[rel=icon]"),
  HTMLLinkElement,
);

const icons = {
  online: "/assets/icons/caduceus-staff.svg",
  offline: "/",
};

class OnlineObserver {
  #offlineSufix = " - offline";

  /** @param {boolean} connected */
  notify(connected) {
    if (connected) {
      window.document.title = window.document.title.replaceAll(
        this.#offlineSufix,
        "",
      );
      link.href = icons.online;
      return;
    }
    if (window.document.title.endsWith(this.#offlineSufix)) {
      return;
    }
    window.document.title += this.#offlineSufix;
    link.href = icons.offline;
  }
}

(() => {
  const statusPublisher = new Publisher(ServerEvents.connected);
  ServerEvents.on("connection-status-change", (data) =>
    statusPublisher.update(data.payload.connected),
  );
  const onlineObserver = new OnlineObserver();
  statusPublisher.subscribe(onlineObserver);
  statusPublisher.notify();
})();
