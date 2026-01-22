package mcp

import (
	"fmt"

	"github.com/spf13/cobra"
)

// Note: This command is a wrapper. The actual server implementation
// is in main.go and can be run standalone with: go run cmd/mcp/main.go
var ServerCmd = &cobra.Command{
	Use:   "mcp-server",
	Short: "Start the OpenCode MCP server",
	Long: `Start the Model Context Protocol (MCP) server that provides tools and resources
for AI assistants. The server communicates via stdio and can be used by MCP clients
like OpenCode or other compatible applications.

The server provides various tools including:
- echo: Echo back input text
- calculator: Perform basic arithmetic operations
- system_info: Get system information
- get_time: Get current time in various formats
- execute_command: Execute shell commands
- format_json: Format and validate JSON
- string_operations: Perform string manipulations

It also provides resources and prompts for extended functionality.`,
	Example: `
  # Start the MCP server (communicates via stdio)
  opencode mcp-server

  # The server is typically used by configuring it in OpenCode's config:
  # {
  #   "mcpServers": {
  #     "opencode": {
  #       "type": "stdio",
  #       "command": "opencode",
  #       "args": ["mcp-server"]
  #     }
  #   }
  # }`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Import and run the server
		// We need to call the runServer function from main.go
		// Since we can't directly import main package, we'll need to refactor
		// For now, this is a placeholder - users should run: go run cmd/mcp/main.go
		return fmt.Errorf("MCP server subcommand is not yet integrated. Please run: go run cmd/mcp/main.go")
	},
}
