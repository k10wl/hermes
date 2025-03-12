import { Bind, escapeMarkup, html } from "/assets/scripts/lib/libdim.mjs";

import { AssertInstance } from "../../assert.mjs";
import { RequestReadTemplatesEvent } from "../../events/client-events-list.mjs";
import { ServerEvents } from "../../events/server-events.mjs";
import { ServerErrorEvent } from "../../events/server-events-list.mjs";

customElements.define(
  "hermes-templates-list-scene",
  class extends HTMLElement {
    /** @type {(() => void)[]} */
    #cleanup = [];

    /** @type {Map<string, HTMLElement>} */
    #elements = new Map();

    templatesContainer = new Bind((el) => AssertInstance.once(el, HTMLElement));
    newTemplate = new Bind((el) => AssertInstance.once(el, HTMLElement));

    constructor() {
      super();
      this.shadow = this.attachShadow({ mode: "open" });
      this.shadow.append(html`
        <style>
          * {
            box-sizing: border-box;
          }

          main {
            height: 100%;
            overflow: auto;
          }

          section {
            padding: 2rem;
            max-width: 100%;
            overflow: hidden;
            display: flex;
            flex-direction: column;
          }

          a {
            /* NOTE uuuuh this could be separate reusable class or component */
            text-align: start;
            padding: 0.5rem 1rem;
            margin: 0.1rem 0rem;
            border: 1px solid rgb(from var(--text-0) r g b / 0.25);
            text-decoration: none;
            color: rgb(from var(--text-0) r g b / 0.5);
            transition: border-color var(--color-transition-duration);
            display: flex;

            &:hover {
              border-color: var(--primary);
            }

            .name {
              flex-shrink: 0;
              color: var(--text-0);
            }

            .content {
              padding: 0;
              margin: 0;
              overflow: hidden;
              text-overflow: ellipsis;
              white-space: nowrap;
            }
          }
        </style>

        <main>
          <section bind="${this.templatesContainer}">
            <a
              bind="${this.newTemplate}"
              is="hermes-link"
              href="/templates/new"
            >
              <span class="name"> // Create new template </span>
            </a>
          </section>
        </main>
      `);
    }

    connectedCallback() {
      const readTemplates = new RequestReadTemplatesEvent({
        name: "",
        limit: -1,
        start_before_id: -1,
      });

      const offRead = ServerEvents.on(
        ["read-templates", "server-error"],
        (event) => {
          if (event.id !== readTemplates.id) {
            return;
          }
          offRead();
          if (event instanceof ServerErrorEvent) {
            alert(`smth went wrong - ${event.payload}`);
            return;
          }
          this.templatesContainer.current.append(
            ...event.payload.templates.map((template) =>
              this.#createLink(template),
            ),
          );
        },
      );
      ServerEvents.send(readTemplates);

      this.#cleanup.push(
        offRead,
        () => this.#elements.clear(),
        ServerEvents.on("template-created", (event) => {
          this.newTemplate.current.after(
            this.#createLink(event.payload.template),
          );
        }),
        ServerEvents.on("template-changed", (event) => {
          const el = this.#elements.get(event.payload.template.name);
          if (!el) {
            return;
          }
          el.replaceChildren(this.#linkContents(event.payload.template));
        }),
        ServerEvents.on("template-deleted", (event) => {
          this.#elements.get(event.payload.name)?.remove();
        }),
      );
    }

    /**
     * @param {import("/assets/scripts/models.mjs").Template} template
     * @returns {DocumentFragment}
     */
    #linkContents(template) {
      return html`
        <span class="name">${template.name}</span>:&nbsp;
        <p class="content">
          ${escapeMarkup(template.content.replaceAll("\n", " "))}
        </p>
      `;
    }

    /**
     * @param {import("/assets/scripts/models.mjs").Template} template
     * @returns {DocumentFragment}
     */
    #createLink(template) {
      const link = new Bind((el) => AssertInstance.once(el, HTMLAnchorElement));
      const fragment = html`
        <a bind="${link}" is="hermes-link" href="/templates/${template.id}">
          ${this.#linkContents(template)}
        </a>
      `;
      this.#elements.set(template.name, link.current);
      return fragment;
    }

    disconnectedCallback() {
      this.#cleanup.forEach((cb) => cb());
    }
  },
);
