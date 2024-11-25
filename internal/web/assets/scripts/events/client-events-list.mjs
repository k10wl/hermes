import {
  ValidateNumber,
  ValidateObject,
  ValidateOptional,
  ValidateString,
} from "/assets/scripts/utils/validate.mjs";

import { ServerEvent } from "./server-events-list.mjs";

class ClientEvent extends ServerEvent {}

export class RequestReadChatEvent extends ClientEvent {
  /** @param {number} chatId  */
  constructor(chatId) {
    super({ type: "request-read-chat", payload: ValidateNumber.parse(chatId) });
  }
}

export class CreateCompletionMessageEvent extends ClientEvent {
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
