import { Bind, html, Signal } from "/assets/scripts/lib/libdim.mjs";

import { AssertInstance, AssertNumber } from "../assert.mjs";
import { CreateCompletionMessageEvent } from "../events/client-events-list.mjs";
import { ServerEvents } from "../events/server-events.mjs";
import {
  ChatCreatedEvent,
  ServerErrorEvent,
} from "../events/server-events-list.mjs";
import { FocusOnKeydown } from "../focus-on-keydown.mjs";
import { LocationControll } from "../location-control.mjs";
import { ResizableTextInput } from "./content-editable-plain-text.mjs";
import { AlertDialog } from "./dialog.mjs";

export class MessageForm extends HTMLElement {
  #form = new Bind((el) => AssertInstance.once(el, HTMLFormElement));
  #content = new Bind((el) => AssertInstance.once(el, ResizableTextInput));
  #focusOnInput = new FocusOnKeydown();
  #empty = new Signal(true);

  constructor() {
    super();
  }

  #submit = (/** @type {Event} */ e) => {
    e.preventDefault();
    const content = this.#content.current.value;
    if (content.trim() === "") {
      return;
    }
    const message = new CreateCompletionMessageEvent({
      chat_id: AssertNumber.check(
        LocationControll.chatId ? +LocationControll.chatId : -1,
      ),
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
          return;
        }
        if (event instanceof ServerErrorEvent) {
          AlertDialog.instance.alert({
            title: "Failed to send message",
            description: event.payload,
          });
          return;
        }
        this.#content.current.value = "";
        this.#empty.value = true;
      },
    );
  };

  connectedCallback() {
    this.attachShadow({ mode: "open" }).append(html`
      <style>
        :host {
          --bg: var(--bg-2);
          --text: var(--text-0);
          --radius: 1rem;
        }

        form {
          display: flex;
          justify-content: center;
          align-items: flex-end;
          gap: 0.5rem;
        }

        form:has([data-empty="true"]) button[type="submit"] {
          background: var(--bg);
          color: rgb(from var(--text) r g b / 0.25);
          cursor: auto;
        }

        button {
          --_size: calc(2rem + 2px);
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

        h-resizable-text-input {
          --padding: 0.5rem 1rem;
          --border-color: var(--bg);
          --border: 1px solid var(--border-color);
          width: 100%;

          &:focus-within {
            --border-color: var(--primary);
          }

          &::part(wrapper) {
            border: var(--border);
            border-radius: 1.25rem;
            overflow: hidden;
            color: var(--text);
          }

          &::part(content) {
            padding: var(--padding);
            max-height: max(50vh);
          }

          &::part(placeholder) {
            padding: var(--padding);
          }
        }
      </style>

      <form
        bind="${this.#form}"
        onsubmit="${this.#submit}"
        onkeydown="${(e) => {
          const event = AssertInstance.once(e, KeyboardEvent);
          if (event.key === "Enter" && !event.shiftKey) {
            this.#form.current.requestSubmit();
            event.preventDefault();
          }
        }}"
      >
        <h-resizable-text-input
          id="content"
          placeholder="${this.getAttribute("placeholder") ?? "Message"}"
          bind="${this.#content}"
        ></h-resizable-text-input>

        <button id="submit-message" type="submit">↑</button>
      </form>
    `);
    this.#focusOnInput.attach(this.#content.current.content);
  }

  disconnectedCallback() {
    this.#focusOnInput.detach();
  }
}

customElements.define("hermes-message-form", MessageForm);
