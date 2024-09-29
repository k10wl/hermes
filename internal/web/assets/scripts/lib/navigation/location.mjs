import { Publisher } from "/assets/scripts/utils/publisher.mjs";

export class LocationControll {
  static #publisher = new Publisher(LocationControll.#getPathname());
  static #ready = false;

  /**
   * @param {{notify: (route: string) => void}} observer
   */
  static attach(observer) {
    observer.notify(LocationControll.#getPathname());
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

  static #getPathname() {
    return window.location.pathname;
  }

  static #update() {
    LocationControll.#publisher.update(LocationControll.#getPathname());
  }
}

LocationControll.__init();
