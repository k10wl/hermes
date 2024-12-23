import "./form.mjs";
import "./header.mjs";
import "./message-form.mjs";
import "./scenes/chats-list.mjs";
import "./scenes/existing-chat.mjs";
import "./scenes/new-chat.mjs";
import "./scenes/templates-list.mjs";
import "./scenes/view-template.mjs";

// XXX this was such a shitty idea tbh
import { Chats } from "./chats.mjs";
import { ControlPanel } from "./control-panel.mjs";
import { Link } from "./link.mjs";
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
  ["control-panel", ControlPanel],
];

for (const [name, instance, options] of elements) {
  customElements.define(withPrefix(name), instance, options);
}
