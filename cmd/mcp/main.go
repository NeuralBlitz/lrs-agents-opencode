package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/opencode-ai/opencode/internal/version"
)

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"OpenCode MCP Server",
		version.Version,
		server.WithToolCapabilities(true),
	)

	// Add echo tool
	echoTool := mcp.NewTool("echo",
		mcp.WithDescription("Echo back the input text"),
		mcp.WithString("text",
			mcp.Description("The text to echo back"),
		),
	)
	s.AddTool(echoTool, handleEcho)

	// Add calculator tool
	calcTool := mcp.NewTool("calculator",
		mcp.WithDescription("Perform basic arithmetic operations"),
		mcp.WithString("expression",
			mcp.Description("Mathematical expression to evaluate (e.g., '2 + 2', '10 * 5')"),
		),
	)
	s.AddTool(calcTool, handleCalculator)

	// Add system info tool
	sysInfoTool := mcp.NewTool("system_info",
		mcp.WithDescription("Get system information (OS, architecture, Go version)"),
	)
	s.AddTool(sysInfoTool, handleSystemInfo)

	// Add time tool
	timeTool := mcp.NewTool("get_time",
		mcp.WithDescription("Get the current time in various formats"),
		mcp.WithString("format",
			mcp.Description("Time format: 'unix', 'rfc3339', 'iso8601', or 'human' (default)"),
			mcp.DefaultString("human"),
		),
	)
	s.AddTool(timeTool, handleTime)

	// Add command execution tool (with safety)
	execTool := mcp.NewTool("execute_command",
		mcp.WithDescription("Execute a shell command and return the output"),
		mcp.WithString("command",
			mcp.Description("The command to execute"),
		),
		mcp.WithNumber("timeout",
			mcp.Description("Timeout in seconds (default: 30)"),
			mcp.DefaultNumber(30),
		),
	)
	s.AddTool(execTool, handleExecuteCommand)

	// Add JSON formatter tool
	jsonTool := mcp.NewTool("format_json",
		mcp.WithDescription("Format and validate JSON string"),
		mcp.WithString("json",
			mcp.Description("The JSON string to format"),
		),
		mcp.WithBoolean("indent",
			mcp.Description("Whether to indent the output (default: true)"),
			mcp.DefaultBoolean(true),
		),
	)
	s.AddTool(jsonTool, handleFormatJSON)

	// Add string manipulation tool
	stringTool := mcp.NewTool("string_operations",
		mcp.WithDescription("Perform string operations (upper, lower, reverse, length)"),
		mcp.WithString("text",
			mcp.Description("The text to operate on"),
		),
		mcp.WithString("operation",
			mcp.Description("Operation to perform: 'upper', 'lower', 'reverse', 'length'"),
			mcp.Enum("upper", "lower", "reverse", "length"),
		),
	)
	s.AddTool(stringTool, handleStringOperations)

	// Add a resource template for greeting
	greetingResource := mcp.NewResourceTemplate(
		"greeting://{name}",
		"getGreeting",
		mcp.WithTemplateDescription("Get a personalized greeting"),
	)
	s.AddResourceTemplate(greetingResource, handleGreetingResource)

	// Add a prompt template
	translationPrompt := mcp.NewPrompt(
		"translate",
		mcp.WithPromptDescription("Translate text to another language"),
		mcp.WithString("text",
			mcp.Description("Text to translate"),
		),
		mcp.WithString("target_language",
			mcp.Description("Target language code (e.g., 'es', 'fr', 'de')"),
		),
	)
	s.AddPrompt(translationPrompt, handleTranslationPrompt)

	// Start the server with stdio transport
	if err := server.ServeStdio(s); err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}

// Tool handlers

func handleEcho(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, ok := request.Params.Arguments["text"].(string)
	if !ok {
		return mcp.NewToolResultError(mcp.InvalidParams, "text parameter is required and must be a string"), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Echo: %s", text)), nil
}

func handleCalculator(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	expression, ok := request.Params.Arguments["expression"].(string)
	if !ok {
		return mcp.NewToolResultError(mcp.InvalidParams, "expression parameter is required and must be a string"), nil
	}

	// Simple expression evaluator (for safety, only basic operations)
	expression = strings.TrimSpace(expression)
	
	// Very basic calculator - in production, use a proper expression parser
	var result float64
	var err error
	
	// Try to parse as a simple arithmetic expression
	if strings.Contains(expression, "+") {
		parts := strings.Split(expression, "+")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				result = a + b
			} else {
				err = fmt.Errorf("invalid numbers in expression")
			}
		}
	} else if strings.Contains(expression, "-") {
		parts := strings.Split(expression, "-")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				result = a - b
			} else {
				err = fmt.Errorf("invalid numbers in expression")
			}
		}
	} else if strings.Contains(expression, "*") {
		parts := strings.Split(expression, "*")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				result = a * b
			} else {
				err = fmt.Errorf("invalid numbers in expression")
			}
		}
	} else if strings.Contains(expression, "/") {
		parts := strings.Split(expression, "/")
		if len(parts) == 2 {
			a, err1 := strconv.ParseFloat(strings.TrimSpace(parts[0]), 64)
			b, err2 := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
			if err1 == nil && err2 == nil {
				if b == 0 {
					err = fmt.Errorf("division by zero")
				} else {
					result = a / b
				}
			} else {
				err = fmt.Errorf("invalid numbers in expression")
			}
		}
	} else {
		// Try to parse as a single number
		result, err = strconv.ParseFloat(expression, 64)
	}

	if err != nil {
		return mcp.NewToolResultError(mcp.InvalidParams, fmt.Sprintf("Error evaluating expression: %v", err)), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Result: %g", result)), nil
}

func handleSystemInfo(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	info := map[string]string{
		"os":           runtime.GOOS,
		"architecture": runtime.GOARCH,
		"go_version":   runtime.Version(),
		"num_cpu":      strconv.Itoa(runtime.NumCPU()),
	}

	jsonData, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(mcp.InternalError, fmt.Sprintf("Error marshaling system info: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonData)), nil
}

func handleTime(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	format, _ := request.Params.Arguments["format"].(string)
	if format == "" {
		format = "human"
	}

	now := time.Now()
	var result string

	switch format {
	case "unix":
		result = strconv.FormatInt(now.Unix(), 10)
	case "rfc3339":
		result = now.Format(time.RFC3339)
	case "iso8601":
		result = now.Format(time.RFC3339)
	case "human":
		result = now.Format("2006-01-02 15:04:05 MST")
	default:
		return mcp.NewToolResultError(mcp.InvalidParams, "Invalid format. Use 'unix', 'rfc3339', 'iso8601', or 'human'"), nil
	}

	return mcp.NewToolResultText(result), nil
}

func handleExecuteCommand(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	command, ok := request.Params.Arguments["command"].(string)
	if !ok {
		return mcp.NewToolResultError(mcp.InvalidParams, "command parameter is required and must be a string"), nil
	}

	timeout := 30.0
	if t, ok := request.Params.Arguments["timeout"].(float64); ok {
		timeout = t
	}

	// Create context with timeout
	cmdCtx, cancel := context.WithTimeout(ctx, time.Duration(timeout)*time.Second)
	defer cancel()

	// Determine shell based on OS
	var shell string
	var args []string
	if runtime.GOOS == "windows" {
		shell = "cmd"
		args = []string{"/C", command}
	} else {
		shell = "/bin/sh"
		args = []string{"-c", command}
	}

	cmd := exec.CommandContext(cmdCtx, shell, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return mcp.NewToolResultError(mcp.InternalError, fmt.Sprintf("Command failed: %v\nOutput: %s", err, string(output))), nil
	}

	return mcp.NewToolResultText(string(output)), nil
}

func handleFormatJSON(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	jsonStr, ok := request.Params.Arguments["json"].(string)
	if !ok {
		return mcp.NewToolResultError(mcp.InvalidParams, "json parameter is required and must be a string"), nil
	}

	indent := true
	if i, ok := request.Params.Arguments["indent"].(bool); ok {
		indent = i
	}

	// Parse JSON to validate
	var jsonData interface{}
	if err := json.Unmarshal([]byte(jsonStr), &jsonData); err != nil {
		return mcp.NewToolResultError(mcp.InvalidParams, fmt.Sprintf("Invalid JSON: %v", err)), nil
	}

	// Format JSON
	var formatted []byte
	var err error
	if indent {
		formatted, err = json.MarshalIndent(jsonData, "", "  ")
	} else {
		formatted, err = json.Marshal(jsonData)
	}

	if err != nil {
		return mcp.NewToolResultError(mcp.InternalError, fmt.Sprintf("Error formatting JSON: %v", err)), nil
	}

	return mcp.NewToolResultText(string(formatted)), nil
}

func handleStringOperations(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	text, ok := request.Params.Arguments["text"].(string)
	if !ok {
		return mcp.NewToolResultError(mcp.InvalidParams, "text parameter is required and must be a string"), nil
	}

	operation, ok := request.Params.Arguments["operation"].(string)
	if !ok {
		return mcp.NewToolResultError(mcp.InvalidParams, "operation parameter is required and must be a string"), nil
	}

	var result string
	switch operation {
	case "upper":
		result = strings.ToUpper(text)
	case "lower":
		result = strings.ToLower(text)
	case "reverse":
		runes := []rune(text)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		result = string(runes)
	case "length":
		result = strconv.Itoa(len(text))
	default:
		return mcp.NewToolResultError(mcp.InvalidParams, "Invalid operation. Use 'upper', 'lower', 'reverse', or 'length'"), nil
	}

	return mcp.NewToolResultText(result), nil
}

// Resource handler

func handleGreetingResource(ctx context.Context, request mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	name := request.Params.URI
	// Extract name from URI (greeting://name)
	if strings.HasPrefix(name, "greeting://") {
		name = strings.TrimPrefix(name, "greeting://")
	}

	greeting := fmt.Sprintf("Hello, %s! Welcome to the OpenCode MCP Server.", name)
	return &mcp.ReadResourceResult{
		Contents: []mcp.ResourceContents{
			mcp.NewTextResourceContents(greeting),
		},
	}, nil
}

// Prompt handler

func handleTranslationPrompt(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	text, _ := request.Params.Arguments["text"].(string)
	targetLang, _ := request.Params.Arguments["target_language"].(string)

	if text == "" {
		text = "Hello, world!"
	}
	if targetLang == "" {
		targetLang = "es"
	}

	// This is a simple example - in production, you'd use a real translation service
	translations := map[string]string{
		"es": "Hola, mundo!",
		"fr": "Bonjour, le monde!",
		"de": "Hallo, Welt!",
		"ja": "こんにちは、世界！",
	}

	translated, ok := translations[targetLang]
	if !ok {
		translated = fmt.Sprintf("[Translation to %s: %s]", targetLang, text)
	}

	prompt := fmt.Sprintf("Translate the following text to %s:\n\nOriginal: %s\nTranslation: %s", targetLang, text, translated)

	return &mcp.GetPromptResult{
		Messages: []mcp.PromptMessage{
			mcp.NewPromptMessage(mcp.User, prompt),
		},
	}, nil
}
