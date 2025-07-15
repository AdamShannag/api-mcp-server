# API MCP Server

A lightweight and extensible server for defining and executing API tools using
the [MCP protocol](https://modelcontextprotocol.io).
Supports dynamic HTTP requests, typed arguments, and multiple transports (`stdio`, `sse`).

## Installation

You can install the server in one of two ways:

### 1. Install via `go install`

If you have Go installed, run:

```bash
go install github.com/AdamShannag/api-mcp-server/cmd/api-mcp-server@v0.1.0
```

### 2. Download a pre-built binary

Download a binary from the [Releases page](https://github.com/AdamShannag/api-mcp-server/releases).

## CLI Flags

| Flag                | Description                      | Default               |
|---------------------|----------------------------------|-----------------------|
| `--transport`, `-t` | Transport type: `stdio` or `sse` | `stdio`               |
| `--config`, `-c`    | Path to JSON tool configuration  | `./tools.config.json` |

## Environment Variables (for SSE transport)

| Variable              | Description                    | Default     |
|-----------------------|--------------------------------|-------------|
| `API_MCP_HOST`        | SSE host to bind to            | `127.0.0.1` |
| `API_MCP_PORT`        | SSE port to bind to            | `13080`     |
| `API_MCP_SSE_API_KEY` | Optional Bearer token for auth | *(none)*    |

## Tool Config Placeholders

The tool configuration file supports environment variable substitution using the syntax:

```text
{{env VAR_NAME:default_value}}
```

At runtime, placeholders are replaced with the corresponding environment variable. If the variable is not set, the
`default_value` is used instead.

### Examples

```json
"host": "{{env API_HOST:jsonplaceholder.typicode.com}}"
```

Uses the value of `API_HOST` if set, otherwise falls back to the default.

```json
"Authorization": "Bearer {{env API_KEY:not-set}}"
```

Uses `API_KEY` or `"not-set"` if not defined.

Before running the server, set any required environment variables:

```bash
export API_HOST=myapi.example.com
export API_KEY=supersecret123
```

## Authentication

When using the SSE transport, you can optionally secure the endpoint using a Bearer token by setting the
`API_MCP_SSE_API_KEY` environment variable.

Example:

```bash
export API_MCP_SSE_API_KEY="your-secret-key"
api-mcp-server --transport sse --config ./demo.tools.json
```

Incoming SSE connections must then provide the matching token in the `Authorization` header.

## Examples

This repository includes several example tool configurations to demonstrate different use cases. These are not
ready-to-use APIs but serve as references:

* **GitLab Issues**
  Tools for listing and creating issues in a GitLab project.
  See [examples/gitlab](./examples/gitlab/README.MD)

* **JSONPlaceholder Todos**
  Tools interacting with the [JSONPlaceholder](https://jsonplaceholder.typicode.com/) public API.
  See [examples/jsonplaceholder](./examples/jsonplaceholder/README.MD)

## License

[MIT](LICENSE)

## Contribution

If you have any questions or want to contribute, feel free to open an issue or PR.

