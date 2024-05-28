document.addEventListener("DOMContentLoaded", () => {
  const messageForm = assertInstance(
    document.getElementById("message-form"),
    HTMLFormElement,
  );
  const messagesList = assertInstance(
    document.getElementById("messages-list"),
    HTMLDivElement,
  );
  const messageInput = assertInstance(
    document.getElementById("message-input"),
    HTMLTextAreaElement,
  );
  const messageSubmitButton = assertInstance(
    document.getElementById("message-submit-button"),
    HTMLButtonElement,
  );

  document.addEventListener("keydown", (e) => {
    if (
      [document.body, null].find((el) => el === document.activeElement) &&
      e.key === "Enter"
    ) {
      assertInstance(messageInput, HTMLTextAreaElement).focus();
      assertInstance(messageSubmitButton, HTMLButtonElement).click();
      e.preventDefault();
      return;
    }
    if (
      document.activeElement === null ||
      ["INPUT", "TEXTAREA"].includes(document.activeElement.tagName) ||
      [e.shiftKey, e.altKey, e.metaKey, e.ctrlKey].find(Boolean) ||
      [
        "Enter",
        "Tab",
        "ArrowLeft",
        "ArrowRight",
        "ArrowTop",
        "ArrowBottom",
        " ",
      ].includes(e.key)
    ) {
      return;
    }
    assertInstance(messageInput, HTMLTextAreaElement).focus();
  });

  messageForm.addEventListener("submit", async (e) => {
    e.preventDefault();
    const data = new FormData(assertInstance(messageForm, HTMLFormElement));
    const content = data.get("content");
    if (!content || typeof content !== "string" || content.trim() === "") {
      return;
    }
    assertInstance(messageInput, HTMLTextAreaElement).value = "";
    assertInstance(messagesList, HTMLDivElement).append(
      Templates.createMessage(content, "user"),
    );
    let pathname = window.location.pathname;
    if (pathname === "/") {
      pathname = "/chats";
    }
    const res = await fetch(pathname, {
      method: "POST",
      body: data,
    });
    const html = await res.text();
    if (res.status === 301 && res.headers.get("Eval") === "js") {
      // yeah nah, this should not be done, but whatever
      eval(html);
    }
    const parsed = new DOMParser().parseFromString(html, "text/html");
    assertInstance(messagesList, HTMLDivElement).append(
      ...parsed.body.childNodes,
    );
  });

  messageInput.addEventListener("keydown", (e) => {
    if (
      e.key === "Enter" &&
      [e.shiftKey, e.metaKey, e.ctrlKey].every((el) => el === false)
    ) {
      e.preventDefault();
      assertInstance(messageSubmitButton, HTMLButtonElement).click();
    }
  });
});

class Templates {
  static getMessage() {
    return assertInstance(
      document.getElementById("template-message"),
      HTMLTemplateElement,
    );
  }

  /**
   * @param {string} content
   * @param {"user" | "assistant" | "system" } role
   * */
  static createMessage(content, role) {
    const message = assertInstance(
      Templates.getMessage().content.cloneNode(true),
      DocumentFragment,
    );
    const div = assertInstance(message.querySelector("div"), HTMLDivElement);
    const pre = assertInstance(message.querySelector("pre"), HTMLPreElement);
    pre.innerText = content;
    div.classList.add(`role-${role}`);
    return message;
  }
}

/**
 * @template T
 * @returns {T}
 * @param {unknown} obj
 * @param {new () => T} type
 */
function assertInstance(obj, type) {
  if (obj instanceof type) {
    /** @type {any} */
    const any = obj;
    /** @type {T} */
    const t = any;
    return t;
  }
  throw new Error(`Object ${obj} does not have the right type '${type}'!`);
}
