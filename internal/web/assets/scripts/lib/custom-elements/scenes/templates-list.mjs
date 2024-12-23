import { AssertInstance } from "../../assert.mjs";
import { RequestReadTemplatesEvent } from "../../events/client-events-list.mjs";
import { ServerEvents } from "../../events/server-events.mjs";
import { html } from "../../html.mjs";

customElements.define(
  "hermes-templates-list-scene",
  class extends HTMLElement {
    /** @type {(() => void)[]} */
    #cleanup = [];

    constructor() {
      super();
      this.shadow = this.attachShadow({ mode: "open" });
      this.shadow.innerHTML = html`
        <style>
          main {
            height: 100%;
            display: grid;
            place-items: center;
            text-align: center;
          }

          section {
            display: grid;
            height: 100%;
          }
        </style>

        <main>
          <section id="templates"></section>
        </main>
      `;
    }

    connectedCallback() {
      this.templates = AssertInstance.once(
        this.shadow.querySelector("#templates"),
        HTMLElement,
      );
      ServerEvents.on("read-templates", (event) => {
        AssertInstance.once(this.templates, HTMLElement).innerHTML =
          event.payload.templates
            .map(
              (template) => html`
                <a is="hermes-link" href="/templates/${template.id}">
                  ${template.name}: ${template.content.slice(0, 79)}
                </a>
              `,
            )
            .join("");
      });
      ServerEvents.send(
        new RequestReadTemplatesEvent({
          name: "",
          limit: -1,
          start_before_id: -1,
        }),
      );
    }

    disconnectedCallback() {
      this.#cleanup.forEach((cb) => cb());
    }
  },
);
