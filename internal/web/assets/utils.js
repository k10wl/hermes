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
 * @param {new (data: any) => T} type
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

/** @param {HTMLElement} el */
function isScrollable(el) {
  return el.scrollWidth > el.clientWidth || el.scrollHeight > el.clientHeight;
}
/** @param {HTMLTextAreaElement} el */
const autoresize = (el) => () => {
  if (isScrollable(el) && 10 > el.rows) {
    el.rows = el.rows + 1;
  }
  if (!el.value) {
    el.rows = 1;
  }
};
