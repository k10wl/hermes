import { ServerEvents } from "./events/server-events.mjs";

ServerEvents.on("reload", () => window.location.reload());
