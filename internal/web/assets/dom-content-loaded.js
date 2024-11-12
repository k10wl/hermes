document.addEventListener("DOMContentLoaded", () => {
  const messageForm = assertInstance(
    document.getElementById("message-form"),
    HTMLFormElement,
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
    const input = assertInstance(messageInput, HTMLTextAreaElement);
    if (e.key === "Escape") {
      if (document.activeElement === input) {
        input.blur();
      } else if (window.location.pathname !== "/") {
        setTimeout(() =>
          window.location.replace(
            window.location.protocol + "//" + window.location.host,
          ),
        );
      }
      return;
    }

    if (e.altKey && e.key === "ArrowUp") {
      return chatNavigation("prev");
    }

    if (e.altKey && e.key === "ArrowDown") {
      return chatNavigation("next");
    }

    if (
      [document.body, null].find((el) => el === document.activeElement) &&
      e.key === "Enter"
    ) {
      input.focus();
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
      el.scrollIntoView(false);
    }
  }

  const resizeObs = new ResizeObserver(autoresize(messageInput));
  messageInput.addEventListener("input", autoresize(messageInput));
  window.addEventListener("resize", autoresize(messageInput));
  resizeObs.observe(messageInput);
});

/** @param {"prev" | "next"} dir */
function chatNavigation(dir) {
  const el = assertInstance(
    document.querySelector(".chat-link.primary-bg"),
    HTMLAnchorElement,
  );
  let reqEl =
    dir === "next" ? el.nextElementSibling : el.previousElementSibling;
  if (!reqEl) {
    return;
  }
  window.location.replace(assertInstance(reqEl, HTMLAnchorElement).href);
}
