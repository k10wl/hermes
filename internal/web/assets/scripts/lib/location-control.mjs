import { Publisher } from "/assets/scripts/lib/publisher.mjs";

export class LocationControll {
  static get pathname() {
    return window.location.pathname;
  }

  static #publisher = new Publisher(LocationControll.pathname);
  static #ready = false;

  /**
   * @param {Parameters<Publisher<string>['subscribe']>[0]} observer
   */
  static attach(observer) {
    return this.#publisher.subscribe(observer);
  }

  /**
   * @param {string} target
   * @param {boolean} [notify = true]
   */
  static navigate(target, notify = true) {
    if (
      window.location.href === target ||
      window.location.pathname === target
    ) {
      return;
    }
    window.history.pushState({}, "", target);
    if (notify) {
      LocationControll.#update();
    }
  }

  static __init() {
    if (LocationControll.#ready) {
      return;
    }
    LocationControll.#ready = true;
    window.addEventListener("popstate", () => {
      LocationControll.#update();
    });
  }

  static #update() {
    LocationControll.#publisher.update(LocationControll.pathname);
  }

  static #templateIdRegexp = /\/templates\/(?<templateId>\d+)/;
  static get templateId() {
    try {
      const templateId = +(
        LocationControll.#templateIdRegexp.exec(LocationControll.pathname)
          ?.groups?.templateId ?? "no"
      );
      if (Number.isNaN(templateId)) {
        return null;
      }
      return templateId;
    } catch {
      return null;
    }
  }

  static #chatIdRegexp = /\/chats\/(?<chatId>\d+)/;
  static get chatId() {
    try {
      const chatId = +(
        LocationControll.#chatIdRegexp.exec(LocationControll.pathname)?.groups
          ?.chatId ?? "no"
      );
      if (Number.isNaN(chatId)) {
        return null;
      }
      return chatId;
    } catch {
      return null;
    }
  }
}

LocationControll.__init();
