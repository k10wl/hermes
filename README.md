 ```
 _    _ 
| |  | |
| |__| |  ___  _ __  _ __ ___    ___  ___
|  __  | / _ \|  __||  _ \ _ \  / _ \/ __|
| |  | ||  __/| |   | | | | | ||  __/\__ \
|_|  |_| \___||_|   |_| |_| |_| \___||___/
```

# Host-based Extensible Response Management Execution System (HERMES)

## Overview

HERMES, which stands for **Host-based Extensible Response Management Execution System**, is a powerful and customizable platform designed to run AI clients with local control over data and processes. It combines flexibility, performance, and cost-efficiency, offering a range of features that make it ideal for various AI-related tasks.

## Key Features

### HERMES
- **Host-based**: Operates locally on your machine to ensure data security and privacy.
- **Extensible**: Provides a customizable API that covers essential use cases and allows users to update functionalities as needed.
- **Response Management**: Manages and streamlines AI responses with precise control and custom templates.
- **Execution System**: Includes robust capabilities for managing tasks, chaining API calls, and setting up automated workflows.

## Summary
HERMES offers a comprehensive solution for running AI clients with local control over data and efficient management of responses and tasks. It ensures a powerful, flexible, and secure environment for AI interactions.

## Installation
Install app using
```sh
go install github.com/k10wl/hermes@latest
```
Add open AI API key to your environment. Hermes expects HERMES_OPENAI_API_KEY env variable to operate.

## Usage
Usage of hermes (from `hermes -help`):
```
  -h string
        shorthand for "host" (default "127.0.0.1")
  -host string
        Host for web server. Optional, does nothing if "-web" was not provided (default "127.0.0.1")
  -last
        Opens last chat in web. Optional, does nothing if "-web" was not provided
  -m string
        shorthand for "message"
  -message string
        Inline prompt message attached to end of Stdin string, or used as standalone prompt string
  -model string
        ai model name (default "gpt-3.5-turbo")
  -p string
        shorthand for "port" (default "8123")
  -port string
        Port for web server. Optional, does nothing if "-web" was not provided (default "8123")
  -web
        Starts web server
```
