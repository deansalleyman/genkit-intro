# Genkit Intro – Recipe Generator

A small Go app that uses [Genkit Go](https://firebase.google.com/docs/genkit-go) to generate recipes from an ingredient and optional dietary restrictions. It demonstrates flows, structured output, and the Genkit Developer UI.

## What this app does

- **Flow**: `recipeGeneratorFlow` takes `ingredient` and `dietaryRestrictions` and returns a structured `Recipe` (title, description, prep/cook time, servings, ingredients, instructions, tips).
- **Structured output**: The flow uses `genkit.GenerateData[Recipe]` so the model returns typed JSON matching the `Recipe` struct.
- **HTTP API**: The flow is exposed as `POST /recipeGeneratorFlow`. Request body can be either `{"ingredient": "...", "dietaryRestrictions": "..."}` or Genkit’s `{"data": { ... }}` format.
- **Prompts**: Prompts can be loaded from the `prompts/` directory (e.g. `create_recipe.prompt`) via `genkit.WithPromptDir("./prompts")`.

## Why Genkit Go?

[Genkit Go](https://firebase.google.com/docs/genkit-go) is a Go-native, open-source framework for building AI apps. As described in [Mastering Genkit: Go Edition – Introduction](https://mastering-genkit.github.io/mastering-genkit-go/chapters/01-introduction-genkit-go.html), it brings several benefits for Go developers:

- **Go-native and type-safe**  
  Clear, readable APIs and explicit error handling. Inputs and outputs are Go types with compile-time guarantees, so you avoid runtime surprises from untyped JSON.

- **Unified model interface**  
  One API for multiple providers (Google AI, Vertex AI, OpenAI, Anthropic, etc.). Switching or adding providers doesn’t require rewriting your flows.

- **Production-oriented**  
  Built with deployment in mind: tracing, observability, and error handling. You get a single binary, minimal dependencies, and straightforward deployment to Cloud Run, Kubernetes, or any Go-friendly host.

- **Developer experience**  
  The Genkit Developer UI (run with `genkit start`) gives a local playground to test flows, inspect traces, and iterate on prompts without leaving your app.

- **Fits existing Go systems**  
  Genkit Go plugs into existing Go backends and uses the same concurrency and deployment patterns you already use.

For a deeper dive into architecture, flows, and production patterns, see [Mastering Genkit: Go Edition](https://mastering-genkit.github.io/mastering-genkit-go).

## Prerequisites

- **Go 1.25+**
- **Genkit CLI** (for the Developer UI):
  ```bash
  curl -sL cli.genkit.dev | bash
  ```
  Ensure the install directory (e.g. `~/.local/bin`) is on your `PATH`.
- **Google AI (Gemini) API key** for `GEMINI_API_KEY` (e.g. from [Google AI Studio](https://aistudio.google.com/apikey)).

## Setup

1. **Clone and install dependencies**
   ```bash
   cd genkit-intro
   go mod download
   ```

2. **Configure environment**
   - Copy or create a `.env` in the project root with:
     ```bash
     GEMINI_API_KEY=your_api_key_here
     ```
   - The app loads `.env` at startup (via `godotenv`).

## Running the app

**With Genkit Developer UI (recommended for development):**

```bash
genkit start -- go run .
```

Then open:

- **Genkit Developer UI**: http://localhost:4001 (or the URL printed in the terminal)
- **Telemetry API**: http://localhost:4033
- **App server**: http://localhost:3400

**Without the Dev UI (app only):**

```bash
go run .
```

The app listens on `http://localhost:3400` and runs a one-off recipe generation at startup, then serves the flow over HTTP.

## Using the API

**Endpoint:** `POST http://localhost:3400/recipeGeneratorFlow`

**Request body (either format):**

```json
{
  "ingredient": "pasta",
  "dietaryRestrictions": "dairy-free"
}
```

or Genkit’s wrapped form:

```json
{
  "data": {
    "ingredient": "pasta",
    "dietaryRestrictions": "dairy-free"
  }
}
```

**Example with curl:**

```bash
curl -X POST http://localhost:3400/recipeGeneratorFlow \
  -H "Content-Type: application/json" \
  -d '{"ingredient": "avocado", "dietaryRestrictions": "vegetarian"}'
```

## Project layout

- `main.go` – Genkit init, flow definition, HTTP handler, and server.
- `prompts/` – Dotprompt files (e.g. `create_recipe.prompt`) loaded with `WithPromptDir`.
- `.env` – Local env vars (e.g. `GEMINI_API_KEY`); not committed.

## Resources

- [Genkit Go – Get started](https://firebase.google.com/docs/genkit-go/get-started-go)
- [Mastering Genkit: Go Edition – Introduction](https://mastering-genkit.github.io/mastering-genkit-go/chapters/01-introduction-genkit-go.html)
- [Genkit Developer Tools](https://firebase.google.com/docs/genkit/devtools)
