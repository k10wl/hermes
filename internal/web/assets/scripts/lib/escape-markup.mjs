/**
 * Escapes a string for use in HTML output.
 * @param {string} unsafe
 * @returns {string}
 */
export function escapeMarkup(unsafe) {
  return unsafe
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;");
}
