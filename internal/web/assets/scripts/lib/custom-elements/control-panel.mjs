import { AssertInstance } from "../assert.mjs";
import { html } from "../html.mjs";
import { LocationControll } from "../location-control.mjs";
import { MovableList } from "../movable-list.mjs";
import { Publisher } from "../publisher.mjs";
import { ShortcutManager } from "../shortcut-manager.mjs";
import { TextAreaAutoresize } from "./textarea-autoresize.mjs";

export class Action {
  /**
   * @param {string} name
   * @param {() => (void | Promise<void>)} action
   */
  constructor(name, action) {
    this.name = name;
    this.action = action;
  }
}

const assertInput = new AssertInstance(TextAreaAutoresize);

export const controlPalanelVisibility = new Publisher(false);

export class ControlPanel extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanup = [];
  #visible = controlPalanelVisibility;
  #movableList;

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.shadow.innerHTML = this.#content;
    this.input = assertInput.check(this.shadow.querySelector("#input"));
    this.matchesContainer = AssertInstance.once(
      this.shadow.querySelector("#matches"),
      HTMLDivElement,
    );
    this.#movableList = new MovableList(
      this.matchesContainer,
      (current, previous) => {
        if (previous) {
          this.matchesContainer.children
            .item(previous)
            ?.classList.remove("under-cursor");
        }
        this.matchesContainer.children
          .item(current)
          ?.classList.add("under-cursor");
      },
    );
  }

  /** @param {Action[]} actions  */
  #updateMatchList(actions) {
    this.matchesContainer.innerHTML = "";
    const elements = actions.map((el) => {
      const anchor = document.createElement("a");
      anchor.onclick = async () => {
        await el.action();
        this.#visible.update(false);
      };
      anchor.innerHTML = el.name.replace(
        this.input.value,
        html`<span class="highlight">${this.input.value}</span>`,
      );
      return anchor;
    });
    this.matchesContainer.append(...elements);
    this.#movableList.cursor = 0;
  }

  connectedCallback() {
    this.style.visibility = "hidden";
    this.#visible.attach({
      notify: (visible) => {
        this.#updateMatchList(ActionStore.search(""));
        this.style.setProperty("visibility", visible ? "visible" : "hidden");
        this.input.focus();
        this.input.value = "";
      },
    });

    this.#hotkeys();
    this.#closeOnClick();
    this.input.addEventListener("input", () => {
      this.#updateMatchList(ActionStore.search(this.input.value));
    });
  }

  #closeOnClick() {
    AssertInstance.once(
      this.shadow.querySelector("#container"),
      HTMLDivElement,
    ).addEventListener("click", (e) => e.stopPropagation());
    AssertInstance.once(
      this.shadow.querySelector("#background"),
      HTMLDivElement,
    ).addEventListener("click", (e) => {
      this.#visible.update(false);
      e.stopPropagation();
    });
  }

  /** @type {(() => void)[]} */
  #visibleHotkeys = [];
  #hotkeys() {
    this.#visible.attach({
      notify: (visible) => {
        if (visible) {
          this.#visibleHotkeys.push(
            ShortcutManager.keydown(
              ["<Escape>", "<C-BracketLeft>", "<C-KeyC>"],
              (event) => {
                event.stopPropagation();
                event.preventDefault();
                this.#visible.update(false);
              },
              { priority: 999 },
            ),

            ShortcutManager.keydown(
              ["<Enter>", "<C-KeyM>"],
              (e) => {
                e.stopPropagation();
                e.preventDefault();
                const match = ActionStore.search(this.input.value).at(
                  this.#movableList.cursor,
                );
                if (match) {
                  match.action();
                  this.#visible.update(false);
                }
              },
              { priority: 990 },
            ),

            ShortcutManager.keydown(
              "<*>",
              (event) => {
                if (document.activeElement === this.input) {
                  return;
                }
                this.input.focus();
                event.stopPropagation();
              },
              { priority: 990 },
            ),

            ShortcutManager.keydown(
              [
                "<ArrowUp>",
                "<ArrowDown>",
                "<C-KeyN>",
                "<C-KeyP>",
                "<Tab>",
                "<S-Tab>",
              ],
              (event) => {
                /** @type {1 | -1} */
                let dir;
                switch (event.notation) {
                  case "<ArrowDown>":
                  case "<C-KeyN>":
                  case "<Tab>":
                    dir = 1;
                    break;
                  case "<ArrowUp>":
                  case "<C-KeyP>":
                  case "<S-Tab>":
                    dir = -1;
                    break;
                  default:
                    throw new Error("unhandled key");
                }
                this.#movableList.move(dir);
                event.stopPropagation();
                event.preventDefault();
              },
              { priority: 991 },
            ),
          );
          return;
        }
        this.#visibleHotkeys.forEach((cb) => cb());
      },
    });

    this.#cleanup.push(
      ShortcutManager.keydown(
        "<C-KeyP>",
        (event) => {
          if (!this.#visible.value) {
            this.#visible.update(true);
            event.preventDefault();
            event.stopPropagation();
          }
        },
        {
          priority: 999,
        },
      ),
    );
  }

  disconnectedCallback() {
    this.#cleanup.forEach((cb) => cb());
    this.#visibleHotkeys.forEach((cb) => cb());
  }

  #content = html`
    <style>
      #background {
        position: fixed;
        inset: 0;
        background: rgb(from var(--bg-0) r g b / 0.75);
        z-index: 99999;
      }

      #container {
        position: absolute;
        left: 50%;
        top: 10%;
        translate: -50% 0;
        width: 60%;
        display: flex;
        flex-direction: column;
        align-items: center;
        border-radius: 1rem;
        overflow: hidden;
        padding: 0 calc(1rem + 1px);
        background: var(--bg-2);
        border: 1px solid rgb(from var(--text-0) r g b / 0.25);
      }

      #input {
        padding: 1rem;
        width: 100%;
        margin: 0;
        color: var(--text-0);
        resize: none;
        outline: none;
        background: transparent;
        border: none;
      }

      #matches {
        display: grid;
        gap: 0.1rem;
        padding: 1rem;
        width: 100%;
        border-top: 1px solid rgb(from var(--text-0) r g b / 0.25);
        max-height: 60vh;
        overflow: auto;

        & > * {
          scroll-margin: 1rem;
        }
      }

      a {
        cursor: pointer;
        padding: 0.5rem 1rem;
        border-radius: 0.5rem;
        border: 1px solid transparent;
      }

      a:is(:hover, :focus) {
        background: rgb(from var(--text-0) r g b / 0.1);
      }

      .highlight {
        color: var(--primary);
      }

      .under-cursor {
        border-color: var(--primary);
      }
    </style>

    <div id="background">
      <div id="container">
        <textarea
          id="input"
          is="hermes-textarea-autoresize"
          focus-on-input="true"
          max-rows="1"
          placeholder="Control Panel"
          autofocus
          required
        ></textarea>
        <div id="matches"></div>
      </div>
    </div>
  `;
}

export class ActionStore {
  /** @type {Set<Action[]>} */
  static #all = new Set();

  /**
   * @param {Action[]} actions
   * @returns {() => void} teardown
   */
  static add(...actions) {
    ActionStore.#all.add(actions);
    return () => {
      ActionStore.#all.delete(actions);
    };
  }

  /** @returns {number} */
  static get size() {
    return ActionStore.#all.values().toArray().flat().length;
  }

  /**
   * @param {string} query
   * @returns {Action[]}
   */
  static search(query) {
    return ActionStore.#all
      .values()
      .toArray()
      .flat()
      .filter((action) => action.name.includes(query));
  }
}

ActionStore.add(
  new Action("create: new chat", () => LocationControll.navigate("/")),
  new Action("open: last chat", async () => {
    const res = await fetch(`/api/v1/chats?limit=1`);
    const data = await res.json();
    LocationControll.navigate("/chats/" + data[0].id);
  }),
  new Action("open: chats history", () => LocationControll.navigate("/chats")),
);
