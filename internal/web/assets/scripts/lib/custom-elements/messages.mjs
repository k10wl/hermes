import { AssertInstance } from "/assets/scripts/lib/assert.mjs";
import { RequestReadChatEvent } from "/assets/scripts/lib/events/client-events-list.mjs";
import { ServerEvents } from "/assets/scripts/lib/events/server-events.mjs";
import { html } from "/assets/scripts/lib/html-v2.mjs";
import { LocationControll } from "/assets/scripts/lib/location-control.mjs";
import { Message } from "/assets/scripts/models.mjs";

customElements.define(
  "h-chat-message",
  class extends HTMLElement {
    messageContent = "";

    constructor() {
      super();
    }

    connectedCallback() {
      this.attachShadow({ mode: "open" }).append(html`
        <style>
          #wrapper {
            display: grid;
            padding: 0.5rem 0;
          }
          #content {
            user-select: text;
          }
          slot {
            white-space: pre-wrap;
            word-break: break-word;
            color: var(--_text);
          }
          :host(.role-user) {
            #content {
              --_border-radius: 0.66rem;
              background: rgb(from var(--_primary) r g b / 0.05);
              border: 1px solid rgb(from var(--_primary) r g b / 0.25);
              justify-self: end;
              width: fit-content;
              max-width: 80%;
              border-radius: var(--_border-radius) var(--_border-radius)
                calc(var(--_border-radius) / 3) var(--_border-radius);
              padding: 0.33rem;
              padding-bottom: 0.2rem;
            }
            h-message-actions {
              justify-self: end;
            }
          }
          h-message-actions {
            --_duration: 100ms;
            transition: all var(--_duration);
            transition-delay: var(--_duration);
            margin-top: 0.2rem;
            opacity: 0;
          }
          #wrapper:hover h-message-actions,
          #wrapper:focus-within h-message-actions {
            opacity: 1;
            transition-delay: 0ms;
          }
        </style>
        <div id="wrapper">
          <div id="content">
            <slot></slot>
          </div>
          <h-message-actions
            oncopycontent="${() => {
              navigator.clipboard.writeText(this.messageContent);
            }}"
            part="actions"
          ></h-message-actions>
        </div>
      `);
    }
  },
);

customElements.define(
  "h-message-actions",
  class extends HTMLElement {
    constructor() {
      super();
      this.attachShadow({ mode: "closed" }).append(html`
        <style>
          :host {
            --_text: var(--text-0);
            --_primary: var(--primary);
          }

          button {
            --_size: 2ch;
            display: flex;
            align-items: center;
            justify-content: center;
            width: var(--_size);
            height: var(--_size);
            padding: cacl(var(--_size) / 2);
            color: rgb(from var(--_text) r g b / 0.5);
            background: transparent;
            font-size: 1rem;
            border: 1px solid transparent;
            border-radius: calc(var(--_size) / 4);
            &:hover {
              color: var(--_text);
              border-color: rgb(from var(--_text) r g b / 0.25);
            }
          }
        </style>

        <button onclick="${this.#handleCopy}">⎘</button>
      `);
    }

    #handleCopy = (/** @type {Event} */ event) => {
      this.dispatchEvent(
        new CustomEvent("copycontent", {
          bubbles: true,
          composed: true,
        }),
      );
      const button = AssertInstance.once(event.target, HTMLButtonElement);
      const prev = button.textContent;
      button.textContent = "✓";
      setTimeout(() => {
        button.textContent = prev;
      }, 2000);
    };
  },
);

export class Messages extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanupOnDisconnect = [];

  /** @type {HTMLElement | undefined} */
  #messages;

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.shadow.append(html`
      <style>
        :host {
          --_text: var(--text-0);
          --_primary: var(--primary);
        }

        section {
          padding: 0.75rem;
        }

        h-chat-message.role-assistant:last-child::part(actions) {
          opacity: 1;
        }
      </style>

      <section
        bind="${(/** @type {unknown} */ el) => {
          this.#messages = AssertInstance.once(el, HTMLElement);
        }}"
      ></section>
    `);
  }

  connectedCallback() {
    const messagesContainer = AssertInstance.once(this.#messages, HTMLElement);
    const routeObserver = new RouteObserver();
    routeObserver.notify();
    this.#cleanupOnDisconnect.push(
      LocationControll.attach(routeObserver),
      ServerEvents.on("message-created", (data) => {
        if (data.payload.chat_id === LocationControll.chatId) {
          messagesContainer.append(this.#messageToHtml(data.payload.message));
        }
      }),
      ServerEvents.on("read-chat", (data) => {
        messagesContainer.replaceChildren(
          ...data.payload.messages.map((message) =>
            this.#messageToHtml(message),
          ),
        );
      }),
    );
  }

  disconnectedCallback() {
    for (const cleanup of this.#cleanupOnDisconnect) {
      cleanup();
    }
  }

  /** @param {Message} message */
  #messageToHtml(message) {
    const { role, content } = Message.validator.check(message);
    return html`
      <h-chat-message
        class="role-${role}"
        bind="${(/** @type {unknown} */ el) => {
          const element = AssertInstance.once(el, customElements.get("h-chat-message"));
          element.messageContent = content;
        }}"
        >${content
          .replace(/&/g, "&amp;")
          .replace(/</g, "&lt;")
          .replace(/>/g, "&gt;")}</h-chat-message
      >
    `;
  }
}

class RouteObserver {
  /** @type {number | null} */
  #last = null;

  constructor() {}

  notify() {
    const chatId = LocationControll.chatId;
    const prev = this.#last;
    this.#last = chatId;
    if (prev === chatId || !chatId) {
      return;
    }
    ServerEvents.send(new RequestReadChatEvent(chatId));
  }
}
