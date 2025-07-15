# API MCP Server

A lightweight, extensible MCP server for defining and executing API tools using
the [MCP protocol](https://modelcontextprotocol.io). Supports dynamic HTTP requests, typed arguments, and multiple
transports (`stdio`, `sse`).

## Installation

You can install the server in one of two ways:

### 1. Install via `go install`

If you have Go installed, run:

```bash
go install github.com/AdamShannag/api-mcp-server/cmd/api-mcp-server@v0.0.1
```

### 2. Download pre-built binary

You can also download a pre-built binary from
the [Releases page](https://github.com/AdamShannag/api-mcp-server/releases).

## Flags

| Flag                | Description                      | Default |
|---------------------|----------------------------------|---------|
| `--transport`, `-t` | Transport type: `stdio` or `sse` | `stdio` |
| `--config`, `-c`    | Path to JSON tools config        | `./`    |

## Environment Variables (SSE only)

| Variable              | Description                    | Default     |
|-----------------------|--------------------------------|-------------|
| `API_MCP_HOST`        | SSE host to bind               | `127.0.0.1` |
| `API_MCP_PORT`        | SSE port to bind               | `13080`     |
| `API_MCP_SSE_API_KEY` | Optional bearer token for auth | *(none)*    |

## Demo Tools Configuration

A sample tools configuration file `demo.tools.json` is provided to help you get started quickly. It defines example
tools that interact with the [JSONPlaceholder](https://jsonplaceholder.typicode.com/) public API:

```json
[
  {
    "name": "ListTodos",
    "description": "Retrieves a list of all todo items. The response is a JSON array where each item contains details such as ID, title, completion status, and associated user ID.",
    "request": {
      "host": "jsonplaceholder.typicode.com",
      "endpoint": "/todos",
      "method": "GET",
      "secure": true
    }
  },
  {
    "name": "GetTodo",
    "description": "Retrieves a specific todo item by its ID",
    "request": {
      "host": "{{env API_HOST:jsonplaceholder.typicode.com}}",
      "endpoint": "/todos/:id",
      "method": "GET",
      "secure": true,
      "headers": {
        "API_KEY": "Authorization {{env API_KEY:not-set}}"
      },
      "pathParams": [
        "id"
      ]
    },
    "args": [
      {
        "name": "id",
        "description": "The unique identifier of the todo item to retrieve.",
        "required": true,
        "type": "int"
      }
    ]
  }
]
```

> **Note:** Setting `"secure": true` in a tool's request configuration ensures the request is made over **HTTPS**
> instead of **HTTP**.

Place this file somewhere accessible and provide its path via the `--config` flag when running the server:

```bash
api-mcp-server -t sse -c ./demo.tools.json
```

This loads the example tools and enables interaction with the JSONPlaceholder API through your MCP server.

### Environment Variable Placeholders

The tool configuration supports environment variable interpolation using the following syntax:

```text
{{env VAR_NAME:default_value}}
```

At runtime, these placeholders are resolved to the corresponding environment variable value. If the environment variable
is not set, the `default_value` will be used instead.

#### Examples:

* `"host": "{{env API_HOST:jsonplaceholder.typicode.com}}"`
  → Uses the `API_HOST` environment variable if set, otherwise defaults to `jsonplaceholder.typicode.com`.

* `"Authorization": "Bearer {{env API_KEY:not-set}}"`
  → Substitutes the value of `API_KEY`, or uses `"not-set"` if the variable is missing.

Make sure to export any required variables before running the server:

```bash
export API_HOST=myapi.example.com
export API_KEY=supersecret123
```

## Authentication

For SSE transport, API key authentication is supported via the environment variable `API_MCP_SSE_API_KEY`. If set,
incoming requests must provide a matching Bearer token in the `Authorization` header.

Example:

```bash
export API_MCP_SSE_API_KEY="your-secret-key"
api-mcp-server --transport sse --config ./demo.tools.json
```

## License

[MIT](LICENSE)

## Contribution

If you have any questions or want to contribute, feel free to open an issue or PR.

