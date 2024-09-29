export class Chat {
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
  /** @param {{
   *   id: number
   *   chatId: number
   *   role: "user" | "assistant" | "system"
   *   content: string
   * }} message */
  constructor(message) {
    this.id = message.id;
    this.chatId = message.chatId;
    this.role = message.role;
    this.content = message.content;
  }
}
