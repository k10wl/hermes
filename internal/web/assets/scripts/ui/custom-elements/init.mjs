import { Chats } from "./chats.mjs";
import { Link } from "./link.mjs";
import { Messages } from "./messages.mjs";
import { PaginatedList } from "./paginated-list.mjs";

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
];

export function initCustomElements() {
  for (const element of elements) {
    customElements.define(withPrefix(element[0]), element[1], element[2]);
  }
}
