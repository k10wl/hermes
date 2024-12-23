import { AssertInstance, AssertString } from "../../assert.mjs";
import { RequestReadTemplateEvent } from "../../events/client-events-list.mjs";
import { ServerEvents } from "../../events/server-events.mjs";
import { ServerErrorEvent } from "../../events/server-events-list.mjs";
import { html } from "../../html.mjs";
import { LocationControll } from "../../location-control.mjs";
import { TextAreaAutoresize } from "../textarea-autoresize.mjs";

customElements.define(
  "hermes-view-template-scene",
  class extends HTMLElement {
    constructor() {
      super();
      this.shadow = this.attachShadow({ mode: "closed" });
      this.shadow.innerHTML = this.#html;
    }

    connectedCallback() {
      const readEvent = new RequestReadTemplateEvent({
        id: parseInt(
          AssertString.check(
            LocationControll.pathname.match(/\d+$/)?.[0],
            "pathname should have id",
          ),
          10,
        ),
      });
      ServerEvents.send(readEvent);
      const off = ServerEvents.on(
        ["read-template", "server-error"],
        (event) => {
          if (event.id !== readEvent.id) {
            return;
          }
          off();
          if (event instanceof ServerErrorEvent) {
            // TODO show user that something exploded
            LocationControll.navigate("/templates");
            console.error("event");
            return;
          }
          const textarea = AssertInstance.once(
            this.shadow.querySelector("textarea"),
            TextAreaAutoresize,
            "expected resizable text area",
          );
          textarea.value = event.payload.template.content;
          textarea.autoresize();
        },
      );
    }

    #html = html`
      <style>
        main {
          height: 100%;
          display: grid;
          place-items: center;
        }

        textarea {
          width: 60ch;
          border: none;
          outline: none;
          background: transparent;
          color: var(--text-0);
          resize: none;
        }
      </style>

      <main>
        <textarea is="hermes-textarea-autoresize" autofocus></textarea>
      </main>
    `;
  },
);
