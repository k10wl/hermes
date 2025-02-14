class _Theme {
  /** @type _Settings */
  #settings;

  /** @param {_Settings} settings */
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

class _Settings {
  /**
   * @typedef settings
   * @property {boolean} dark_mode
   * @property {boolean} initted
   */

  /** @type _Settings | undefined */
  static instance;
  #settings;

  /** @param {settings} settings */
  constructor(settings) {
    this.#settings = settings.initted ? settings : this.#init();
  }

  /** @param {settings} settings */
  static init(settings) {
    if (this.instance) {
      throw new Error("settings already initialized");
    }
    this.instance = new _Settings(settings);
    return this.instance;
  }

  /**
   * @public
   * @returns {settings}
   */
  get() {
    return this.#settings;
  }

  /** @returns {settings} */
  #init() {
    /** @type settings */
    const settings = {
      initted: true,
      dark_mode:
        window.matchMedia &&
        window.matchMedia("(prefers-color-scheme: dark)").matches,
    };
    this.update(settings);
    return settings;
  }

  /**
   * @param {settings} settings
   * @returns {Promise<boolean>}
   */
  async update(settings) {
    try {
      await fetch("/settings", {
        method: "PUT",
        body: JSON.stringify(settings),
      });
      this.#settings = settings;
      return true;
    } catch (e) {
      console.error(e);
      return false;
    }
  }
}
