export class SoundManager {
  /** @typedef {keyof typeof SoundManager['availableSounds']} SoundName */

  static #path = "/assets/sounds/";
  static availableSounds = {
    "message-in-global": "message-in-global.mp3",
    "message-in-local": "message-in-local.mp3",
  };

  /** @param {SoundName} name */
  static play(name) {
    const src = SoundManager.#getSrc(name);
    const audio = document.createElement("audio");
    SoundManager.#setupAudioTeardown(audio);
    audio.src = src;
    audio.play();
  }

  /** @param {SoundName} name */
  static #getSrc(name) {
    const filename = SoundManager.availableSounds[name];
    if (!filename) {
      throw new Error(`unknown audio "${name}"`);
    }
    return this.#path + filename;
  }

  /** @param {HTMLAudioElement} audio  */
  static #setupAudioTeardown(audio) {
    const teardown = () => audio.remove();
    audio.onended = teardown;
    audio.onerror = teardown;
  }
}
