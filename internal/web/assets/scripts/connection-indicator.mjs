import { ServerEvents } from "./events/server-events.mjs";

const OFFLINE_SUFIX = " - offline";

function hasOfflineSufix() {
  return window.document.title.endsWith(OFFLINE_SUFIX);
}

function onOnline() {
  if (!hasOfflineSufix()) {
    return;
  }
  window.document.title.replace(OFFLINE_SUFIX, "");
}

function onOffline() {
  if (hasOfflineSufix()) {
    return;
  }
  window.document.title += OFFLINE_SUFIX;
}

export function initOfflineIndicator() {
  ServerEvents.onClose(onOnline);
  ServerEvents.onClose(onOffline);
}
