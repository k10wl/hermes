/**
 * @param {string} leaf
 * @returns {string}
 */
function apiPathnameV1(leaf) {
  return "/api/v1/" + leaf;
}

export const config = {
  server: {
    pathnames: {
      healthCheck: apiPathnameV1("health-check"),
      events: apiPathnameV1("ws"),
    },
  },
};
