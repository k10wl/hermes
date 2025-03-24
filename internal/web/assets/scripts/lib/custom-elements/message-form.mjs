import { AssertInstance, AssertNumber, AssertString } from "../assert.mjs";
import { CreateCompletionMessageEvent } from "../events/client-events-list.mjs";
import { ServerEvents } from "../events/server-events.mjs";
import {
  ChatCreatedEvent,
  ServerErrorEvent,
} from "../events/server-events-list.mjs";
import { FocusOnKeydown } from "../focus-on-keydown.mjs";
import { Bind, html } from "../libdim.mjs";
import { LocationControll } from "../location-control.mjs";
import { ResizableTextInput } from "./content-editable-plain-text.mjs";
import { AlertDialog } from "./dialog.mjs";

export class MessageForm extends HTMLElement {
  /** @type {Provider} */
  #provider;
  #focusOnInput = new FocusOnKeydown();
  #content = new Bind((el) => AssertInstance.once(el, ResizableTextInput));
  #temperature = new Bind((el) => AssertInstance.once(el, HTMLInputElement));
  #max_tokens = new Bind((el) => AssertInstance.once(el, HTMLInputElement));
  #model = new Bind((el) => AssertInstance.once(el, HTMLInputElement));

  /** @param {Provider} provider */
  constructor(provider) {
    super();
    this.#provider = provider;
  }

  #submit = async (/** @type {Event} */ e) => {
    e.preventDefault();
    const content = AssertString.check(this.#content.current.value);
    if (content.trim() === "") {
      return;
    }
    try {
      const params = {
        model: this.#model.current.value,
        max_tokens: +this.#max_tokens.current.value || undefined,
        temperature: +this.#temperature.current.value || undefined,
      };
      await this.#provider.submit({
        content,
        params,
      });
      this.#content.current.value = "";
    } catch (error) {
      console.error("caught error, not reseting content", error);
    }
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
        onsubmit="${this.#submit}"
        onkeydown="${(e) => {
          const event = AssertInstance.once(e, KeyboardEvent);
          const form = AssertInstance.once(e.currentTarget, HTMLFormElement);
          if (event.key === "Enter" && !event.shiftKey) {
            form.requestSubmit();
            event.preventDefault();
          }
        }}"
      >
        <input
          id="temperature"
          type="range"
          min="0"
          value="0.8"
          max="1"
          step="0.01"
          bind="${this.#temperature}"
        />
        <input
          id="max_tokens"
          type="text"
          inputmode="numeric"
          pattern="[0-9]*"
          value="0"
          bind="${this.#max_tokens}"
        />
        <input id="model" value="openai/gpt-4o-mini" bind="${this.#model}" />

        <h-resizable-text-input
          id="content"
          placeholder="${this.getAttribute("placeholder") ?? "Message"}"
          bind="${this.#content}"
        ></h-resizable-text-input>

        <button type="submit">↑</button>
      </form>
    `);
    this.#focusOnInput.attach(this.#content.current.content);
  }

  disconnectedCallback() {
    this.#focusOnInput.detach();
  }
}

class BusinessLogic {
  /** @type {Provider['submit']} */
  static submit(message) {
    const { content } = message;
    const { resolve, reject, promise } =
      /** @type {PromiseWithResolvers<void>} */
      (Promise.withResolvers());
    const createChatCompletionMessageEvent = new CreateCompletionMessageEvent({
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
    ServerEvents.send(createChatCompletionMessageEvent);
    const off = ServerEvents.on(
      ["chat-created", "message-created", "server-error"],
      (event) => {
        if (event.id !== createChatCompletionMessageEvent.id) {
          resolve();
          return;
        }
        off();
        if (event instanceof ChatCreatedEvent) {
          LocationControll.navigate(`/chats/${event.payload.chat.id}`);
          resolve();
          return;
        }
        if (event instanceof ServerErrorEvent) {
          AlertDialog.instance.alert({
            title: "Failed to send message",
            description: event.payload,
          });
          reject(event.payload);
          return;
        }
        resolve();
      },
    );
    return promise;
  }
}

/**
 * @typedef Params
 * @property {string} model
 * @property {number | undefined} [temperature]
 * @property {number | undefined} [max_tokens]
 */

/**
 * @typedef Provider
 * @property {(param: { content: string, params: Params }) => Promise<void>} submit
 */

/**
 * @param {Provider} provider
 * @returns {typeof MessageForm}
 */
export function creator(provider) {
  return class extends MessageForm {
    constructor() {
      super(provider);
    }
  };
}

customElements.define("hermes-message-form", creator(BusinessLogic));
