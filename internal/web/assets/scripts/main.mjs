import "./ui/scene.mjs";
import "./lib/sound-manager.mjs";
import "./lib/custom-elements/init.mjs";
import "./lib/html-v2.mjs";
import "./ui/connection-indicator.mjs";

import { ServerEvents } from "./lib/events/server-events.mjs";

// TODO move to UI errors after those are ready
ServerEvents.on("server-error", (event) =>
  console.error("server error:", event.payload),
);
