import "./ui/scene.mjs";
import "./lib/sound-manager.mjs";
import "./lib/custom-elements/init.mjs";
import "./ui/connection-indicator.mjs";

import { AlertDialog } from "./lib/custom-elements/dialog.mjs";
import { ServerEvents } from "./lib/events/server-events.mjs";

ServerEvents.on("server-error", (event) => {
  console.error("server error:", event.payload);
  AlertDialog.instance.alert({
    title: "Server Error",
    description: event.payload,
  });
});
