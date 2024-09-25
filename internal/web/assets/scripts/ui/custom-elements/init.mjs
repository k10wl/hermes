import { PaginatedList } from "./paginated-list.mjs";

/** @param {string} name  */
function withPrefix(name) {
  return "hermes-" + name;
}

/** @type {{name: string, customElement: CustomElementConstructor}[]} */
const elements = [{ name: "paginated-list", customElement: PaginatedList }];

export function initCustomElements() {
  for (let i = 0; i < elements.length; i++) {
    const el = elements[i];
    if (!el) {
      throw new Error("failed to get element");
    }
    customElements.define(withPrefix(el.name), el.customElement);
  }
}
