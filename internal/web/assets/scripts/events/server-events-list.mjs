import { Chat } from "/assets/scripts/models.mjs";
import {
  ValidateBoolean,
  ValidateNumber,
  ValidateObject,
  ValidateString,
} from "/assets/scripts/utils/validate.mjs";

export class ServerEvent {
  static #eventValidation = new ValidateObject({ type: ValidateString });

  /** @type {string} */
  type;

  /** @param { { type: string } } event - The type of the event. */
  constructor(event) {
    this.type = event.type;
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
