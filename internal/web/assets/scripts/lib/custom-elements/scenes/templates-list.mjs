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

          #templates {
            padding: 2rem;
            max-width: 100%;
            overflow: hidden;
            display: flex;
            flex-direction: column;
          }

          #templates a {
            /* NOTE uuuuh this could be separate reusable class or component */
            text-align: start;
            padding: 0.5rem 1rem;
            margin: 0.1rem 0rem;
            border: 1px solid rgb(from var(--text-0) r g b / 0.25);
            text-decoration: none;
            color: var(--text-0);
            transition: border-color var(--color-transition-duration);
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
            &:hover {
              border-color: var(--primary);
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
            id="templates"
          ></section>
        </main>
      `);
    }

    connectedCallback() {
      this.templates = AssertInstance.once(
        this.shadow.querySelector("#templates"),
        HTMLElement,
      );

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
          AssertInstance.once(this.templates, HTMLElement).append(
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
          this.templatesContainer.prepend(
            this.#createLink(event.payload.template),
          );
        }),
        ServerEvents.on("template-changed", (event) => {
          const el = this.#elements.get(event.payload.template.name);
          if (!el) {
            return;
          }
          el.textContent = this.#linkText(event.payload.template);
        }),
        ServerEvents.on("template-deleted", (event) => {
          this.#elements.get(event.payload.name)?.remove();
        }),
      );
    }

    /**
     * @param {import("/assets/scripts/models.mjs").Template} template
     * @returns {string}
     */
    #linkText(template) {
      return `${template.name}: ${template.content}`;
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
            asserted.textContent = this.#linkText(template);
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
