class Settings {
  /**
   * @typedef settings
   * @property {boolean} dark_mode
   * @property {boolean} initted
   */

  /** @type Settings | undefined */
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
    this.instance = new Settings(settings);
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
