import { ServerEvents } from "./events/server-events.mjs";
import { initConnectionIndicator } from "./ui/connection-indicator.mjs";
import { initCustomElements } from "./ui/custom-elements/init.mjs";

initConnectionIndicator();
initCustomElements();
ServerEvents.on("server-error", (event) =>
  console.error("server error:", event.payload),
);
