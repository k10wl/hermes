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
    const input = assertInstance(messageInput, HTMLTextAreaElement);
    input.value = "";
    input.rows = 1;
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

  const chatsSelector = document.querySelectorAll("#chats-list a");
  for (let i = 0; i < chatsSelector.length; i++) {
    const el = assertInstance(chatsSelector.item(i), HTMLAnchorElement);
    if (el.pathname === window.location.pathname) {
      el.classList.add("primary-bg");
    }
  }

  const resizeObs = new ResizeObserver(autoresize(messageInput));
  messageInput.addEventListener("input", autoresize(messageInput));
  window.addEventListener("resize", autoresize(messageInput));
  resizeObs.observe(messageInput);
});
