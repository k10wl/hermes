import { AssertInstance } from "/assets/scripts/lib/assert.mjs";
import { debounce } from "/assets/scripts/lib/debounce.mjs";
import { Bind, html } from "/assets/scripts/lib/libdim.mjs";

export class ExistingChatScene extends HTMLElement {
  constructor() {
    super();
  }

  #scrollableWrapper = new Bind((el) =>
    AssertInstance.once(el, HTMLDivElement),
  );

  connectedCallback() {
    this.attachShadow({ mode: "open" }).append(html`
      <style>
        main {
          height: 100vh;
          max-height: 100vh;
          overflow: auto;
          display: grid;
          grid-template-rows: 1fr auto;
        }

        #scrollable-message-wrapper {
          max-height: 100%;
          overflow: auto;
          display: flex;
          flex-direction: column-reverse;

          border-bottom: 1px solid transparent;
          &.scrolled {
            border-image: linear-gradient(
                to right,
                transparent 0%,
                rgb(from var(--text-0) r g b / 0.25) 50%,
                transparent 100%
              )
              50% 0%;
          }
        }

        #messages-width-wrapper {
          max-width: var(--container);
          display: flex;
          justify-content: center;
          align-self: center;
          width: 100%;
          position: relative;
        }

        .scrolled #to-bottom-wrapper {
          visibility: visible;
        }
        #to-bottom-wrapper {
          visibility: hidden;
          margin: auto;
          width: 100%;
          max-width: var(--container);
          position: fixed;
          bottom: 4rem;
          right: 1rem;
          height: 0;
          display: flex;
          justify-content: flex-end;
          z-index: 1;
          #to-bottom {
            --size: 2rem;
            width: var(--size);
            height: var(--size);
            border-radius: var(--size);
            cursor: pointer;
            translate: 0 -100%;
            background: var(--bg-2);
            color: rgb(from var(--text-0) r g b / 0.25);
            &:hover {
              color: var(--primary);
            }
            display: grid;
            place-items: center;
            outline: none;
            border: none;
            position: relative;
          }
        }

        #messages-list {
          width: 100%;
        }

        .input-form-wrapper {
          display: flex;
          justify-content: center;
        }

        hermes-message-form {
          --_pad: 1rem;
          padding: calc(var(--_pad) / 4) var(--_pad) var(--_pad);
          max-width: var(--container);
          width: 100%;
        }
      </style>

      <main>
        <div
          id="scrollable-message-wrapper"
          bind="${this.#scrollableWrapper}"
          onscroll="${this.#onScroll}"
        >
          <div id="to-bottom-wrapper">
            <button
              type="button"
              aria-label="Scroll to latest message"
              id="to-bottom"
              onclick="${this.#handleScrollToBottom}"
            >
              ⇊
            </button>
          </div>
          <div id="messages-width-wrapper">
            <hermes-messages id="messages-list"></hermes-messages>
          </div>
        </div>

        <div class="input-form-wrapper">
          <hermes-message-form placeholder="How can I help you?">
          </hermes-message-form>
        </div>
      </main>
    `);
  }

  #onScroll = debounce(() => {
    const scrollable = this.#scrollableWrapper.current;
    // 0 is top because flex-direction: column-reverse flips scroll calculation
    if (scrollable.scrollTop >= 0) {
      scrollable.classList.remove("scrolled");
    } else {
      scrollable.classList.add("scrolled");
    }
  }, 1000 / 60);

  #handleScrollToBottom = () => {
    this.#scrollableWrapper.current.scrollTo({
      top: 0,
      behavior: "smooth",
    });
  };
}

customElements.define("hermes-existing-chat-scene", ExistingChatScene);
