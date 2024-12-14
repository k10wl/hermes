import { AssertInstance } from "./lib/assert.mjs";
import { LocationControll } from "./lib/location-control.mjs";
import { ShortcutManager } from "./lib/shortcut-manager.mjs";

ShortcutManager.keydown("<Escape>", (event) => {
  try {
    const target = AssertInstance.once(event.target, HTMLElement);
    if (target === window.document.body) {
      LocationControll.navigate("/");
      return;
    }
    event.stopPropagation();
    event.preventDefault();
    target.blur();
  } catch {
    // just don't explode
  }
});
