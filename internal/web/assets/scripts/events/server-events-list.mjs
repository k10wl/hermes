import { Chat, Message } from "/assets/scripts/models.mjs";
import {
  ValidateArray,
  ValidateBoolean,
  ValidateNumber,
  ValidateObject,
  ValidateOptional,
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
      chat: new ValidateObject({
        id: ValidateNumber,
        name: ValidateString,
      }),
      message: Message.validator,
      redirect: new ValidateOptional(ValidateBoolean),
    }),
  });

  /** @param { ReturnType<ChatCreatedEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = {
      ...data.payload,
      chat: new Chat(data.payload.chat.id, data.payload.chat.name),
      message: new Message(data.payload.message),
    };
  }

  /** @param {unknown} data */
  static parse(data) {
    return new ChatCreatedEvent(
      ChatCreatedEvent.validate(JSON.parse(ValidateString.parse(data))),
    );
  }

  /** @param {unknown} data */
  static validate(data) {
    return ChatCreatedEvent.#eventValidation.parse(data);
  }
}

export class ReadChatEvent extends ServerEvent {
  static #eventValidation = new ValidateObject({
    type: ValidateString,
    payload: new ValidateObject({
      messages: new ValidateArray(Message.validator),
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

export class MessageCreatedEvent extends ServerEvent {
  static #eventValidation = new ValidateObject({
    type: ValidateString,
    payload: new ValidateObject({
      chat_id: ValidateNumber,
      message: Message.validator,
    }),
  });

  /**
   * @param {{
   *   type: string,
   *   payload: {
   *     chat_id: number
   *     message: Message
   *   }
   * }} event
   */
  constructor(event) {
    super(event);
    this.payload = event.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    const parsed = MessageCreatedEvent.#eventValidation.parse(
      JSON.parse(ValidateString.parse(data)),
    );
    return new MessageCreatedEvent({
      ...parsed,
      payload: {
        ...parsed.payload,
        message: new Message(parsed.payload.message),
      },
    });
  }
}
