import { AssertInstance, AssertNumber, AssertString } from "../assert.mjs";
import { CreateCompletionMessageEvent } from "../events/client-events-list.mjs";
import { ServerEvents } from "../events/server-events.mjs";
import {
  ChatCreatedEvent,
  ServerErrorEvent,
} from "../events/server-events-list.mjs";
import { html } from "../html-v2.mjs";
import { LocationControll } from "../location-control.mjs";

export class MessageForm extends HTMLElement {
  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
  }

  /** @type {HTMLTextAreaElement | null} */
  #textarea = null;

  connectedCallback() {
    this.#render();

    const form = AssertInstance.once(
      this.shadow.querySelector("form"),
      HTMLFormElement,
    );

    form.addEventListener("submit", (e) => {
      e.preventDefault();
      /** @type {string | number | undefined} */
      let chat_id = LocationControll.pathname.split("/").at(-1);
      chat_id = AssertNumber.check(chat_id ? +chat_id : -1);
      const content = AssertString.check(new FormData(form).get("content"));
      if (content.trim() === "") {
        return;
      }
      const message = new CreateCompletionMessageEvent({
        chat_id,
        content: content,
        parameters: {
          model: "openai/gpt-4o-mini",
          max_tokens: undefined,
          temperature: undefined,
        },
      });
      ServerEvents.send(message);
      const off = ServerEvents.on(
        ["chat-created", "message-created", "server-error"],
        (event) => {
          if (event.id !== message.id) {
            return;
          }
          off();
          if (event instanceof ChatCreatedEvent) {
            LocationControll.navigate(`/chats/${event.payload.chat.id}`);
          }
          if (event instanceof ServerErrorEvent) {
            return;
          }
          form.reset();
        },
      );
    });
  }

  #render() {
    this.shadow.append(html`
      <style>
        :host {
          --bg: var(--bg-2);
          --text: var(--text-0);
          --radius: 1rem;
        }

        form {
          display: flex;
          gap: 0.5rem;
          border-radius: var(--radius);
          padding: 0 calc(var(--radius) / 2) 0 var(--radius);
          background: var(--bg);
          color: var(--text);
        }

        .actions {
          height: 3rem;
          display: flex;
          justify-content: flex-end;
          align-items: center;
          align-self: flex-end;
          gap: 0.5rem;
        }

        textarea {
          background: transparent;
          color: var(--text);
          padding: 0;
          padding-top: var(--radius);
          max-height: 50vh;
          width: 100%;
          margin: 0;
          resize: none;
          outline: none;
          border: none;
        }

        form:has(textarea:invalid) button[type="submit"] {
          --_col: rgb(from var(--text) r g b / 0.25);
          background: var(--bg);
          color: var(--_col);
          border-color: var(--_col);
          cursor: auto;
        }

        button {
          --_size: 2rem;
          transition: all var(--color-transition-duration);
          flex-shrink: 0;
          background: var(--primary);
          font-size: calc(var(--_size) * 0.66);
          color: var(--text);
          outline-color: transparent;
          border-color: transparent;
          border-radius: var(--_size);
          width: var(--_size);
          height: var(--_size);
          cursor: pointer;
        }
      </style>

      <form
        onclick="${() => {
          return this.#textarea?.focus();
        }}"
        is="hermes-form"
      >
        <textarea
          id="message-content-input"
          is="hermes-textarea-autoresize"
          focus-on-input="true"
          max-rows="12"
          name="content"
          placeholder="${this.getAttribute("placeholder") ?? "Message..."}"
          autofocus
          required
          bind="${(/** @type {unknown} */ el) => {
            this.#textarea = AssertInstance.once(el, HTMLTextAreaElement);
          }}"
        ></textarea>
        <div class="actions">
          <button
            onclick="${(/** @type {Event} */ e) => e.stopPropagation()}"
            id="submit-message"
            type="submit"
          >
            ↑
          </button>
        </div>
      </form>
    `);
  }
}

customElements.define("hermes-message-form", MessageForm);
