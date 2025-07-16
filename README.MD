# API MCP Server

A lightweight and extensible server for defining and executing API tools using
the [MCP protocol](https://modelcontextprotocol.io).
Supports dynamic HTTP requests, typed arguments, and multiple transports (`stdio`, `sse`).

## Installation

You can install the server in one of two ways:

### 1. Install via `go install`

If you have Go installed, run:

```bash
go install github.com/AdamShannag/api-mcp-server/cmd/api-mcp-server@v0.2.0
```

### 2. Download a pre-built binary

Download a binary from the [Releases page](https://github.com/AdamShannag/api-mcp-server/releases).

## CLI Flags

| Flag                | Description                                       | Default         |
|---------------------|---------------------------------------------------|-----------------|
| `--transport`, `-t` | Transport type: `stdio` or `sse`                  | `stdio`         |
| `--config`, `-c`    | Path to JSON tool configuration                   | `./config.json` |
| `--version`, `-v`   | API MCP Server version                            | `-`             |
| `--metrics`, `-m`   | Enable Prometheus metrics and health check server | `-`             |
| `--metrics-port`    | Metrics and health check server port              | `8080`          |

## Environment Variables

| Variable              | Description                                                       | Default     |
|-----------------------|-------------------------------------------------------------------|-------------|
| `API_MCP_HOST`        | SSE host to bind to                                               | `127.0.0.1` |
| `API_MCP_PORT`        | SSE port to bind to                                               | `13080`     |
| `API_MCP_SSE_API_KEY` | Optional Bearer token for auth                                    | *(none)*    |
| `LOG_LEVEL`           | Sets the minimum logging level (`DEBUG`, `INFO`, `WARN`, `ERROR`) | `INFO`      |

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

## Tool Configuration (JSON)

The server accepts a JSON configuration file defining one or more tools. Each tool includes metadata (`name`,
`description`), HTTP request information, and a list of `args` that define the input values required from the LLM.

Each tool config includes:

* A unique `name`
* A `description` shown to the LLM
* A `request` object describing the HTTP call
* A list of `args` to define expected inputs

### Tool Arguments (`args`)

Arguments are inputs collected from the LLM. Each one may be used in one or more parts of the request:

* **Path Parameters** → replaced directly inside the `endpoint`
* **Query Parameters** → automatically added to the URL
* **Body** → inserted as the raw request body string

### Host (`host`)

The `host` is the target API domain (e.g., `gitlab.com`).

```json
"host": "{{env GITLAB_API_HOST:gitlab.com}}"
```

### Method (`method`)

The HTTP method (`GET`, `POST`, `PUT`, etc.).

```json
"method": "POST"
```

### Path Parameters (`pathParams`)

Used to replace variables in the `endpoint`. For example:

```json
"endpoint": "/api/v4/projects/:project_id/pipeline",
"pathParams": ["project_id"]
```

With `project_id = "MyProject/test"`, the final endpoint becomes:

```
/api/v4/projects/MyProject%2Ftest/pipeline
```

### Query Parameters (`queryParams`)

Keys listed here are automatically appended to the endpoint as query string parameters:

```json
"endpoint": "/api/v4/projects/:project_id/pipeline",
"queryParams": ["ref"]
```

With `ref = "main"`, the result is:

```
/api/v4/projects/MyProject%2Ftest/pipeline?ref=main
```

### Request Body (`body`)

The `body` field maps to a single argument name. Its value will be used as the raw request body string.

Example:

```json
"body": "issue_payload"
```

Given this argument:

```json
{
  "name": "issue_payload",
  "type": "string",
  "required": true,
  "description": "A JSON object like {\"title\": \"Bug report\"}"
}
```

The body sent to the server will be the exact string value of `issue_payload`.

### Headers (`headers`)

Define static or dynamic HTTP headers to include in the request, e.g., tokens or content type:

```json
"headers": {
"PRIVATE-TOKEN": "{{env GITLAB_TOKEN:not-set}}",
"Content-Type": "application/json"
}
```

### Secure (`secure`)

If `secure: true`, the request uses `https`. If omitted or `false`, it uses `http`.

### Full Example

```json
{
  "name": "TriggerPipeline",
  "description": "Triggers a pipeline in GitLab for a specific project and branch.",
  "request": {
    "host": "{{env GITLAB_API_HOST:gitlab.com}}",
    "method": "POST",
    "secure": true,
    "endpoint": "/api/v4/projects/:project_id/pipeline",
    "headers": {
      "PRIVATE-TOKEN": "{{env GITLAB_TOKEN:not-set}}"
    },
    "queryParams": [
      "ref"
    ],
    "pathParams": [
      "project_id"
    ]
  },
  "args": [
    {
      "name": "project_id",
      "type": "string",
      "required": true,
      "description": "..."
    },
    {
      "name": "ref",
      "type": "string",
      "required": true,
      "description": "..."
    }
  ]
}
```

If `project_id = "MyProject/test"` and `ref = "main"`, the final HTTP request is:

```
POST https://gitlab.com/api/v4/projects/MyProject%2Ftest/pipeline?ref=main
Headers:
  PRIVATE-TOKEN: <value from GITLAB_TOKEN>
```

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

