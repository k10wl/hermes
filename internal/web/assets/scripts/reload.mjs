import { ServerEvents } from "./lib/events/server-events.mjs";

ServerEvents.on("reload", () => window.location.reload());
