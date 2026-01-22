package main

import (
	"fmt"
	"os"

	"github.com/opencode-ai/opencode/cmd/mcp"
)

func main() {
	if err := mcp.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
