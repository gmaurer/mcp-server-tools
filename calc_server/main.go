package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	fmt.Println("This is an MCP Server")

	s := server.NewMCPServer(
		"Basic Actions Server",
		"0.1.0",
		server.WithToolCapabilities(false),
		server.WithRecovery(),
	)

	calculatorTool := mcp.NewTool("calculator",
		mcp.WithDescription("Performs basic arithmetic operations"),
		mcp.WithString("operation",
			mcp.Required(),
			mcp.Description("Operations include add, subtract, multiply, divide"),
			mcp.Enum("add", "subtract", "multiply", "divide"),
		),
		mcp.WithNumber("x",
			mcp.Required(),
			mcp.Description("First number"),
		),
		mcp.WithNumber("y",
			mcp.Required(),
			mcp.Description("Second number"),
		),
	)

	s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		op, err := request.RequireString("operation")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		x, err := request.RequireFloat("x")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		y, err := request.RequireFloat("y")
		if err != nil {
			return mcp.NewToolResultError(err.Error()), nil
		}

		var result float64
		switch op {
		case "add":
			result = x + y
		case "subtract":
			result = x - y
		case "multiply":
			result = x * y
		case "divide":
			if y == 0 {
				return mcp.NewToolResultError("division by zero is not allowed"), nil
			}
			result = x / y
		default:
			return mcp.NewToolResultError(fmt.Sprintf("unknown operation: %s", op)), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
	})

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Error running MCP server: %v\n", err)
	}

}
