import { AssertInstance } from "../assert.mjs";
import { html } from "../html-v2.mjs";
import { LocationControll } from "../location-control.mjs";
import { ContextMenu } from "./context-menu.mjs";
import { controlPalanelVisibility } from "./control-panel.mjs";
import { HermesViewTemplateScene, template } from "./scenes/view-template.mjs";

customElements.define(
  "hermes-header",
  class Header extends HTMLElement {
    /** @type {(() => void)[]} */
    #cleanup = [];

    #contextNavigation = [
      {
        name: "New Chat",
        action: () => {
          LocationControll.navigate("/");
          ContextMenu.instance.close();
        },
      },
      {
        name: "History",
        action: () => {
          LocationControll.navigate("/chats");
          ContextMenu.instance.close();
        },
      },
      {
        name: "Templates",
        action: () => {
          LocationControll.navigate("/templates");
          ContextMenu.instance.close();
        },
      },
    ];

    /** @type {Parameters<typeof ContextMenu.instance.open>[0][number]} */
    #contextLocationSpecific = [];

    #contextShared = [
      {
        name: "Control Panel",
        action: () => {
          controlPalanelVisibility.update(true);
          ContextMenu.instance.close();
        },
      },
    ];

    /**
     * @param {MouseEvent} e
     */
    #openContext = (e) => {
      /** @type {Parameters<typeof ContextMenu.instance.open>[0]} */
      const actions = [];
      actions.push(this.#contextNavigation);
      if (this.#contextLocationSpecific.length > 0) {
        actions.push(this.#contextLocationSpecific);
      }
      actions.push(this.#contextShared);
      const { x, y, height } = AssertInstance.once(
        e.target,
        HTMLElement,
      ).getBoundingClientRect();
      ContextMenu.instance.open(actions, {
        x,
        y: y + height,
      });
    };

    constructor() {
      super();
      this.shadow = this.attachShadow({ mode: "open" });
      this.shadow.append(html`
        <style>
          button {
            font-size: 1.25rem;
            position: fixed;
            top: 0;
            left: 0;
            z-index: 10;
            cursor: pointer;
            outline-color: transparent;
            border-color: transparent;
            background-color: transparent;
            color: rgb(from var(--text-0) r g b / 0.5);
            transition: color var(--color-transition-duration);
            &:hover {
              color: var(--text-0);
            }
          }
        </style>

        <header>
          <button onclick="${this.#openContext}">☰</button>
        </header>
      `);
    }

    connectedCallback() {
      this.notify();
      this.#cleanup.push(LocationControll.attach(this));
    }

    disconnectedCallback() {
      this.#cleanup.forEach((cb) => cb());
    }

    notify() {
      if (LocationControll.templateId !== null) {
        this.#contextLocationSpecific = [
          {
            name: "Save Template",
            action: () => {
              AssertInstance.once(
                template,
                HermesViewTemplateScene,
                "scene must share template and make it accessible",
              ).submit();
              ContextMenu.instance.close();
            },
          },
          {
            name: "Delete Template",
            action: () => {
              AssertInstance.once(
                template,
                HermesViewTemplateScene,
                "scene must share template and make it accessible",
              ).delete();
              ContextMenu.instance.close();
            },
          },
        ];
        return;
      }
      this.#contextLocationSpecific = [];
    }
  },
);
