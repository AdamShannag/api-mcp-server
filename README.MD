# API MCP Server

A lightweight, extensible MCP server for defining and executing API tools using
the [MCP protocol](https://modelcontextprotocol.io). Supports dynamic HTTP requests, typed arguments, and multiple
transports (`stdio`, `sse`).

## Installation

````bash
go install github.com/AdamShannag/api-mcp-server@v0.0.1
````

### Flags

| Flag                | Description                      | Default |
|---------------------|----------------------------------|---------|
| `--transport`, `-t` | Transport type: `stdio` or `sse` | `stdio` |
| `--config`, `-c`    | Path to JSON tools config        | `./`    |

### Environment Variables (SSE only)

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
      "host": "jsonplaceholder.typicode.com",
      "endpoint": "/todos/:id",
      "method": "GET",
      "secure": true,
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

Place this file somewhere accessible and provide its path via the `--config` flag when running the server:

```bash
api-mcp-server -t sse -c ./demo.tools.json
```

This loads the example tools and enables interaction with the JSONPlaceholder API through your MCP server.

## Authentication

For SSE transport, API key authentication is supported via the environment variable `API_MCP_SSE_API_KEY`. If set,
incoming requests must provide a matching Bearer token in the `Authorization` header.

Example:

```bash
export API_MCP_SSE_API_KEY="your-secret-key"
api-mcp-server --transport sse --config ./demo.tools.json
```

## Development

* The server code is organized under the `internal/mcp` package.
* Tool management and execution are handled by the `tool` and `request` packages.
* Argument resolution and validation are in the `resolver` package.
* Authentication middleware is in the `auth` package.

## License

[MIT](LICENSE)

## Contribution

If you have any questions or want to contribute, feel free to open an issue or PR.

