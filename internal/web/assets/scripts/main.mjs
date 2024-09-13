import { config } from "./config.mjs";
import { ServerEvents } from "./events/server-events.mjs";

const serverEvents = new ServerEvents(config.server.pathnames.events);
serverEvents.on("reload", () => window.location.reload());
