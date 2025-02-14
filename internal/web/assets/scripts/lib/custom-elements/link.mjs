import { AssertInstance } from "/assets/scripts/lib/assert.mjs";
import { LocationControll } from "/assets/scripts/lib/location-control.mjs";

export class Link extends HTMLAnchorElement {
  constructor() {
    super();
  }
  /** @param {MouseEvent} event  */
  #onclick(event) {
    if (event.ctrlKey || event.metaKey || event.altKey || event.shiftKey) {
      return;
    }
    LocationControll.navigate(
      AssertInstance.once(event.currentTarget, HTMLAnchorElement).href,
    );
    event.preventDefault();
  }
  connectedCallback() {
    this.addEventListener("click", this.#onclick);
  }
  disconnectedCallback() {
    this.removeEventListener("click", this.#onclick);
  }
}
