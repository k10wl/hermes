import {
  AssertInstance,
  AssertNumber,
  AssertString,
} from "/assets/scripts/lib/assert.mjs";
import { CreateCompletionMessageEvent } from "/assets/scripts/lib/events/client-events-list.mjs";
import { ServerEvents } from "/assets/scripts/lib/events/server-events.mjs";
import {
  ChatCreatedEvent,
  ServerErrorEvent,
} from "/assets/scripts/lib/events/server-events-list.mjs";
import { LocationControll } from "/assets/scripts/lib/location-control.mjs";

import { TextAreaAutoresize } from "./textarea-autoresize.mjs";

export class MessageContentForm extends HTMLFormElement {
  /** @type (() => void)[] */
  #messageContentFormCleanup = [];

  constructor() {
    super();
    this.detectKeyboardSubmit = this.detectKeyboardSubmit.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
    this.addEventListener("submit", this.handleSubmit);
    this.#messageContentFormCleanup.push(() => {
      this.removeEventListener("submit", this.handleSubmit);
    });
  }

  connectedCallback() {
    window.addEventListener("keydown", this.detectKeyboardSubmit);
    this.#messageContentFormCleanup.push(
      () => {
        window.removeEventListener("keydown", this.detectKeyboardSubmit);
      },
      // TODO maybe persist messages somewhere?
      LocationControll.attach({ notify: () => this.reset() }),
    );
  }

  reset() {
    super.reset();
    this.querySelectorAll("textarea").forEach((el) => {
      try {
        AssertInstance.once(el, TextAreaAutoresize).autoresize();
      } catch {
        // just don't explode
      }
    });
  }

  disconnectedCallback() {
    this.#messageContentFormCleanup.forEach((cb) => cb());
  }

  /** @param {KeyboardEvent} e */
  detectKeyboardSubmit(e) {
    if (e.key !== "Enter" || e.shiftKey || e.metaKey || e.ctrlKey) {
      return;
    }
    e.preventDefault();
    e.stopPropagation();
    this.requestSubmit();
  }

  /** @param {SubmitEvent} e  */
  handleSubmit(e) {
    e.preventDefault();
    /** @type {string | number | undefined} */
    let chat_id = LocationControll.pathname.split("/").at(-1);
    chat_id = AssertNumber.check(chat_id ? +chat_id : -1);
    const content = AssertString.check(new FormData(this).get("content"));
    const message = new CreateCompletionMessageEvent({
      chat_id,
      content: content,
      parameters: {
        model: "gpt-4o-mini",
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
        this.reset();
      },
    );
  }
}
