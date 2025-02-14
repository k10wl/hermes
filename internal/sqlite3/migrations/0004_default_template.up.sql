INSERT OR REPLACE INTO
    templates (name, content)
VALUES (
    'hermes-help',
    '--{{define "hermes-help"}}' || CHAR(10) ||
    '<documentation>' || CHAR(10) ||
    '# Host-based Extensible Response Management Execution System (HERMES)' || CHAR(10) ||
    CHAR(10) ||
    '## Overview' || CHAR(10) ||
    'Hermes provides unified local point entry to popular AI providers.' || CHAR(10) ||
    'You can call completion using desired model in following ways:' || CHAR(10) ||
    '- CLI, app is build to use stdin, stdout, stderr. Commands can be chained and work with other tools from terminal. Use `hermes --help` for more info regarding CLI usage;' || CHAR(10) ||
    '- Web, app ships HTTP server to use all features from browser GUI. Use `hermes serve --help` for more info;' || CHAR(10) ||
    CHAR(10) ||
    '## Name derives from following valuable key features:' || CHAR(10) ||
    '- **Host-based**: Stores data locally, allows single entry point to multiple third party providers.' || CHAR(10) ||
    '- **Extensible**: Provides a customizable API that covers essential use cases and allows to update functionalities as needed.' || CHAR(10) ||
    '- **Response Management**: Manages and streamlines AI responses with intuitive UX.' || CHAR(10) ||
    '- **Execution System**: Includes robust capabilities for templating prompts, chaining API calls, and setting up automated workflows.' || CHAR(10) ||
    CHAR(10) ||
    '--------------------------------------------------------------------------------' || CHAR(10) ||
    CHAR(10) ||
    '## Instalation' || CHAR(10) ||
    CHAR(10) ||
    '### Using golang (requires go 1.23.1+)' || CHAR(10) ||
    'You can add it to you golang managed library using' || CHAR(10) ||
    '```bash' || CHAR(10) ||
    'make install' || CHAR(10) ||
    '```' || CHAR(10) ||
    'Or build binaries into ./bin/ directory by yourself' || CHAR(10) ||
    '```bash' || CHAR(10) ||
    'make build' || CHAR(10) ||
    '```' || CHAR(10) ||
    CHAR(10) ||
    '### Environmental variables:' || CHAR(10) ||
    '```bash' || CHAR(10) ||
    '# Models' || CHAR(10) ||
    'HERMES_OPENAI_API_KEY=    # key to OpenAI API, optional' || CHAR(10) ||
    'HERMES_ANTHROPIC_API_KEY= # key to Anthropic API, optional' || CHAR(10) ||
    CHAR(10) ||
    '# Files' || CHAR(10) ||
    'HERMES_DB_DNS=            # sqlite3 dns entry to persist chats, templates, messages' || CHAR(10) ||
    '                          # optional, defaults to /hermes/main.db in your config dir' || CHAR(10) ||
    '                          # rel: https://pkg.go.dev/os#UserConfigDir' || CHAR(10) ||
    CHAR(10) ||
    '# DEBUG' || CHAR(10) ||
    'HERMES_MOCK_COMPLETION=   # debug mode, chats will return inputted message' || CHAR(10) ||
    '                          # for development `make watch` this value to true' || CHAR(10) ||
    '                          # not to be use in any real scenario' || CHAR(10) ||
    '```' || CHAR(10) ||
    CHAR(10) ||
    '--------------------------------------------------------------------------------' || CHAR(10) ||
    CHAR(10) ||
    '## Templates' || CHAR(10) ||
    'One of core features of Hermes are templates. Golang provides flexible and extensible way to manage text contents. It is used as templating engine. Having such powerful utility allows to create reusable prompts with custom instructions. You can ask questions and lead users to answers based on documentation. Tool is flexible, take a minute to implement smart solutions' || CHAR(10) ||
    CHAR(10) ||
    'After executing this CLI command:' || CHAR(10) ||
    '```bash' || CHAR(10) ||
    'hermes template upsert --content' || CHAR(10) ||
    '    ''−−{{define "short"}}−−{{.}} (Short answer please)−−{{end}}''' || CHAR(10) ||
    '```' || CHAR(10) ||
    'You will have stored your first template named "short". It can be now used in chats!' || CHAR(10) ||
    '```bash' || CHAR(10) ||
    'hermes chat --content ''what is the mass of the Sun?'' --template short' || CHAR(10) ||
    '#                     name inherited from template definition   ^^^^^' || CHAR(10) ||
    '```' || CHAR(10) ||
    'AI will receive following - `what is the mass of the Sun? (Short answer please)`.' || CHAR(10) ||
    'And now we have reusable template that will make response concise. You can you any prompts with any instructions, AI is really pleasant to use when there is no need to repeat yourself' || CHAR(10) ||
    CHAR(10) ||
    CHAR(10) ||
    'Hermes ships with copy of this README as one of templates, you can always ask questions related this application using `hermes-help`' || CHAR(10) ||
    '```bash' || CHAR(10) ||
    'hermes chat --content ''how to I change database?'' --template hermes-help' || CHAR(10) ||
    '```' || CHAR(10) ||
    CHAR(10) ||
    '--------------------------------------------------------------------------------' || CHAR(10) ||
    CHAR(10) ||
    '## Development' || CHAR(10) ||
    CHAR(10) ||
    '### Internals' || CHAR(10) ||
    'Application has decent code coverage over most valuable and sensitive parts, you can always run `make test` to make sure that defined parts of application work as expected. You can skip JS testing part if you run `make pre-test && make test-app`.' || CHAR(10) ||
    CHAR(10) ||
    '### Workflows' || CHAR(10) ||
    'CI is configured, pull requests will run tests, refer to github actions result as to source of truth.' || CHAR(10) ||
    CHAR(10) ||
    '### Server (ignore this if you are not going to touch js parts)' || CHAR(10) ||
    'For developing server part I would advice to use `make watch`. Even tho web part can be accessed without installing dependencies, it contains tests which are expected to run in node (see `make test-web`, ./internal/web/package.json). The best way to carry on current work would be to run `npm i` from ./internal/web/. Don''t worry, this is optional, but recommended. If you need some functionality - write it from scratch, this will be a great testing field. I value process more than speed of delivery, anything can be DIY.' || CHAR(10) ||
    CHAR(10) ||
    '--------------------------------------------------------------------------------' || CHAR(10) ||
    CHAR(10) ||
    '### ...Bottom line' || CHAR(10) ||
    'Did not expect you to see here tbh. Take a breath, relax and enjoy some mindful moments before continuing with your busy-busy-business. Thank you for your time, stranger...' || CHAR(10) ||
    '</documentation>' || CHAR(10) ||
    CHAR(10) ||
    CHAR(10) ||
    'You will be asked a question or completion. Make you answer short, do not state whole documentation. Do not create step by step instructions if not asked. Answer clear, and stick to documentation. This documentation describes one use case, extend on it. Do not hallucinate new documentation. Do not cite documentation if answer can be flexible and based on user input. Feel free to improvise your own templates' || CHAR(10) ||
    CHAR(10) ||
    CHAR(10) ||
    'Based on <documentation> answer this: --{{.}}' || CHAR(10) ||
    '--{{end}}'
); 
