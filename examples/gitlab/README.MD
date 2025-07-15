## GitLab Tool Example

This example demonstrates how the MCP server can interact with GitLab’s API to manage project issues.

### Running the Example

Run the server with:

```bash
api-mcp-server -t sse -c ./gitlab.tools.json
```

Make sure to set:

* `GITLAB_TOKEN`: your GitLab personal access token
* `GITLAB_API_HOST`: optional, defaults to `gitlab.com`

### Prerequisite

This MCP server requires a compatible LLM host (such as Claude, GPT, or similar) configured and connected, which
processes your natural language prompts and triggers these API calls.

### Examples

The following examples were run using [mcphost](https://github.com/mark3labs/mcphost) together with Ollama running the
Mistral model.

Example MCP host configuration:

```yaml
mcpServers:
  api-mcp-server:
    url: http://127.0.0.1:13080/sse
```

**List Issues Prompt:**

> list issues for project id AdamShannag/test

**Sample Response:**

```
The following issues are found for the project id AdamShannag/test:

1. ID: 170599967, Title: "add auto retry on failing requests", Description: "External HTTP requests are not retried after failure.", State: "opened", Created At: "2025-07-15T08:40:18Z", Assignee: Adam Shannag  
2. ID: 170119942, Title: "test", Description: "", State: "opened", Created At: "2025-07-07T08:33:57Z", Assignee: None

The first issue was created by Adam Shannag and is currently open. The second issue has no description, is also open, and has no assignee.
```

---

**Create Issue Prompt:**

> create an issue in project AdamShannag/test about a bug report

**Sample Response:**

```
The issue was successfully created in the AdamShannag/test project with an ID of 4. Here's some additional information about the issue:

• Title: Bug report  
• Description: Something is broken  
• Status: opened  
• Created at: 2025-07-15T09:12:58Z  
• Web URL: https://gitlab.com/AdamShannag/test/-/issues/4  
```