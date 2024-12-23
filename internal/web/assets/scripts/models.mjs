import {
  AssertNumber,
  AssertObject,
  AssertString,
} from "/assets/scripts/lib/assert.mjs";

export class Chat {
  static validator = new AssertObject({
    id: AssertNumber,
    name: AssertString,
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
  static validator = new AssertObject({
    id: AssertNumber,
    chat_id: AssertNumber,
    content: AssertString,
    role: AssertString,
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

export class Template {
  static validator = new AssertObject({
    id: AssertNumber,
    name: AssertString,
    content: AssertString,
  });

  /** @param {ReturnType<Template.validator['check']>} template  */
  constructor(template) {
    this.id = template.id;
    this.name = template.name;
    this.content = template.content;
  }
}
