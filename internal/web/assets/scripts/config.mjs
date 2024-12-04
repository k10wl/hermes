import { currentUrl } from "./lib/current-url.mjs";

/**
 * @param {string} leaf
 * @returns {string}
 */
function apiPathnameV1(leaf) {
  return currentUrl("/api/v1/" + leaf);
}

export const config = {
  server: {
    pathnames: {
      healthCheck: apiPathnameV1("health-check"),
      webSocket: apiPathnameV1("ws"),
    },
  },
  chats: {
    paginationLimit: 5,
  },
};
