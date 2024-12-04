import {
  ValidateNumber,
  ValidateObject,
  ValidateOptional,
  ValidateString,
} from "/assets/scripts/lib/validate.mjs";

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
      payload: ValidateNumber.parse(chatId),
    });
  }
}

export class CreateCompletionMessageEvent extends ClientEvent {
  /** @type {"create-completion"} */
  static canonicalType = "create-completion";

  static #eventValidation = new ValidateObject({
    chat_id: ValidateNumber,
    content: ValidateString,
    parameters: new ValidateObject({
      model: ValidateString,
      max_tokens: new ValidateOptional(ValidateNumber),
      temperature: new ValidateOptional(ValidateNumber),
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
    return CreateCompletionMessageEvent.#eventValidation.parse(data);
  }
}
