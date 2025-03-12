# Host-based Extensible Response Management Execution System (HERMES)

## Overview
Hermes provides unified local point entry to popular AI providers.
You can call completion using desired model in following ways:
- CLI, app is build to use stdin, stdout, stderr. Commands can be chained and work with other tools from terminal. Use `hermes --help` for more info regarding CLI usage;
- Web, app ships HTTP server to use all features from browser GUI. Use `hermes serve --help` for more info;

## Name derives from following valuable key features:
- **Host-based**: Stores data locally, allows single entry point to multiple third party providers.
- **Extensible**: Provides a customizable API that covers essential use cases and allows to update functionalities as needed.
- **Response Management**: Manages and streamlines AI responses with intuitive UX.
- **Execution System**: Includes robust capabilities for templating prompts, chaining API calls, and setting up automated workflows.

--------------------------------------------------------------------------------

## Instalation

### Using golang (requires go 1.24.0)
You can add it to you golang managed library using
```bash
make install
```
Or build binaries into ./bin/ directory by yourself
```bash
make build
```

### Environmental variables:
```bash
# Models
HERMES_OPENAI_API_KEY=    # key to OpenAI API, optional
HERMES_ANTHROPIC_API_KEY= # key to Anthropic API, optional

# Files
HERMES_DB_DNS=            # sqlite3 dns entry to persist chats, templates, messages
                          # optional, defaults to /hermes/main.db in your config dir
                          # rel: https://pkg.go.dev/os#UserConfigDir

# DEBUG
HERMES_MOCK_COMPLETION=   # debug mode, chats will return inputted message
                          # for development `make watch` this value to true
                          # not to be use in any real scenario
```

--------------------------------------------------------------------------------

## Templates
One of core features of Hermes are templates. Golang provides flexible and extensible way to manage text contents. It is used as templating engine. Having such powerful utility allows to create reusable prompts with custom instructions. 

After executing this CLI command:
```bash
hermes template upsert --content\
    '--{{define "short"}}--{{.}} (Short answer please)--{{end}}'
```
You will have stored your first template named "short". It can be now used in chats!
```bash
hermes chat --content 'what is the mass of the Sun?' --template short
#                     name inherited from template definition   ^^^^^
```
AI will receive following - `what is the mass of the Sun? (Short answer please)`.
And now we have reusable template that will make response concise. You can you any prompts with any instructions, AI is really pleasant to use when there is no need to repeat yourself

This can be any set of instructions that are repeated or needs to be tested. This is Golang templating. The single difference is syntax. Golang uses `[[` as opening brackets and `]]` as closing, but it conflicted with a lot of formatting, therefore I chose `--{{` as opening and `}}` as closing. Any other expected build in feature works


Hermes ships with copy of this README as one of templates, you can always ask questions related this application using `hermes-help`
```bash
hermes chat --content 'how to I change database?' --template hermes-help
```

--------------------------------------------------------------------------------

## Development

### Internals
Application has decent code coverage over most valuable and sensitive parts, you can always run `make test` to make sure that defined parts of application work as expected. You can skip JS testing part if you run `make pre-test && make test-app`.

### Workflows
CI is configured, pull requests will run tests, refer to github actions result as to source of truth.

### Server (ignore this if you are not going to touch js parts)
For developing server part I would advice to use `make watch`. Even tho web part can be accessed without installing dependencies, it contains tests which are expected to run in node (see `make test-web`, ./internal/web/package.json). The best way to carry on current work would be to run `npm i` from ./internal/web/. Don't worry, this is optional, but recommended. If you need some functionality - write it from scratch, this will be a great testing field. I value process more than speed of delivery, anything can be DIY.

--------------------------------------------------------------------------------

### ...Bottom line
Did not expect you to see here tbh. Take a breath, relax and enjoy some mindful moments before continuing with your busy-busy-business. Thank you for your time, stranger...
