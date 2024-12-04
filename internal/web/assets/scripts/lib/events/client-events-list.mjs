import {
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
