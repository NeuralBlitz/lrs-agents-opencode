# OpenCode MCP Server

This is a Model Context Protocol (MCP) server implementation that provides various tools and resources for AI assistants. The server can be used standalone or integrated with OpenCode.

## Overview

The OpenCode MCP Server implements the Model Context Protocol, allowing AI assistants to interact with external tools and services through a standardized interface. The server communicates via stdio and provides a collection of useful tools for development and general-purpose tasks.

## Features

### Tools

The server provides the following tools:

#### 1. **echo**
Echo back the input text.

**Parameters:**
- `text` (string, required): The text to echo back

**Example:**
```json
{
  "name": "echo",
  "arguments": {
    "text": "Hello, World!"
  }
}
```

#### 2. **calculator**
Perform basic arithmetic operations.

**Parameters:**
- `expression` (string, required): Mathematical expression to evaluate (e.g., "2 + 2", "10 * 5")

**Example:**
```json
{
  "name": "calculator",
  "arguments": {
    "expression": "15 * 3 + 7"
  }
}
```

#### 3. **system_info**
Get system information (OS, architecture, Go version, CPU count).

**Parameters:** None

**Example:**
```json
{
  "name": "system_info",
  "arguments": {}
}
```

#### 4. **get_time**
Get the current time in various formats.

**Parameters:**
- `format` (string, optional): Time format - "unix", "rfc3339", "iso8601", or "human" (default: "human")

**Example:**
```json
{
  "name": "get_time",
  "arguments": {
    "format": "rfc3339"
  }
}
```

#### 5. **execute_command**
Execute a shell command and return the output.

**Parameters:**
- `command` (string, required): The command to execute
- `timeout` (number, optional): Timeout in seconds (default: 30)

**Example:**
```json
{
  "name": "execute_command",
  "arguments": {
    "command": "ls -la",
    "timeout": 10
  }
}
```

#### 6. **format_json**
Format and validate JSON string.

**Parameters:**
- `json` (string, required): The JSON string to format
- `indent` (boolean, optional): Whether to indent the output (default: true)

**Example:**
```json
{
  "name": "format_json",
  "arguments": {
    "json": "{\"name\":\"test\",\"value\":123}",
    "indent": true
  }
}
```

#### 7. **string_operations**
Perform string operations (upper, lower, reverse, length).

**Parameters:**
- `text` (string, required): The text to operate on
- `operation` (string, required): Operation to perform - "upper", "lower", "reverse", or "length"

**Example:**
```json
{
  "name": "string_operations",
  "arguments": {
    "text": "Hello, World!",
    "operation": "upper"
  }
}
```

### Resources

#### **greeting://{name}**
Get a personalized greeting.

**URI Format:** `greeting://{name}`

**Example:**
```
greeting://Alice
```

### Prompts

#### **translate**
Translate text to another language.

**Parameters:**
- `text` (string): Text to translate
- `target_language` (string): Target language code (e.g., "es", "fr", "de")

**Supported Languages:**
- `es` - Spanish
- `fr` - French
- `de` - German
- `ja` - Japanese

## Usage

### Standalone Execution

Run the server as a standalone program:

```bash
go run cmd/mcp/main.go
```

Or build and run:

```bash
go build -o mcp-server cmd/mcp/main.go
./mcp-server
```

### As OpenCode Subcommand

Run the server using the OpenCode CLI:

```bash
opencode mcp-server
```

### Integration with OpenCode

To use this MCP server with OpenCode, add it to your OpenCode configuration file (`.opencode.json`):

```json
{
  "mcpServers": {
    "opencode": {
      "type": "stdio",
      "command": "opencode",
      "args": ["mcp-server"]
    }
  }
}
```

Or if you've built a standalone binary:

```json
{
  "mcpServers": {
    "opencode": {
      "type": "stdio",
      "command": "/path/to/mcp-server"
    }
  }
}
```

## Protocol

The server implements the Model Context Protocol (MCP) specification and communicates via stdio using JSON-RPC 2.0. It supports:

- **Tools**: Callable functions that can be invoked by AI assistants
- **Resources**: Readable data sources identified by URIs
- **Prompts**: Template-based prompt generation

## Development

### Adding New Tools

To add a new tool, follow this pattern in `server_impl.go`:

```go
newTool := mcp.NewTool("tool_name",
    mcp.WithDescription("Tool description"),
    mcp.WithString("param1",
        mcp.Description("Parameter description"),
    ),
)
s.AddTool(newTool, handleNewTool)
```

Then implement the handler:

```go
func handleNewTool(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Extract parameters
    param1, ok := request.Params.Arguments["param1"].(string)
    if !ok {
        return mcp.NewToolResultError(mcp.InvalidParams, "param1 is required"), nil
    }
    
    // Process and return result
    result := processTool(param1)
    return mcp.NewToolResultText(result), nil
}
```

### Adding Resources

```go
resource := mcp.NewResourceTemplate(
    "resource://{param}",
    "resourceName",
    mcp.WithTemplateDescription("Resource description"),
)
s.AddResourceTemplate(resource, handleResource)
```

### Adding Prompts

```go
prompt := mcp.NewPrompt(
    "prompt_name",
    mcp.WithPromptDescription("Prompt description"),
    mcp.WithString("text",
        mcp.Description("Text parameter"),
    ),
)
s.AddPrompt(prompt, handlePrompt)
```

## Testing

The server can be tested using any MCP-compatible client. For manual testing, you can use tools like:

- The OpenCode application itself
- Other MCP clients that support stdio transport

## Security Considerations

⚠️ **Important Security Notes:**

1. **Command Execution**: The `execute_command` tool can execute arbitrary shell commands. Use with caution and only in trusted environments.

2. **Input Validation**: All tools validate input parameters, but additional validation may be needed for production use.

3. **Resource Access**: Resources are read-only and don't expose sensitive system information, but be mindful of what data is accessible.

4. **Network Access**: The server doesn't make network requests by default, but tools that do should be carefully reviewed.

## License

This MCP server is part of the OpenCode project and is licensed under the MIT License.

## Contributing

Contributions are welcome! When adding new tools or features:

1. Follow the existing code patterns
2. Add appropriate error handling
3. Include parameter validation
4. Update this README with documentation
5. Test thoroughly before submitting
