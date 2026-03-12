package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/hmsoft0815/mlc-markitdown/internal/presentation/handler"
	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
	"github.com/hmsoft0815/mlcartifact"
)

const (
	name    = "mlc-markitdown"
	version = "1.0.0"
)

func main() {
	threshold := flag.Int("threshold", 10000, "Character threshold for auto-artifact storage")
	flag.Parse()

	// 1. Initialize Artifact Client
	artifactCli, err := mlcartifact.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to artifact server: %v", err)
	}
	defer artifactCli.Close()

	// 2. Initialize UseCase
	convertUC := usecase.NewConvertUseCase(artifactCli, *threshold)

	// 3. Initialize MCP Server
	mcpServer := server.NewMCPServer(name, version)

	// 4. Register Tools
	convertHandler := handler.NewConvertHandler(convertUC)
	convertArtifactHandler := handler.NewConvertArtifactHandler(convertUC, artifactCli)
	quickInspectHandler := handler.NewQuickInspectHandler()

	mcpServer.AddTool(convertHandler.GetTool(), convertHandler.Handle)
	mcpServer.AddTool(convertArtifactHandler.GetTool(), convertArtifactHandler.Handle)
	mcpServer.AddTool(quickInspectHandler.GetTool(), quickInspectHandler.Handle)

	// 5. Start Server
	fmt.Fprintf(os.Stderr, "MLC MarkItDown MCP Server starting (version %s)\n", version)
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}
