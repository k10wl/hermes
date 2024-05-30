class Theme {
  /** @type Settings */
  #settings;

  /** @param {Settings} settings */
  constructor(settings) {
    this.#settings = settings;
  }

  load() {
    this.adjustColors(this.isDark() ? "dark" : "light");
  }

  isDark() {
    return this.#settings.get().dark_mode;
  }

  /**
   * @param {"light" | "dark"} variant
   */
  adjustColors(variant) {
    document.documentElement.style.setProperty(
      "--bg-0",
      `var(--${variant}-bg-0)`,
    );
    document.documentElement.style.setProperty(
      "--bg-1",
      `var(--${variant}-bg-1)`,
    );
    document.documentElement.style.setProperty(
      "--bg-2",
      `var(--${variant}-bg-2)`,
    );
    document.documentElement.style.setProperty(
      "--text-0",
      `var(--${variant}-text-0)`,
    );
  }

  swithTheme() {
    const wasDark = this.isDark();
    this.adjustColors(wasDark ? "light" : "dark");
    this.#settings.update({ ...this.#settings.get(), dark_mode: !wasDark });
  }
}
