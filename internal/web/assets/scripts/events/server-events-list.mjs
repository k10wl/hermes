import { Chat, Message } from "/assets/scripts/models.mjs";
import {
  ValidateArray,
  ValidateBoolean,
  ValidateNumber,
  ValidateObject,
  ValidateString,
} from "/assets/scripts/utils/validate.mjs";

export class ServerEvent {
  static #eventValidation = new ValidateObject({ type: ValidateString });

  /** @type {string} */
  type;

  /** @type {unknown} */
  payload;

  /** @param { { type: string, payload?: unknown } } event - The type of the event. */
  constructor(event) {
    this.type = event.type;
    this.payload = event.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    return new ServerEvent(
      ServerEvent.#eventValidation.parse(
        JSON.parse(ValidateString.parse(data)),
      ),
    );
  }
}

export class ConnectionStatusChangeEvent extends ServerEvent {
  /**
   * @param {  boolean  } connected
   */
  constructor(connected) {
    super({ type: "connection-status-change" });
    this.payload = { connected: ValidateBoolean.parse(connected) };
  }
}

export class ChatCreatedEvent extends ServerEvent {
  static #eventValidation = new ValidateObject({
    type: ValidateString,
    payload: new ValidateObject({
      id: ValidateNumber,
      name: ValidateString,
    }),
  });

  /**
   * @param { { type: string, payload: { id: number, name: string } } } event
   */
  constructor(event) {
    super(event);
    this.payload = new Chat(event.payload.id, event.payload.name);
  }

  /** @param {unknown} data */
  static parse(data) {
    return new ChatCreatedEvent(
      ChatCreatedEvent.#eventValidation.parse(
        JSON.parse(ValidateString.parse(data)),
      ),
    );
  }
}

export class ReadChatEvent extends ServerEvent {
  static #eventValidation = new ValidateObject({
    type: ValidateString,
    payload: new ValidateObject({
      messages: new ValidateArray(
        new ValidateObject({
          id: ValidateNumber,
          chat_id: ValidateNumber,
          content: ValidateString,
          role: ValidateString,
        }),
      ),
    }),
  });

  /**
   * @param {{
   *   type: string,
   *   payload: {
   *     messages: Message[]
   *   }
   * }} event
   */
  constructor(event) {
    super(event);
    this.payload = event.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    const parsed = ReadChatEvent.#eventValidation.parse(
      JSON.parse(ValidateString.parse(data)),
    );
    return new ReadChatEvent({
      ...parsed,
      payload: {
        ...parsed.payload,
        messages: parsed.payload.messages.map(
          (message) => new Message(message),
        ),
      },
    });
  }
}

export class ServerErrorEvent extends ServerEvent {
  static #eventValidation = new ValidateObject({
    type: ValidateString,
    payload: ValidateString,
  });

  /**
   * @param {{
   *   type: string,
   *   payload: string
   * }} event
   */
  constructor(event) {
    super(event);
    this.payload = event.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    return new ServerErrorEvent(
      ServerErrorEvent.#eventValidation.parse(
        JSON.parse(ValidateString.parse(data)),
      ),
    );
  }
}
