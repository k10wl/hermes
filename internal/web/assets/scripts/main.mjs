import { initChats } from "./ui/chats/init.mjs";
import { initConnectionIndicator } from "./ui/connection-indicator.mjs";
import { initCustomElements } from "./ui/custom-elements/init.mjs";

initConnectionIndicator();
initCustomElements();

document.addEventListener("DOMContentLoaded", function () {
  initChats();
});
