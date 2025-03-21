import { Bind, html } from "/assets/scripts/lib/libdim.mjs";

import { AssertInstance } from "../assert.mjs";
import { LocationControll } from "../location-control.mjs";
import { MovableList } from "../movable-list.mjs";
import { Publisher } from "../publisher.mjs";
import { ShortcutManager } from "../shortcut-manager.mjs";
import { stringMatching } from "../string-matching.mjs";
import { withCache } from "../with-cache.mjs";
import { ResizableTextInput } from "./content-editable-plain-text.mjs";
import { HermesDialog } from "./dialog.mjs";

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

export const controlPalanelVisibility = new Publisher(false);

const stringMatchingWithCache = withCache(stringMatching);

export class ControlPanel extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanup = [];
  #visible = controlPalanelVisibility;
  #movableList;
  #input = new Bind((el) => AssertInstance.once(el, ResizableTextInput));

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.shadow.replaceChildren(this.#content);
    this.matchesContainer = AssertInstance.once(
      this.shadow.querySelector("#matches"),
      HTMLDivElement,
    );
    this.#movableList = new MovableList(
      this.matchesContainer,
      (current, previous) => {
        if (previous !== undefined) {
          this.matchesContainer.children
            .item(previous)
            ?.classList.remove("under-cursor");
        }
        const cur = this.matchesContainer.children.item(current);
        cur?.classList.add("under-cursor");
        cur?.scrollIntoView({ block: "nearest" });
      },
    );
  }

  /** @param {Action[]} actions  */
  #updateMatchList(actions) {
    const elements = actions.map((el) => {
      const anchor = document.createElement("a");
      anchor.onclick = async () => {
        await el.action();
        this.#visible.update(false);
      };
      anchor.replaceChildren(
        ...el.name.split("").map((char, i) => {
          const { ok, matches } = stringMatchingWithCache(
            el.name,
            this.#input.current.value.trim(),
          );
          if (!ok || !matches[i]) {
            return char;
          }
          return html`<span class="highlight">${char}</span>`;
        }),
      );
      return anchor;
    });
    this.matchesContainer.replaceChildren(...elements);
    this.#movableList.cursor = 0;
  }

  connectedCallback() {
    AssertInstance.once(
      this.shadow.querySelector("h-dialog"),
      HermesDialog,
      "expected contorl panel to render inside of hermes dialog",
    ).element.addEventListener("close", () => this.#visible.update(false));

    this.#cleanup.push(
      this.#visible.subscribe({
        notify: () => {
          this.#updateMatchList(ActionStore.search(""));
          this.#input.current.content.focus();
          this.#input.current.value = "";
        },
      }),
      this.#visible.subscribe({
        notify: (visible) => {
          const dialog = AssertInstance.once(
            this.shadow.querySelector("h-dialog"),
            HermesDialog,
            "expected contorl panel to render inside of hermes dialog",
          );
          if (visible) {
            dialog.element.showModal();
          } else {
            dialog.element.close();
          }
        },
      }),
    );

    this.#hotkeys();
    this.#input.current.addEventListener("change", () => {
      this.#updateMatchList(ActionStore.search(this.#input.current.value.trim()));
    });
  }

  /** @type {(() => void)[]} */
  #visibleHotkeys = [];
  #hotkeys() {
    this.#visible.subscribe({
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
                const match = ActionStore.search(this.#input.current.value.trim()).at(
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
                if (document.activeElement === this.#input.current) {
                  return;
                }
                this.#input.current.content.focus();
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
      * {
        color: var(--text-0);
        box-sizing: border-box;
      }

      h-dialog-card {
        display: block;
        width: min(80ch, 80vw);
      }

      h-resizable-text-input {
        --padding: 1rem;
        &::part(wrapper) {
          padding: var(--padding);
          color: var(--text-0);
        }
        &::part(placeholder) {
          padding: var(--padding);
        }
        &::part(content) {
          outline: none;
          border: none;
        }
      }

      #matches {
        display: grid;
        gap: 0.1rem;
        padding: 1rem;
        border-top: 1px solid rgb(from var(--text-0) r g b / 0.25);
        max-height: 60vh;
        overflow: auto;

        &:is(:empty):before {
          color: rgb(from var(--text-0) r g b / 0.5);
          content: "no matching actions";
        }

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

    <h-dialog dialog-style="margin-top: 10vh;">
      <h-dialog-card section-style="max-width: unset">
        <div id="content">
          <h-resizable-text-input
            bind="${this.#input}"
            placeholder="Control Panel"
          ></h-resizable-text-input>
          <div id="matches"></div>
        </div>
      </h-dialog-card>
    </h-dialog>
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
      .filter((action) => stringMatchingWithCache(action.name, query).ok);
  }
}

ActionStore.add(
  new Action("create: new chat", () => LocationControll.navigate("/")),
  new Action("create: new template", () =>
    LocationControll.navigate("/templates/new"),
  ),
  new Action("open: last chat", async () => {
    const res = await fetch(`/api/v1/chats?limit=1`);
    const data = await res.json();
    LocationControll.navigate("/chats/" + data[0].id);
  }),
  new Action("open: chats history", () => LocationControll.navigate("/chats")),
  new Action("open: templates history", () =>
    LocationControll.navigate("/templates"),
  ),
);
