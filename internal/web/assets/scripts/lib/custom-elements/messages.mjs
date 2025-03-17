import { AssertInstance } from "/assets/scripts/lib/assert.mjs";
import { escapeMarkup } from "/assets/scripts/lib/escape-markup.mjs";
import { RequestReadChatEvent } from "/assets/scripts/lib/events/client-events-list.mjs";
import { ServerEvents } from "/assets/scripts/lib/events/server-events.mjs";
import { Bind, html, Signal } from "/assets/scripts/lib/libdim.mjs";
import { LocationControll } from "/assets/scripts/lib/location-control.mjs";
import { Message } from "/assets/scripts/models.mjs";
import Prism from "/assets/scripts/third-party/prism.js";

import { AlertDialog } from "./dialog.mjs";

const COPY_SYMBOL = "⎘";
const COPY_SUCCESS_SYMBOL = "✓";

/**
 * @param {string} text
 * @param {() => void} onSuccess
 * @param {(error: unknown) => void} onError
 */
async function clipboardWriteText(text, onSuccess, onError) {
  try {
    await navigator.clipboard.writeText(text);
    onSuccess();
  } catch (error) {
    console.error(error);
    console.trace();
    onError(error);
  }
}

/** @param {unknown} error */
function clipboardReadError(error) {
  let msg = "unknown error";
  if (error instanceof Error) {
    msg = error.message;
  }
  AlertDialog.instance.alert({
    title: "Copy failed",
    description: msg,
  });
}

/**
 * @param {HTMLElement} element
 * @returns {() => void}
 */
function clipboardReadSuccess(element) {
  return () => {
    element.textContent = COPY_SUCCESS_SYMBOL;
    setTimeout(() => {
      element.textContent = COPY_SYMBOL;
    }, 2000);
  };
}

class CodeBlock extends HTMLElement {
  #slot = new Bind((el) => AssertInstance.once(el, HTMLSlotElement));
  #raw = "";
  #wrapLines = new Signal(false);

  constructor() {
    super();
    this.attachShadow({ mode: "closed" }).append(html`
      <style>
        /* PrismJS 1.30.0 */
        code[class*="language-"],
        pre[class*="language-"] {
          color: #f8f8f2;
          background: 0 0;
          text-shadow: 0 1px rgba(0, 0, 0, 0.3);
          font-family: Consolas, Monaco, "Andale Mono", "Ubuntu Mono", monospace;
          font-size: 1em;
          text-align: left;
          white-space: pre;
          word-spacing: normal;
          word-break: normal;
          word-wrap: normal;
          line-height: 1.5;
          -moz-tab-size: 4;
          -o-tab-size: 4;
          tab-size: 4;
          -webkit-hyphens: none;
          -moz-hyphens: none;
          -ms-hyphens: none;
          hyphens: none;
        }
        pre[class*="language-"] {
          padding: 1em;
          margin: 0.5em 0;
          overflow: auto;
          border-radius: 0.3em;
        }
        :not(pre) > code[class*="language-"],
        pre[class*="language-"] {
          background: #272822;
        }
        :not(pre) > code[class*="language-"] {
          padding: 0.1em;
          border-radius: 0.3em;
          white-space: normal;
        }
        .token.cdata,
        .token.comment,
        .token.doctype,
        .token.prolog {
          color: #8292a2;
        }
        .token.punctuation {
          color: #f8f8f2;
        }
        .token.namespace {
          opacity: 0.7;
        }
        .token.constant,
        .token.deleted,
        .token.property,
        .token.symbol,
        .token.tag {
          color: #f92672;
        }
        .token.boolean,
        .token.number {
          color: #ae81ff;
        }
        .token.attr-name,
        .token.builtin,
        .token.char,
        .token.inserted,
        .token.selector,
        .token.string {
          color: #a6e22e;
        }
        .language-css .token.string,
        .style .token.string,
        .token.entity,
        .token.operator,
        .token.url,
        .token.variable {
          color: #f8f8f2;
        }
        .token.atrule,
        .token.attr-value,
        .token.class-name,
        .token.function {
          color: #e6db74;
        }
        .token.keyword {
          color: #66d9ef;
        }
        .token.important,
        .token.regex {
          color: #fd971f;
        }
        .token.bold,
        .token.important {
          font-weight: 700;
        }
        .token.italic {
          font-style: italic;
        }
        .token.entity {
          cursor: help;
        }
      </style>

      <style>
        :host {
          --_spacing: 0.5rem;
          white-space: wrap;
          overflow: auto;
        }

        #container {
          background: var(--bg-1);
          border-radius: var(--_spacing);
          overflow: hidden;
          border: 1px solid rgb(from var(--text-0) r g b / 0.25);
        }

        #tab {
          background: var(--bg-2);
          display: flex;
          align-items: center;
          justify-content: space-between;
          padding: calc(var(--_spacing) / 2) var(--_spacing);
          * {
            margin: 0;
          }
        }

        ::part(button) {
          display: flex;
          align-items: center;
          justify-content: center;
          --_size: 1rem;
          width: var(--_size);
          height: var(--_size);
        }
        #wrap::part(button) {
          font-size: 0.66rem;
          padding-top: 0.25rem;
        }

        pre {
          padding: var(--_spacing);
          margin: 0;
          overflow: auto;
        }

        code[data-wrap="true"] {
          white-space: pre-wrap;
        }

        #actions {
          display: flex;
          align-items: center;
          justify-content: space-between;
          gap: 0.5rem;
        }
      </style>

      <div id="container">
        <div id="tab">
          <p>${this.getAttribute("data-language")}</p>
          <div id="actions">
            <h-button
              id="wrap"
              title="Toggle wrap lines"
              onclick="${() => {
                this.#wrapLines.value = !this.#wrapLines.value;
              }}"
              >↵</h-button
            >
            <h-button
              id="copy"
              title="Copy to clipboard"
              onclick="${(event) => {
                clipboardWriteText(
                  this.#raw,
                  clipboardReadSuccess(
                    AssertInstance.once(event.target, HTMLElement),
                  ),
                  clipboardReadError,
                );
              }}"
              >${COPY_SYMBOL}</h-button
            >
          </div>
        </div>
        <pre><code
            data-wrap="${this.#wrapLines}"
        ><slot bind="${this.#slot}"></slot></code></pre>
      </div>
    `);
    const replacement = this.#slot.current.assignedNodes();
    this.#raw = replacement
      .map((el) => el.textContent)
      .join("")
      .replace(/\n$/, "");
    this.#slot.current.replaceWith(...replacement);
  }
}
customElements.define("h-message-code-block", CodeBlock);

class ChatMessage extends HTMLElement {
  static #unset = "******** UNSET ********";
  #slot = new Bind((el) => AssertInstance.once(el, HTMLSlotElement));
  #raw = ChatMessage.#unset;

  constructor() {
    super();
    this.attachShadow({ mode: "open" }).append(html`
      <style>
        :host {
          --_spacing: 0.5rem;
        }

        #wrapper {
          display: grid;
          padding: var(--_spacing) 0;
        }
        #content {
          user-select: text;
          display: grid;
          gap: var(--_spacing);
          white-space: pre-wrap;
          word-break: break-word;
          color: var(--_text);
        }
        :host([data-role="user"]) {
          #content {
            --_border-radius: 0.66rem;
            background: rgb(from var(--_primary) r g b / 0.05);
            border: 1px solid rgb(from var(--_primary) r g b / 0.25);
            justify-self: end;
            width: fit-content;
            max-width: 80%;
            border-radius: var(--_border-radius) var(--_border-radius)
              calc(var(--_border-radius) / 3) var(--_border-radius);
            padding: 0.33rem;
            padding-bottom: 0.2rem;
          }
          #actions {
            justify-self: end;
          }
        }

        #actions {
          --_duration: 100ms;
          transition: all var(--_duration);
          transition-delay: var(--_duration);
          margin-top: 0.2rem;
          opacity: 0;
          #copy::part(button) {
            display: inline-flex;
            align-items: center;
            justify-content: center;
            --_size: 2ch;
            width: var(--_size);
            height: var(--_size);
          }
        }
        #wrapper:hover #actions,
        #wrapper:focus-within #actions {
          opacity: 1;
          transition-delay: 0ms;
        }
      </style>

      <div id="wrapper">
        <div id="content"><slot bind="${this.#slot}"></slot></div>
        <div part="actions" id="actions">
          <h-button
            id="copy"
            onclick="${(event) => {
              clipboardWriteText(
                this.#raw,
                clipboardReadSuccess(
                  AssertInstance.once(event.target, HTMLElement),
                ),
                clipboardReadError,
              );
            }}"
            >${COPY_SYMBOL}</h-button
          >
        </div>
      </div>
    `);
  }

  connectedCallback() {
    this.#slot.current.outerHTML = this.#parseContent(this.#extractContent());
  }

  #extractContent() {
    if (this.#raw !== ChatMessage.#unset) {
      return this.#raw;
    }
    const content = this.#slot.current.assignedNodes()[0]?.textContent;
    if (typeof content !== "string") {
      throw new Error("unexpected content");
    }
    this.#raw = content;
    return this.#raw;
  }

  /**
   * @param {string} content
   * @returns {string}
   */
  #parseContent(content) {
    if (this.getAttribute("data-role") === "user") {
      return this.#parseUser(content);
    }
    return this.#parseAssistant(content);
  }

  /** @param {string} content
   * @returns {string}
   */
  #parseUser(content) {
    return escapeMarkup(content);
  }

  static #matchMarkdown =
    /(?<markup>```(?<lang>\w+)\n(?<code>(.*\n)+?)```(\n)?)/gm;
  /**
   * @param {string} content
   * @returns {string}
   */
  #parseAssistant(content) {
    const matches = content.matchAll(ChatMessage.#matchMarkdown);

    let pointer = 0;
    let escapedMessage = "";
    for (const execArray of matches) {
      const escapedPrefix = escapeMarkup(
        execArray.input.slice(pointer, execArray.index),
      );
      const code = execArray.groups?.code;
      if (!code) {
        throw new Error("code not found");
      }
      const markup = execArray.groups?.markup;
      if (!markup) {
        throw new Error("markup not found");
      }
      pointer += execArray.index + markup.length;
      const lang = execArray.groups?.lang;
      if (!lang) {
        throw new Error("lang not found");
      }
      // @ts-expect-error languages are not complete
      const pr = Prism.languages[lang ?? "plain"];
      if (!pr) {
        escapedMessage += escapedPrefix;
      }
      const res = Prism.highlight(code, pr, lang);
      escapedMessage +=
        escapedPrefix +
        `<h-message-code-block data-language="${lang}">${res}</h-message-code-block>`;
    }
    escapedMessage += escapeMarkup(content.slice(pointer));

    return escapedMessage;
  }
}
customElements.define("h-chat-message", ChatMessage);

export class Messages extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanupOnDisconnect = [];

  messages = new Bind((el) => AssertInstance.once(el, HTMLElement));

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
    this.shadow.append(html`
      <style>
        :host {
          --_text: var(--text-0);
          --_primary: var(--primary);
        }

        h-chat-message[data-role="assistant"]:last-child::part(actions) {
          opacity: 1;
        }
      </style>

      <section bind="${this.messages}"></section>
    `);
  }

  connectedCallback() {
    const messagesContainer = AssertInstance.once(
      this.messages.current,
      HTMLElement,
    );
    const routeObserver = new RouteObserver();
    routeObserver.notify();
    this.#cleanupOnDisconnect.push(
      LocationControll.attach(routeObserver),

      ServerEvents.on("message-created", async (data) => {
        if (data.payload.chat_id !== LocationControll.chatId) {
          return;
        }
        const newMessage = AssertInstance.once(
          this.#messageToHtml(data.payload.message).firstElementChild,
          HTMLElement,
        );
        messagesContainer.append(newMessage);
        if (data.payload.message.role !== "assistant") {
          return;
        }
        const { top } = newMessage.getBoundingClientRect();
        if (top < 0) {
          newMessage.scrollIntoView({
            behavior: "instant",
            block: "start",
          });
        }
      }),

      ServerEvents.on("read-chat", (data) => {
        messagesContainer.replaceChildren(
          ...data.payload.messages.map((message) =>
            this.#messageToHtml(message),
          ),
        );
      }),
    );
  }

  disconnectedCallback() {
    for (const cleanup of this.#cleanupOnDisconnect) {
      cleanup();
    }
  }

  /** @param {Message} message */
  #messageToHtml(message) {
    const { role, content } = Message.validator.check(message);
    return html`
      <h-chat-message data-role="${role}"
        >${document.createTextNode(content)}</h-chat-message
      >
    `;
  }
}

class RouteObserver {
  /** @type {number | null} */
  #last = null;

  constructor() {}

  notify() {
    const chatId = LocationControll.chatId;
    const prev = this.#last;
    this.#last = chatId;
    if (prev === chatId || !chatId) {
      return;
    }
    ServerEvents.send(new RequestReadChatEvent(chatId));
  }
}
