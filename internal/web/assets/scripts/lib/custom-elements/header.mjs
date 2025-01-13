import { AssertInstance } from "../assert.mjs";
import { html } from "../html.mjs";
import { LocationControll } from "../location-control.mjs";
import { Publisher } from "../publisher.mjs";
import { Action, ActionStore } from "./control-panel.mjs";

export const headerVisiblePublisher = new Publisher(false);

class ActionStoreUpdater {
  #teardown = () => {};
  /** @param {boolean} isVisible  */
  notify(isVisible) {
    this.#teardown();
    this.#teardown = ActionStore.add(
      new Action(`${isVisible ? "hide" : "show"} header`, () =>
        headerVisiblePublisher.update(!isVisible),
      ),
    );
  }
}

class ActiveLinkObserver {
  #container;
  /**
   * @param {HTMLElement} container
   */
  constructor(container) {
    this.#container = container;
  }
  /**
   * @param {string} current
   * @param {string} [previous]
   */
  notify(current, previous) {
    for (const element of this.#container.querySelectorAll(
      `:is([href="${current}"], [href="${previous}"])`,
    )) {
      const link = AssertInstance.once(element, HTMLAnchorElement);
      if (link.href.endsWith(current)) {
        link.classList.add("active");
      } else {
        link.classList.remove("active");
      }
    }
  }
}

customElements.define(
  "hermes-header",
  class Header extends HTMLElement {
    /** @type {(() => void)[]} */
    #cleanup = [];
    /** @type {{name: string, href: string}[]} */
    static #links = [
      {
        name: "new chat",
        href: "/",
      },
      {
        name: "chats history",
        href: "/chats",
      },
      {
        name: "templates",
        href: "/templates",
      },
    ];

    constructor() {
      super();
      this.shadow = this.attachShadow({ mode: "open" });
      this.shadow.innerHTML = html`
        <style>
          :host {
            --button-size: var(--header-height);
            --default-color: rgb(from var(--text-0) r g b / 0.33);
            --active-color: rgb(from var(--text-0) r g b / 0.66);
            --hovered-color: var(--text-0);
          }

          header {
            &[aria-expanded="false"] {
              translate: 0 -100%;
              visibility: hidden;
            }
            z-index: 10;
            background-color: var(--bg-0);
            translate: 0;
            height: var(--header-height);
            transition:
              translate 100ms,
              visibility 100ms;
            position: fixed;
            top: 0;
            left: 0;
            right: 0;
            display: flex;
            gap: 1rem;
            align-items: center;
            font-size: 1rem;
            border-bottom: 1px solid var(--text-0);
          }

          * {
            color: var(--default-color);
          }

          :is(button, a, button::before) {
            &:hover {
              color: var(--hovered-color) !important;
            }
          }

          button {
            width: var(--button-size);
            height: var(--button-size);
            position: relative;
            cursor: pointer;
            background: transparent;
            border-color: transparent;
            outline-color: transparent;
            &#open {
              position: fixed;
              top: 0;
              left: 0;
              [aria-expanded="true"] + & {
                visibility: hidden;
              }
            }
          }

          a {
            text-decoration: none;
            &.active {
              color: var(--active-color);
              text-decoration: underline;
            }
          }
        </style>

        <header>
          <button id="close">✕</button>
          <nav>
            ${Header.#links
              .map(
                (link) => html`
                  <a
                    href="${link.href}"
                    is="hermes-link"
                    ${LocationControll.pathname === link.href
                      ? "class='active'"
                      : ""}
                    >${link.name}</a
                  >
                `,
              )
              .join("|")}
          </nav>
        </header>

        <button id="open">☰</button>
      `;
    }

    connectedCallback() {
      this.header = AssertInstance.once(
        this.shadow.querySelector("header"),
        HTMLElement,
      );
      this.open = AssertInstance.once(
        this.shadow.getElementById("open"),
        HTMLButtonElement,
      );
      this.open.addEventListener("click", () => {
        headerVisiblePublisher.update(true);
      });
      const close = AssertInstance.once(
        this.shadow.getElementById("close"),
        HTMLButtonElement,
      );
      close.addEventListener("click", () => {
        headerVisiblePublisher.update(false);
      });

      const linkObserver = new ActiveLinkObserver(
        AssertInstance.once(this.shadow.querySelector("nav"), HTMLElement),
      );
      linkObserver.notify(LocationControll.pathname, "");

      const actionStoreUpdater = new ActionStoreUpdater();

      actionStoreUpdater.notify(headerVisiblePublisher.value);
      this.notify(headerVisiblePublisher.value);

      this.#cleanup.push(
        LocationControll.attach(linkObserver),
        headerVisiblePublisher.subscribe(this),
        headerVisiblePublisher.subscribe(actionStoreUpdater),
      );
    }

    disconnectedCallback() {
      this.#cleanup.forEach((cb) => cb());
    }

    /**
     * @param {boolean} current
     */
    notify(current) {
      AssertInstance.once(this.header, HTMLElement).ariaExpanded = current
        ? "true"
        : "false";
    }
  },
);
