import { AssertInstance } from "/assets/scripts/lib/assert.mjs";
import { RequestReadChatEvent } from "/assets/scripts/lib/events/client-events-list.mjs";
import { ServerEvents } from "/assets/scripts/lib/events/server-events.mjs";
import { html } from "/assets/scripts/lib/html-v2.mjs";
import { LocationControll } from "/assets/scripts/lib/location-control.mjs";

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

        .message {
          margin: 0.75rem 0;
        }

        .content {
          user-select: text;
          white-space: pre-wrap;
          word-break: break-word;
        }

        .role-user {
          --_border-radius: 0.66rem;
          color: var(--_text);
          background: rgb(from var(--_primary) r g b / 0.05);
          border: 1px solid rgb(from var(--_primary) r g b / 0.25);
          margin-left: auto;
          width: fit-content;
          max-width: 80%;
          border-radius: var(--_border-radius) var(--_border-radius) 0
            var(--_border-radius);
          padding: 0.33rem;
          padding-bottom: 0.2rem;
        }

        .actions {
          display: flex;
          gap: 0.5rem;
          margin-top: 0.15rem;
        }

        .actions button {
          background: var(--bg-2);
          border: 1px solid rgb(from var(--text-0) r g b / 0.25);
          font-size: 1rem;
          --_size: 2ch;
          width: var(--_size);
          height: var(--_size);
          border-radius: calc(var(--_size) / 4);
          display: flex;
          align-items: center;
          justify-content: center;
          padding: 0.25rem;
          color: rgb(from var(--text-0) r g b / 0.5);
          background: transparent;
          border: 1px solid transparent;
          &:hover {
            color: var(--_text);
            border-color: rgb(from var(--text-0) r g b / 0.25);
          }
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
        console.log("read-chat", data);
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

  /** @param {import("/assets/scripts/models.mjs").Message} message */
  #messageToHtml(message) {
    const escaped = message.content
      .replace(/&/g, "&amp;")
      .replace(/</g, "&lt;")
      .replace(/>/g, "&gt;");

    if (message.role === "user") {
      return html`
        <div class="message content role-${message.role}">${escaped}</div>
      `;
    }

    return html`
      <div class="message">
        <div class="content">${escaped}</div>
        <div class="actions">
          <button
            onclick="${(/** @type {Event} */ event) => {
              const button = AssertInstance.once(
                event.target,
                HTMLButtonElement,
              );
              const prev = button.textContent;
              button.textContent = "✓";
              setTimeout(() => {
                button.textContent = prev;
              }, 2000);
              window.navigator.clipboard.writeText(message.content);
            }}"
          >
            ⎘
          </button>
        </div>
      </div>
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
