import { Bind, html } from "/assets/scripts/lib/libdim.mjs";
import { Template } from "/assets/scripts/models.mjs";

import { AssertInstance, AssertString } from "../../assert.mjs";
import {
  DeleteTemplateEvent,
  RequestEditTemplateEvent,
  RequestReadTemplateEvent,
} from "../../events/client-events-list.mjs";
import { ServerEvents } from "../../events/server-events.mjs";
import { ServerErrorEvent } from "../../events/server-events-list.mjs";
import { FocusOnKeydown } from "../../focus-on-keydown.mjs";
import { LocationControll } from "../../location-control.mjs";
import { ShortcutManager } from "../../shortcut-manager.mjs";
import { ResizableTextInput } from "../content-editable-plain-text.mjs";
import { Action, ActionStore } from "../control-panel.mjs";
import { AlertDialog, ConfirmDialog, HermesDialog } from "../dialog.mjs";

/** @type {null | HermesViewTemplateScene} */
export let template = null;

class TemplateUpdatedDialog extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanup = [];

  /** @type {Template | null} */
  template = null;

  dialog = new Bind((el) => AssertInstance.once(el, HermesDialog));

  constructor() {
    super();
    this.attachShadow({ mode: "closed" }).append(html`
      <style>
        * {
          color: var(--text-0);
          user-select: none;
        }

        #actions {
          display: flex;
          justify-content: flex-end;
          gap: 0.5rem;
          h-button::part(button) {
            font-size: 1.1rem;
          }
        }
      </style>

      <h-dialog bind="${this.dialog}">
        <h-dialog-card>
          <h-dialog-title>Template outdated</h-dialog-title>
          <h-dialog-block>
            <p>
              Content was changed. Newer version of content is available on the
              server
            </p>
            <p>Override current content with latest value?</p>
          </h-dialog-block>

          <h-dialog-block>
            <div id="actions">
              <h-button onclick="${() => this.#dispatch("cancel")}"
                >Cancel <h-key>n</h-key></h-button
              >
              <h-button
                onclick="${() => this.#dispatch("confirm")}"
                variant="primary"
                >OK <h-key>y</h-key></h-button
              >
            </div>
          </h-dialog-block>
        </h-dialog-card>
      </h-dialog>
    `);
  }

  /** @param {"cancel" | "confirm"} eventName */
  #dispatch = (eventName) => {
    this.dispatchEvent(new Event(eventName));
  };

  /** @param {Template} template */
  showModal(template) {
    if (this.dialog.current.element.open) {
      return;
    }
    Template.validator.check(template);
    this.template = template;
    this.#cleanup.push(
      ShortcutManager.keydown("<KeyY>", (event) => {
        event.preventDefault();
        event.stopPropagation();
        this.#dispatch("confirm");
      }),
      ShortcutManager.keydown("<KeyN>", (event) => {
        event.preventDefault();
        event.stopPropagation();
        this.#dispatch("cancel");
      }),
    );
    this.dialog.current.element.showModal();
  }

  close() {
    if (!this.dialog.current.element.open) {
      return;
    }
    this.#cleanup.forEach((cb) => cb());
    this.dialog.current.element.close();
  }

  disconnectedCallback() {
    template = null;
    this.#cleanup.forEach((cb) => cb());
  }
}
customElements.define("h-teplate-updated-dialog", TemplateUpdatedDialog);

class NameCollisionDialog extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanup = [];

  dialog = new Bind((el) => AssertInstance.once(el, HermesDialog));

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "open" });
  }

  showModal() {
    this.#cleanup.push(
      ShortcutManager.keydown("<KeyN>", (event) => {
        event.preventDefault();
        event.stopPropagation();
        this.#dispatch("cancel");
      }),
      ShortcutManager.keydown("<KeyC>", (event) => {
        event.preventDefault();
        event.stopPropagation();
        this.#dispatch("clone");
      }),
      ShortcutManager.keydown("<KeyY>", (event) => {
        event.preventDefault();
        event.stopPropagation();
        this.#dispatch("rename");
      }),
    );
    AssertInstance.once(
      this.dialog.current.element,
      HTMLDialogElement,
    ).showModal();
  }

  close() {
    AssertInstance.once(this.dialog.current.element, HTMLDialogElement).close();
    this.#cleanup.forEach((cb) => cb());
  }

  /** @param {"cancel" | "clone" | "rename"} name */
  #dispatch = (name) => {
    this.dispatchEvent(new CustomEvent(name));
  };

  connectedCallback() {
    this.shadow.append(html`
      <style>
        * {
          color: var(--text-0);
          user-select: none;
        }

        #actions {
          display: flex;
          justify-content: flex-end;
          gap: 0.5rem;
          h-button::part(button) {
            font-size: 1.1rem;
          }
        }</style
      ><style>
        * {
          color: var(--text-0);
          user-select: none;
        }

        #actions {
          display: flex;
          justify-content: flex-end;
          gap: 0.5rem;
          h-button::part(button) {
            font-size: 1.1rem;
          }
        }
      </style>

      <h-dialog
        onclose="${() => this.#dispatch("cancel")}"
        bind="${this.dialog}"
      >
        <h-dialog-card>
          <h-dialog-title>Template name changed</h-dialog-title>
          <h-dialog-block>
            <p>
              Edited template has different name than initial, so changes can be
              cloned as new template
            </p>
            <p>Rename original?</p>
          </h-dialog-block>

          <h-dialog-block>
            <div id="actions">
              <h-button onclick="${() => this.#dispatch("cancel")}"
                >Cancel <h-key>n</h-key></h-button
              >
              <h-button
                variant="primary"
                onclick="${() => this.#dispatch("clone")}"
                >Clone <h-key>c</h-key></h-button
              >
              <h-button
                variant="error"
                onclick="${() => this.#dispatch("rename")}"
                >Rename <h-key>y</h-key></h-button
              >
            </div>
          </h-dialog-block>
        </h-dialog-card>
      </h-dialog>
    `);
  }

  disconnectedCallback() {
    this.#cleanup.forEach((cb) => cb());
  }
}

customElements.define("h-name-collision-dialog", NameCollisionDialog);

export class HermesViewTemplateScene extends HTMLElement {
  /** @type {(() => void)[]} */
  #cleanup = [];

  /** @type {import("/assets/scripts/models.mjs").Template | null} */
  #template = null;

  #saveButtonText = new Bind((el) => AssertInstance.once(el, HTMLSpanElement));
  nameCollisionDialog = new Bind((el) =>
    AssertInstance.once(el, NameCollisionDialog),
  );
  templateUpdatedDialog = new Bind((el) =>
    AssertInstance.once(el, TemplateUpdatedDialog),
  );
  form = new Bind((el) => AssertInstance.once(el, HTMLFormElement));
  #content = new Bind((el) => AssertInstance.once(el, ResizableTextInput));
  #focus = new FocusOnKeydown();

  constructor() {
    super();
    this.shadow = this.attachShadow({ mode: "closed" });
    this.shadow.append(this.#html);
  }

  connectedCallback() {
    template = this;
    this.#sendReadRequest();
    this.#cleanup.push(
      ActionStore.add(new Action("template: delete template", this.delete)),
      ActionStore.add(
        new Action("template: save edit", () =>
          this.form.current.requestSubmit(),
        ),
      ),
      ServerEvents.on(["template-changed"], (event) => {
        if (
          !this.#content ||
          this.#content.current.value === event.payload.template.content
        ) {
          return;
        }
        this.templateUpdatedDialog.current.showModal(event.payload.template);
      }),
    );
  }

  disconnectedCallback() {
    template = null;
    this.#cleanup.forEach((cb) => cb());
  }

  #sendReadRequest = () => {
    const readEvent = new RequestReadTemplateEvent({
      id: LocationControll.templateId || -1,
    });
    ServerEvents.send(readEvent);
    const off = ServerEvents.on(["read-template", "server-error"], (event) => {
      if (event.id !== readEvent.id) {
        return;
      }
      off();
      if (event instanceof ServerErrorEvent) {
        LocationControll.navigate("/templates");
        AlertDialog.instance.alert({
          title: "Failed to read template",
          description: `Failed to read template: ${event.payload}`,
        });
        return;
      }
      this.#template = event.payload.template;
      this.#setDelayedContent(html`
        <form bind="${this.form}" onsubmit="${this.submit}">
          <h-resizable-text-input
            bind="${this.#content}"
            name="content"
            placeholder='--{{define "name"}} dynamic value => --{{.}} --{{end}}'
          ></h-resizable-text-input>
          <input
            type="hidden"
            name="initial name"
            value="${event.payload.template.name}"
          />
          <h-button onclick="${() => this.form.current.requestSubmit()}">
            <span bind="${this.#saveButtonText}"> Save </span>
            &nbsp;
            <h-key>Meta-S</h-key>
          </h-button>
        </form>
      `);
      this.#processForm();
      this.#content.current.content.focus();
      this.#content.current.value = event.payload.template.content.trim();
      this.#focus.attach(this.#content.current.content);
      this.#cleanup.push(
        () => this.#focus.detach(),
      );
    });
    this.#cleanup.push(off);
  };

  delete = async () => {
    const ok = await ConfirmDialog.instance.confirm({
      title: "Delete template",
      description: `Are you sure you want to delete template '${this.#template?.name}'?`,
    });
    if (!ok) {
      return;
    }
    const deleteEvent = new DeleteTemplateEvent({
      name: AssertString.check(this.#template?.name),
    });
    ServerEvents.on(["server-error", "template-deleted"], (event) => {
      if (event.id !== deleteEvent.id) {
        return;
      }
      if (event instanceof ServerErrorEvent) {
        AlertDialog.instance.alert({
          title: "Failed to delete template",
          description: `Delete errored: ${event.payload}`,
        });
        return;
      }
      LocationControll.navigate("/templates/");
    });
    ServerEvents.send(deleteEvent);
  };

  /** @param {boolean} clone */
  #save = (clone) => {
    const template = AssertInstance.once(this.#template, Template);
    const content = AssertString.check(this.#content.current.value);
    if (template.content === content) {
      return;
    }
    const edit = new RequestEditTemplateEvent({
      name: template.name,
      content,
      clone: clone || LocationControll.pathname.startsWith("/templates/new"),
    });
    ServerEvents.send(edit);
    const off = ServerEvents.on(
      ["template-changed", "template-created", "server-error"],
      (event) => {
        if (event.id !== edit.id) {
          return;
        }
        if (event instanceof ServerErrorEvent) {
          AlertDialog.instance.alert({
            description: `Edit failed - ${event.payload}`,
          });
          return;
        }
        this.#content.current.value = edit.payload.content;
        this.#template = event.payload.template;
        this.nameCollisionDialog.current.close();
        this.#savedIndicator();
        LocationControll.navigate("/templates/" + this.#template.id, false);
        off();
      },
    );
  };

  #savedIndicator() {
    const text = this.#saveButtonText.current.textContent;
    this.#saveButtonText.current.textContent = "Saved";
    setTimeout(() => {
      this.#saveButtonText.current.textContent = text;
    }, 2000);
  }

  /**
   * @param {DocumentFragment} html
   */
  #setDelayedContent = (html) => {
    AssertInstance.once(
      this.shadow.querySelector("main"),
      HTMLElement,
      "parent element should be present to prevent flickering while loading",
    ).append(html);
  };

  submit = () => {
    const newName = /"(?<name>.*?)"/.exec(
      AssertString.check(
        this.#content.current.value,
        "expected text input to have string value",
      ),
    );
    const nameChanged =
      AssertString.check(
        this.#template?.name,
        "expected component to hold on initial template name",
      ) !==
      AssertString.check(
        newName?.groups?.name,
        "expected new name to be retrieved from content",
      );
    if (nameChanged && LocationControll.templateId) {
      AssertInstance.once(
        this.nameCollisionDialog.current,
        NameCollisionDialog,
      ).showModal();
      return;
    }
    this.#save(false);
  };

  #processForm = () => {
    const form = AssertInstance.once(
      this.form.current,
      HTMLFormElement,
      "form must be present before process call",
    );

    form.addEventListener("submit", (event) => {
      event.preventDefault();
    });

    this.#cleanup.push(
      ShortcutManager.keydown("<M-KeyS>", (event) => {
        event.preventDefault();
        form.requestSubmit();
      }),
    );
  };

  #html = html`
    <style>
      * {
        box-sizing: border-box;
        color: var(--text-0);
      }

      main {
        height: 100vh;
        max-height: 100vh;
        display: grid;
        place-items: center;
        overflow: auto;
      }

      form {
        max-width: var(--container-max-width);
        margin: var(--container-margin);
        display: flex;
        width: 100%;
        flex-flow: column nowrap;
        h-button {
          align-self: flex-end;
        }
      }

      h-resizable-text-input {
        --padding: 0.5rem 1rem;
        --border-color: var(--bg-2);
        --border: 1px solid var(--border-color);
        width: 100%;

        &:focus-within {
          --border-color: var(--primary);
        }

        &::part(wrapper) {
          border: var(--border);
          border-radius: 1.25rem;
          overflow: hidden;
          color: var(--text);
        }

        &::part(content) {
          padding: var(--padding);
        }

        &::part(placeholder) {
          padding: var(--padding);
        }
      }
    </style>

    <main>
      <!--content is delayed until template data is ready to prevent flickering-->
    </main>

    <h-name-collision-dialog
      oncancel="${() => this.nameCollisionDialog.current.close()}"
      onclone="${() => this.#save(true)}"
      onrename="${() => this.#save(false)}"
      bind="${this.nameCollisionDialog}"
    ></h-name-collision-dialog>

    <h-teplate-updated-dialog
      oncancel="${() => this.templateUpdatedDialog.current.close()}"
      onconfirm="${() => {
        this.#content.current.value = AssertString.check(
          this.templateUpdatedDialog.current.template?.content,
          "expected update value to be string",
        );
        this.templateUpdatedDialog.current.close();
      }}"
      bind="${this.templateUpdatedDialog}"
    ></h-teplate-updated-dialog>
  `;
}

customElements.define("hermes-view-template-scene", HermesViewTemplateScene);