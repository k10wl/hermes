import { LocationControll } from "/assets/scripts/lib/navigation/location.mjs";
import { assertInstance } from "/assets/scripts/utils/assert-instance.mjs";

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
      assertInstance(event.currentTarget, HTMLAnchorElement).href,
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
