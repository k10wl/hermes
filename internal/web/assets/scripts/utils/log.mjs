/**
 * @param {any[]} args
 */
export function log(...args) {
  console.log(...args);
}

export function init() {
  document.addEventListener("DOMContentLoaded", () => {
    console.log("DOMContentLoaded");
  });
}
