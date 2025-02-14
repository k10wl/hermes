import {
  AssertBoolean,
  AssertNumber,
  AssertObject,
  AssertOptional,
  AssertString,
} from "/assets/scripts/lib/assert.mjs";

import { ServerEvent } from "./server-events-list.mjs";

class ClientEvent extends ServerEvent {
  /** @type {string} */
  id;

  /** @param {Partial<Omit<ServerEvent, "id">> & {type: string}} data */
  constructor(data) {
    const id = crypto.randomUUID();
    super({ ...data, id });
    this.id = id;
  }
}

export class RequestReadChatEvent extends ClientEvent {
  /** @type {"request-read-chat"} */
  static canonicalType = "request-read-chat";

  /** @param {number} chatId  */
  constructor(chatId) {
    super({
      type: "request-read-chat",
      payload: AssertNumber.check(chatId),
    });
  }
}

export class CreateCompletionMessageEvent extends ClientEvent {
  /** @type {"create-completion"} */
  static canonicalType = "create-completion";

  static #eventValidation = new AssertObject({
    chat_id: AssertNumber,
    content: AssertString,
    parameters: new AssertObject({
      model: AssertString,
      max_tokens: new AssertOptional(AssertNumber),
      temperature: new AssertOptional(AssertNumber),
    }),
  });

  /** @param {ReturnType<CreateCompletionMessageEvent['validatePayload']>} payload  */
  constructor(payload) {
    super({
      type: "create-completion",
    });
    this.payload = this.validatePayload(payload);
  }

  /** @param {unknown} data */
  validatePayload(data) {
    return CreateCompletionMessageEvent.#eventValidation.check(data);
  }
}

export class RequestReadTemplatesEvent extends ClientEvent {
  static canonicalType = /** @type {const} */ "request-read-templates";

  static #eventValidation = new AssertObject({
    name: AssertString,
    start_before_id: AssertNumber,
    limit: AssertNumber,
  });

  /** @param {ReturnType<RequestReadTemplatesEvent['validatePayload']>} payload  */
  constructor(payload) {
    super({
      type: RequestReadTemplatesEvent.canonicalType,
    });
    this.payload = this.validatePayload(payload);
  }

  /** @param {unknown} data */
  validatePayload(data) {
    return RequestReadTemplatesEvent.#eventValidation.check(data);
  }
}

export class RequestReadTemplateEvent extends ClientEvent {
  static canonicalType = /** @type {const} */ "request-read-template";

  static #eventValidation = new AssertObject({
    id: AssertNumber,
  });

  /** @param {ReturnType<RequestReadTemplateEvent['validatePayload']>} payload  */
  constructor(payload) {
    super({
      type: RequestReadTemplateEvent.canonicalType,
    });
    this.payload = this.validatePayload(payload);
  }

  /** @param {unknown} data */
  validatePayload(data) {
    return RequestReadTemplateEvent.#eventValidation.check(data);
  }
}

export class RequestEditTemplateEvent extends ClientEvent {
  static canonicalType = /** @type {const} */ "request-edit-template";

  static #eventValidation = new AssertObject({
    name: AssertString,
    content: AssertString,
    clone: new AssertOptional(AssertBoolean),
  });

  /** @param {ReturnType<RequestEditTemplateEvent['validatePayload']>} payload  */
  constructor(payload) {
    super({
      type: RequestEditTemplateEvent.canonicalType,
    });
    this.payload = this.validatePayload(payload);
  }

  /** @param {unknown} data */
  validatePayload(data) {
    return RequestEditTemplateEvent.#eventValidation.check(data);
  }
}

export class DeleteTemplateEvent extends ClientEvent {
  static canonicalType = /** @type {const} */ "delete-template";

  static #eventValidation = new AssertObject({
    name: AssertString,
  });

  /** @param {ReturnType<DeleteTemplateEvent['validatePayload']>} payload  */
  constructor(payload) {
    super({
      type: DeleteTemplateEvent.canonicalType,
    });
    this.payload = this.validatePayload(payload);
  }

  /** @param {unknown} data */
  validatePayload(data) {
    return DeleteTemplateEvent.#eventValidation.check(data);
  }
}
