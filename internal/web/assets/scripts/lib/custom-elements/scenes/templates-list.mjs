import { AssertInstance } from "../../assert.mjs";
import { RequestReadTemplatesEvent } from "../../events/client-events-list.mjs";
import { ServerEvents } from "../../events/server-events.mjs";
import { ServerErrorEvent } from "../../events/server-events-list.mjs";
import { html } from "../../html-v2.mjs";

customElements.define(
  "hermes-templates-list-scene",
  class extends HTMLElement {
    /** @type {(() => void)[]} */
    #cleanup = [];

    /** @type {Map<string, HTMLElement>} */
    #elements = new Map();

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
          <section
            bind="${(/** @type {unknown} */ element) =>
              (this.templatesContainer = AssertInstance.once(
                element,
                HTMLElement,
              ))}"
          >
            <a
              bind="${(/** @type {unknown} */ element) =>
                (this.newTemplate = AssertInstance.once(element, HTMLElement))}"
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
          this.templatesContainer.append(
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
          this.newTemplate.after(this.#createLink(event.payload.template));
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
        <p
          class="content"
          bind="${(/** @type {unknown} */ element) => {
            AssertInstance.once(element, HTMLElement).innerText =
              template.content.replaceAll("\n", " ");
          }}"
        ></p>
      `;
    }

    /**
     * @param {import("/assets/scripts/models.mjs").Template} template
     * @returns {DocumentFragment}
     */
    #createLink(template) {
      return html`
        <a
          bind="${(/** @type {unknown} */ element) => {
            const asserted = AssertInstance.once(element, HTMLElement);
            // sometimes templates contain HTML, needs not to be interpreted
            asserted.replaceChildren(this.#linkContents(template));
            this.#elements.set(template.name, asserted);
          }}"
          is="hermes-link"
          href="/templates/${template.id}"
        >
        </a>
      `;
    }

    disconnectedCallback() {
      this.#cleanup.forEach((cb) => cb());
    }
  },
);
