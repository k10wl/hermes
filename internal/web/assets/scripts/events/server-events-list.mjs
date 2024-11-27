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
  static #eventValidation = new ValidateObject({
    id: ValidateString,
    type: ValidateString,
  });

  /** @type {string} */
  id;

  /** @type {string} */
  type;

  /** @type {unknown} */
  payload;

  /** @param { { id: string, type: string, payload?: unknown } } data */
  constructor(data) {
    this.id = data.id;
    this.type = data.type;
    this.payload = data.payload;
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
    super({ id: crypto.randomUUID(), type: "connection-status-change" });
    this.payload = { connected: ValidateBoolean.parse(connected) };
  }
}

export class ChatCreatedEvent extends ServerEvent {
  static #eventValidation = new ValidateObject({
    id: ValidateString,
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
    id: ValidateString,
    type: ValidateString,
    payload: new ValidateObject({
      messages: new ValidateArray(Message.validator),
    }),
  });

  /** @param { ReturnType<ReadChatEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    const parsed = ReadChatEvent.validate(
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

  /** @param {unknown} data */
  static validate(data) {
    return ReadChatEvent.#eventValidation.parse(data);
  }
}

export class ServerErrorEvent extends ServerEvent {
  static #eventValidation = new ValidateObject({
    id: ValidateString,
    type: ValidateString,
    payload: ValidateString,
  });

  /** @param { ReturnType<ServerErrorEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    return new ServerErrorEvent(
      ServerErrorEvent.validate(JSON.parse(ValidateString.parse(data))),
    );
  }

  /** @param {unknown} data */
  static validate(data) {
    return ServerErrorEvent.#eventValidation.parse(data);
  }
}

export class MessageCreatedEvent extends ServerEvent {
  static #eventValidation = new ValidateObject({
    id: ValidateString,
    type: ValidateString,
    payload: new ValidateObject({
      chat_id: ValidateNumber,
      message: Message.validator,
    }),
  });

  /** @param { ReturnType<MessageCreatedEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    const parsed = MessageCreatedEvent.validate(
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

  /** @param {unknown} data */
  static validate(data) {
    return MessageCreatedEvent.#eventValidation.parse(data);
  }
}
