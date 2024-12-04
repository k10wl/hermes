import { Chats } from "./chats.mjs";
import { Link } from "./link.mjs";
import { MessageContentForm } from "./message-content-form.mjs";
import { Messages } from "./messages.mjs";
import { PaginatedList } from "./paginated-list.mjs";
import { TextAreaAutoresize } from "./textarea-autoresize.mjs";

/** @param {string} name  */
function withPrefix(name) {
  return "hermes-" + name;
}

/** @type {Parameters<typeof customElements.define>[] } */
const elements = [
  ["paginated-list", PaginatedList],
  ["link", Link, { extends: "a" }],
  ["chats", Chats],
  ["messages", Messages],
  ["textarea-autoresize", TextAreaAutoresize, { extends: "textarea" }],
  ["message-content-form", MessageContentForm, { extends: "form" }],
];

export function initCustomElements() {
  for (const [name, instance, options] of elements) {
    customElements.define(withPrefix(name), instance, options);
  }
}
