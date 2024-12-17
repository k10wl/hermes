import { LocationControll } from "./lib/location-control.mjs";
import { ShortcutManager } from "./lib/shortcut-manager.mjs";

ShortcutManager.keydown("<Escape>", (event) => {
  const target = ShortcutManager.getTarget(event);
  if (target === window.document.body) {
    LocationControll.navigate("/");
    return;
  }
  event.stopPropagation();
  event.preventDefault();
  target?.blur();
});
