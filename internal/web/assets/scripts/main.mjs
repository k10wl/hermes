import { initCustomElements } from "./lib/custom-elements/init.mjs";
import { ServerEvents } from "./lib/events/server-events.mjs";
import { initConnectionIndicator } from "./ui/connection-indicator.mjs";

initConnectionIndicator();
initCustomElements();
// TODO move to UI errors after those are ready
ServerEvents.on("server-error", (event) =>
  console.error("server error:", event.payload),
);
