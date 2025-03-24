import { config } from "./config.mjs";
import { ServerEvents } from "./lib/events/server-events.mjs";

ServerEvents.on("reload", () => window.location.reload());

ServerEvents.__init(config.server.pathnames.webSocket);
Reflect.set(window, "ServerEvents", ServerEvents);
