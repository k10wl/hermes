/**
 * @param {string} pathname
 * @returns {string}
 * */
export function currentUrl(pathname = "") {
  let url = `${window.location.protocol}//${window.location.host}${pathname}`;
  return url;
}
