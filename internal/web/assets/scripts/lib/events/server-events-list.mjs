import {
  AssertArray,
  AssertBoolean,
  AssertInstance,
  AssertNumber,
  AssertObject,
  AssertOptional,
  AssertString,
} from "/assets/scripts/lib/assert.mjs";
import { Chat, Message, Template } from "/assets/scripts/models.mjs";

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

export class ReadTemplatesEvent extends ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
    payload: new AssertObject({
      templates: new AssertArray(new AssertInstance(Template)),
    }),
  });

  static canonicalType = /** @type {const} */ ("read-templates");

  /** @param { ReturnType<ReadTemplatesEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    /** @type {ReadTemplatesEvent} or at least it should be */
    const parsed = JSON.parse(AssertString.check(data));
    parsed.payload.templates = parsed.payload.templates.map(
      (template) => new Template(template),
    );
    return new ReadTemplatesEvent(ReadTemplatesEvent.validate(parsed));
  }

  /** @param {unknown} data */
  static validate(data) {
    return ReadTemplatesEvent.#eventValidation.check(data);
  }
}

export class ReadTemplateEvent extends ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
    payload: new AssertObject({
      template: new AssertInstance(Template),
    }),
  });

  static canonicalType = /** @type {const} */ ("read-template");

  /** @param { ReturnType<ReadTemplateEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    /** @type {ReadTemplateEvent} or at least it should be */
    const parsed = JSON.parse(AssertString.check(data));
    parsed.payload.template = new Template(parsed.payload.template);
    return new ReadTemplateEvent(ReadTemplateEvent.validate(parsed));
  }

  /** @param {unknown} data */
  static validate(data) {
    return ReadTemplateEvent.#eventValidation.check(data);
  }
}

export class TemplateChangedEvent extends ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
    payload: new AssertObject({
      template: new AssertInstance(Template),
    }),
  });

  static canonicalType = /** @type {const} */ ("template-changed");

  /** @param { ReturnType<TemplateChangedEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    /** @type {TemplateChangedEvent} or at least it should be */
    const parsed = JSON.parse(AssertString.check(data));
    parsed.payload.template = new Template(parsed.payload.template);
    return new TemplateChangedEvent(TemplateChangedEvent.validate(parsed));
  }

  /** @param {unknown} data */
  static validate(data) {
    return TemplateChangedEvent.#eventValidation.check(data);
  }
}

export class TemplateCreatedEvent extends ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
    payload: new AssertObject({
      template: new AssertInstance(Template),
    }),
  });

  static canonicalType = /** @type {const} */ ("template-created");

  /** @param { ReturnType<TemplateCreatedEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    /** @type {TemplateCreatedEvent} or at least it should be */
    const parsed = JSON.parse(AssertString.check(data));
    parsed.payload.template = new Template(parsed.payload.template);
    return new TemplateCreatedEvent(TemplateCreatedEvent.validate(parsed));
  }

  /** @param {unknown} data */
  static validate(data) {
    return TemplateCreatedEvent.#eventValidation.check(data);
  }
}

export class TemplateDeletedEvent extends ServerEvent {
  static #eventValidation = new AssertObject({
    id: AssertString,
    type: AssertString,
    payload: new AssertObject({
      name: AssertString,
    }),
  });

  static canonicalType = /** @type {const} */ ("template-deleted");

  /** @param { ReturnType<TemplateDeletedEvent.validate> } data */
  constructor(data) {
    super(data);
    this.payload = data.payload;
  }

  /** @param {unknown} data */
  static parse(data) {
    /** @type {TemplateDeletedEvent} or at least it should be */
    const parsed = JSON.parse(AssertString.check(data));
    return new TemplateDeletedEvent(TemplateDeletedEvent.validate(parsed));
  }

  /** @param {unknown} data */
  static validate(data) {
    return TemplateDeletedEvent.#eventValidation.check(data);
  }
}
