# Ask m' AI

**Ask m' AI** (ask my ai -> ask m' ai) is a little Desktop-Chat-Application for LargeLanguageModels (e.g. OpenAI's GPT).

Unlike many chat applications out there, this application aims to be scriptable. 
Which means you can call it from you terminal and gave all necessary options as arguments. 
The conversations will also be printed out in the terminal, so you can use it in your scripts. 

## Features

* Support different LLM provider
  * Github Copilot
  * OpenAI
  * AnythingLLM
  * Ollama
  * Mistral
  * Anthropic
* Scriptable
  * All settings can be set via command line arguments
  * The users questions and models answers will be printed out in the terminal

## How to Build this application

1. Install dependencies [see wails documentation](https://wails.io/docs/gettingstarted/installation)
2. Build the application:
    ```sh
    wails build
    ```
