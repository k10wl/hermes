import {
  AssertArray,
  AssertBoolean,
  AssertNumber,
  AssertObject,
  AssertOptional,
  AssertString,
} from "/assets/scripts/lib/assert.mjs";
import { Chat, Message } from "/assets/scripts/models.mjs";

export class ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
  });

  /** @type {string} */
  id;

  /** @type {string} */
  type;

  /** @type {unknown} */
  payload;

  /** @type {string} */
  static canonicalType = "__meant_to_be_overriden__";

  /** @param { { id: string, type: string, payload?: unknown } } data */
  constructor(data) {
    this.id = data.id;
    this.type = data.type;
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    return new ServerEvent(
      ServerEvent.#eventValidation.check(JSON.parse(AssertString.check(data))),
    );
  }
}

export class ReloadEvent extends ServerEvent {
  static canonicalType = /** @type {const} */ ("reload");

  /** @param { ServerEvent } data */
  constructor(data) {
    super(data);
  }
}

export class ConnectionStatusChangeEvent extends ServerEvent {
  static canonicalType = /** @type {const} */ ("connection-status-change");

  /**
   * @param {  boolean  } connected
   */
  constructor(connected) {
    super({ id: crypto.randomUUID(), type: "connection-status-change" });
    this.payload = { connected: AssertBoolean.check(connected) };
  }
}

export class ChatCreatedEvent extends ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
    payload: new AssertObject({
      chat: Chat.validator,
      message: Message.validator,
      redirect: new AssertOptional(AssertBoolean),
    }),
  });

  static canonicalType = /** @type {const} */ ("chat-created");

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
      ChatCreatedEvent.validate(JSON.parse(AssertString.check(data))),
    );
  }

  /** @param {unknown} data */
  static validate(data) {
    return ChatCreatedEvent.#eventValidation.check(data);
  }
}

export class ReadChatEvent extends ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
    payload: new AssertObject({
      messages: new AssertArray(Message.validator),
    }),
  });

  static canonicalType = /** @type {const} */ ("read-chat");

  /** @param { ReturnType<ReadChatEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    const parsed = ReadChatEvent.validate(JSON.parse(AssertString.check(data)));
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
    return ReadChatEvent.#eventValidation.check(data);
  }
}

export class ServerErrorEvent extends ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
    payload: AssertString,
  });

  static canonicalType = /** @type {const} */ ("server-error");

  /** @param { ReturnType<ServerErrorEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    return new ServerErrorEvent(
      ServerErrorEvent.validate(JSON.parse(AssertString.check(data))),
    );
  }

  /** @param {unknown} data */
  static validate(data) {
    return ServerErrorEvent.#eventValidation.check(data);
  }
}

export class MessageCreatedEvent extends ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
    payload: new AssertObject({
      chat_id: AssertNumber,
      message: Message.validator,
    }),
  });

  static canonicalType = /** @type {const} */ ("message-created");

  /** @param { ReturnType<MessageCreatedEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    const parsed = MessageCreatedEvent.validate(
      JSON.parse(AssertString.check(data)),
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
    return MessageCreatedEvent.#eventValidation.check(data);
  }
}
