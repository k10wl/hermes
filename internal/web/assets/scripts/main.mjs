import { ServerEvents } from "./events/server-events.mjs";
import { initConnectionIndicator } from "./ui/connection-indicator.mjs";
import { initCustomElements } from "./ui/custom-elements/init.mjs";

initConnectionIndicator();
initCustomElements();
// TODO move to UI errors after those are ready
ServerEvents.on("server-error", (event) =>
  console.error("server error:", event.payload),
);
