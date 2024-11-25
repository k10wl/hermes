import { Publisher } from "/assets/scripts/utils/publisher.mjs";

export class LocationControll {
  static get pathname() {
    return window.location.pathname;
  }

  static #publisher = new Publisher(LocationControll.pathname);
  static #ready = false;

  /**
   * @param {{notify: (route: string) => void}} observer
   */
  static attach(observer) {
    observer.notify(LocationControll.pathname);
    return this.#publisher.attach(observer);
  }

  /** @param {string} href  */
  static navigate(href) {
    if (window.location.href === href) {
      return;
    }
    window.history.pushState({}, "", href);
    LocationControll.#update();
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
}

LocationControll.__init();
