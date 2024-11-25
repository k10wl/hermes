import { CreateCompletionMessageEvent } from "/assets/scripts/events/client-events-list.mjs";
import { ServerEvents } from "/assets/scripts/events/server-events.mjs";
import { LocationControll } from "/assets/scripts/lib/navigation/location.mjs";
import {
  ValidateNumber,
  ValidateString,
} from "/assets/scripts/utils/validate.mjs";

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
    chat_id = ValidateNumber.parse(chat_id ? +chat_id : -1);
    const content = ValidateString.parse(new FormData(this).get("content"));
    ServerEvents.send(
      new CreateCompletionMessageEvent({
        chat_id,
        content: content,
        parameters: {
          model: "gpt-4o-mini",
          max_tokens: undefined,
          temperature: undefined,
        },
      }),
    );
    let eventOff = () => {};
    let locationControllOff = () => {};
    if (chat_id === -1) {
      eventOff = ServerEvents.on("chat-created", (e) => {
        if (e.payload.message.content !== content) {
          return;
        }
        LocationControll.navigate(`/chats/${e.payload.chat.id}`);
        this.reset();
        eventOff();
        locationControllOff();
      });
    } else {
      eventOff = ServerEvents.on("message-created", (e) => {
        if (
          // XXX this is not a best comparison, would be good to have some
          // traveling id to indicate concreate messages instead of... this
          e.payload.message.content !== content ||
          e.payload.chat_id !== chat_id
        ) {
          return;
        }
        this.reset();
        eventOff();
        locationControllOff();
      });
    }

    let skip = true;
    locationControllOff = LocationControll.attach({
      notify: () => {
        if (skip) {
          skip = false;
          return;
        }
        eventOff();
        locationControllOff();
      },
    });
  }
}
