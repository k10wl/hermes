import "./button.mjs";
import "./chats.mjs";
import "./context-menu.mjs";
import "./dialog.mjs";
import "./form.mjs";
import "./header.mjs";
import "./key.mjs";
import "./message-form.mjs";
import "./paginated-list.mjs";
import "./scenes/chats-list.mjs";
import "./scenes/existing-chat.mjs";
import "./scenes/new-chat.mjs";
import "./scenes/templates-list.mjs";
import "./scenes/view-template.mjs";

// XXX this was such a shitty idea tbh
import { ControlPanel } from "./control-panel.mjs";
import { Link } from "./link.mjs";
import { Messages } from "./messages.mjs";
import { TextAreaAutoresize } from "./textarea-autoresize.mjs";

/** @param {string} name  */
function withPrefix(name) {
  return "hermes-" + name;
}

/** @type {Parameters<typeof customElements.define>[] } */
const elements = [
  ["link", Link, { extends: "a" }],
  ["messages", Messages],
  ["textarea-autoresize", TextAreaAutoresize, { extends: "textarea" }],
  ["control-panel", ControlPanel],
];

for (const [name, instance, options] of elements) {
  customElements.define(withPrefix(name), instance, options);
}
