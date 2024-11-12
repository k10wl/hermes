import { ValidateNumber } from "/assets/scripts/utils/validate.mjs";

import { ServerEvent } from "./server-events-list.mjs";

export class RequestReadChatEvent extends ServerEvent {
  /** @param {number} chatId  */
  constructor(chatId) {
    super({ type: "request-read-chat", payload: ValidateNumber.parse(chatId) });
  }
}
