import {
  ValidateNumber,
  ValidateObject,
  ValidateString,
} from "/assets/scripts/lib/validate.mjs";

export class Chat {
  static validator = new ValidateObject({
    id: ValidateNumber,
    name: ValidateString,
  });

  /**
   * @param {number} id
   * @param {string} name
   */
  constructor(id, name) {
    this.id = id;
    this.name = name;
  }
}

export class Message {
  static validator = new ValidateObject({
    id: ValidateNumber,
    chat_id: ValidateNumber,
    content: ValidateString,
    role: ValidateString,
  });

  /** @param {{
   *   id: number
   *   chat_id: number
   *   role: "user" | "assistant" | "system" | string
   *   content: string
   * }} message */
  constructor(message) {
    this.id = message.id;
    this.chat_id = message.chat_id;
    this.role = message.role;
    this.content = message.content;
  }
}
