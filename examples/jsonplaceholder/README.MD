## JSONPlaceholder Todos Example

This example demonstrates how the MCP server can interact with
the [JSONPlaceholder](https://jsonplaceholder.typicode.com/) API to retrieve todo items.

### Running the Example

Run the server with:

```bash
api-mcp-server -t sse -c ./jsonplaceholder.tools.json
```

You may optionally set:

* `API_HOST`: defaults to `jsonplaceholder.typicode.com` if not set

### Prerequisite

This MCP server requires a compatible LLM host (such as Claude, GPT, Mistral, etc.) that processes your natural language
prompts and triggers the appropriate API tools.

### Examples

The following examples were run using [mcphost](https://github.com/mark3labs/mcphost) together with Ollama running the
Mistral model.

Example MCP host configuration:

```yaml
mcpServers:
  api-mcp-server:
    url: http://127.0.0.1:13080/sse
```

**Get Todo by ID Prompt:**

> get me the todo with id of 2

**Sample Response:**

```
The todo with id 2 is:

{ "userId": 1, "id": 2, "title": "quis ut nam facilis et officia qui", "completed": false }
```

---

**List All Todos Prompt:**

> show me all the todos

**Sample Response:**

```
This is an array of objects, each representing a task with a unique ID, content (title), and a boolean indicating whether the task is completed or not. The tasks belong to different users, identified by the "userId" property.

Here's an example of one task:

{ "id": 1, "content": "Task 1", "completed": false, "userId": 1 }

To check the completion status of a specific task for a certain user, look for the object with the matching task and user ID.
```