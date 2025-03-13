import { AssertInstance } from "/assets/scripts/lib/assert.mjs";
import { Bind, html } from "/assets/scripts/lib/libdim.mjs";

export class ExistingChatScene extends HTMLElement {
  constructor() {
    super();
  }

  connectedCallback() {
    const scrollableWrapper = new Bind((el) =>
      AssertInstance.once(el, HTMLDivElement),
    );

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

        .message {
          border: 1px solid var(--bg-2);
          padding: 4px 8px;
          margin: 12px;
          width: fit-content;
          max-width: 80%;
          border-radius: 10px;
          background: var(--bg-1);

          pre {
            margin: 0;
            text-wrap: wrap;
            word-break: break-all;
            user-select: text;
          }
        }

        .role-assistant {
          color: var(--text-0);
          border-bottom-left-radius: 0;
          border-color: rgb(from var(--primary) r g b / 0.33);
        }

        .role-user {
          border-bottom-right-radius: 0;
          margin-left: auto;
          color: var(--light-bg-0);
        }

        .input-form-wrapper {
          display: flex;
          justify-content: center;
        }

        hermes-message-form {
          padding: 4px 16px 16px;
          max-width: var(--container);
          width: 100%;
        }
      </style>

      <main>
        <div
          id="scrollable-message-wrapper"
          bind="${scrollableWrapper}"
          onscroll="${() => {
            const scrollable = scrollableWrapper.current;
            if (scrollable.scrollTop >= 0) {
              scrollable.classList.remove("scrolled");
            } else {
              scrollable.classList.add("scrolled");
            }
          }}"
        >
          <div id="to-bottom-wrapper">
            <button
              type="button"
              aria-label="Scroll to latest message"
              id="to-bottom"
              onclick="${() => {
                scrollableWrapper.current.scrollTo({
                  top: 0,
                  behavior: "smooth",
                });
              }}"
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
}

customElements.define("hermes-existing-chat-scene", ExistingChatScene);
