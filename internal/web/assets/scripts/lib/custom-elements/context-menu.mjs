import { Bind, html } from "/assets/scripts/lib/libdim.mjs";

import { AssertInstance } from "../assert.mjs";
import { ShortcutManager } from "../shortcut-manager.mjs";

export class ContextMenu extends HTMLElement {
  #position = { x: 0, y: 0 };

  backdrop = new Bind((el) => AssertInstance.once(el, HTMLElement));
  element = new Bind((el) => AssertInstance.once(el, HTMLDialogElement));

  constructor() {
    super();
    this.style.setProperty("visibility", "hidden");
    this.attachShadow({ mode: "open" }).append(html`
      <style>
        :host {
          --_z-index: 1000;
        }
        * {
          color: var(--text-0);
          user-select: none;
        }

        dialog,
        button {
          background: var(--bg-2);
          border: 1px solid rgb(from var(--text-0) r g b / 0.25);
          font-size: 1rem;
        }

        #main {
          z-index: var(--_z-index);
          position: fixed;
          top: 0;
          left: 0;
          margin: 0;
          transform: translateX(min(var(--_x), calc(100vw - 100%)))
            translateY(min(var(--_y), calc(100vh - 100% - 1px)));
        }

        dialog {
          --_border-radius: 0.33rem;
          min-width: 15ch;
          padding: 0.25rem;
          border: 1px solid rgba(from var(--text-0) r g b / 0.5);
          border-radius: var(--_border-radius);
          outline: none;
        }

        button {
          border-radius: calc(var(--_border-radius) / 2);
          background: transparent;
          text-align: start;
          display: block;
          width: 100%;
          outline-color: transparent;
          border-color: transparent;
          &:hover {
            background: var(--primary);
          }
        }

        .group {
          position: relative;
          display: flex;
          justify-content: space-between;
          &:has(dialog:hover) {
            background: rgb(from var(--text-0) r g b / 0.1);
          }
        }

        hr {
          margin: 0.1rem 0.25rem;
          border-bottom: none;
        }

        #backdrop {
          z-index: calc(var(--_z-index) - 1);
          position: fixed;
          inset: 0;
        }
      </style>

      <div
        id="backdrop"
        onpointerdown="${() => this.close()}"
        bind="${this.backdrop}"
      ></div>

      <dialog
        id="main"
        tabindex="-1"
        class="translate"
        oncontextmenu="${(/** @type {unknown} */ e) => {
          const event = AssertInstance.once(e, Event);
          event.preventDefault();
          event.stopPropagation();
        }}"
        onclose="${() => this.style.setProperty("visibility", "hidden")}"
        bind="${this.element}"
      ></dialog>
    `);
  }

  /** @typedef {{name: string}} ItemBase */
  /** @typedef {ItemBase & {children: ItemGroup}} NestedItem */
  /** @typedef {{name: string} & {action: () => void}} ActionItem */
  /** @typedef {(NestedItem | ActionItem)} Item */
  /** @typedef {Item[][]} ItemGroup */

  #closeCleanup = () => {};

  /** @type {ItemGroup | null} */
  #currentContentData = null;
  /** @param {ItemGroup} data */
  open(data, position = this.#position) {
    if (this.#currentContentData !== data) {
      this.element.current.replaceChildren(this.#buildContent(data));
      this.#currentContentData = data;
    }
    this.element.current.show();
    this.element.current.focus();
    this.style.setProperty("visibility", "visible");
    this.element.current.style.setProperty("--_x", `${position.x}px`);
    this.element.current.style.setProperty("--_y", `${position.y}px`);
    this.#closeCleanup = ShortcutManager.keydown("<Escape>", () =>
      ContextMenu.instance.close(),
    );
  }

  close() {
    this.element.current.close();
    this.#closeCleanup();
  }

  /**
   * @param {ItemGroup} data
   * @returns {DocumentFragment}
   */
  #buildContent(data) {
    const fragment = new DocumentFragment();
    data.forEach((group, index) => {
      if (index > 0) {
        fragment.append(html`<hr />`);
      }

      group.forEach((item) => {
        if ("action" in item) {
          fragment.append(
            html`<button onclick="${item.action}">${item.name}</button>`,
          );
          return;
        }
        fragment.append(this.#buildNestedContent(item));
      });
    });
    return fragment;
  }

  /**
   * @param {NestedItem} item
   * @returns {DocumentFragment}
   * */
  #buildNestedContent = (item) => {
    const content = new Bind((el) =>
      AssertInstance.once(el, HTMLDialogElement),
    );
    const trigger = new Bind((el) => AssertInstance.once(el, HTMLElement));

    const onPointerLeave = () => {
      const backdrop = this.backdrop.current;
      ContextMenu.instance.addEventListener(
        "pointermove",
        async function onElement(event) {
          // @ts-expect-error fuck Firefox unsupported custom elements target
          const el = event.explicitOriginalTarget ?? event.target;
          if (
            !(el instanceof Node) ||
            el === backdrop ||
            trigger.current.contains(el)
          ) {
            return;
          }
          content.current.close();
          ContextMenu.instance.removeEventListener("pointermove", onElement);
        },
      );
    };

    const onPointerEnter = async () => {
      this.#openNestedDialog(trigger.current, content.current);
    };

    return html`
      <button
        class="group"
        bind="${trigger}"
        onpointerleave="${onPointerLeave}"
        onclick="${onPointerEnter}"
        onpointerenter="${onPointerEnter}"
      >
        ${item.name}<span>&rsaquo;</span>
        <dialog
          tabindex="-1"
          onclick="${(/** @type {unknown} */ e) =>
            AssertInstance.once(e, Event).stopPropagation()}"
          bind="${content}"
        >
          ${this.#buildContent(item.children)}
        </dialog>
      </button>
    `;
  };

  /**
   *@param {HTMLElement} trigger
   *@param {HTMLDialogElement} content
   */
  #openNestedDialog(trigger, content) {
    const triggerElement = AssertInstance.once(
      trigger,
      HTMLElement,
      "trigger needs to be an element for position calculations",
    );
    const contentElement = AssertInstance.once(
      content,
      HTMLDialogElement,
      "content needs to be dialog",
    );
    contentElement.show();
    contentElement.focus();
    const buttonBox = triggerElement.getBoundingClientRect();
    const dialogBox = contentElement.getBoundingClientRect();
    contentElement.style.setProperty("margin", "0");
    contentElement.style.setProperty(
      "translate",
      `var(--_transform-x) var(--_transform-y)`,
    );
    contentElement.style.setProperty(
      "--_transform-x",
      buttonBox.x + buttonBox.width + dialogBox.width < window.innerWidth
        ? `${buttonBox.width - 1}px`
        : `-100%`,
    );
    const verticalOverflow =
      window.innerHeight - (buttonBox.y + dialogBox.height) - 4;
    contentElement.style.setProperty(
      "--_transform-y",
      verticalOverflow > 0 ? "-7px" : `${verticalOverflow}px`,
    );
  }

  /** @type {ContextMenu} */
  static instance;
  connectedCallback() {
    if (ContextMenu.instance) {
      throw new Error("Only one context menu is allowed to exist");
    }
    ContextMenu.instance = this;
    window.addEventListener("contextmenu", (event) => {
      this.#position.x = event.x;
      this.#position.y = event.y;
    });
  }
}

customElements.define("h-context-menu", ContextMenu);
