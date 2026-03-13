package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mark3labs/mcp-go/server"

	"github.com/hmsoft0815/mlc-markitdown/internal/presentation/handler"
	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
	"github.com/hmsoft0815/mlcartifact/client"
)

const (
	name = "mlc-markitdown"
)

var version = "v0.1.0" // Default version, can be overridden by ldflags

func main() {
	threshold := flag.Int("threshold", 10000, "Character threshold for auto-artifact storage")
	showVersion := flag.Bool("version", false, "Show version and exit")
	llmProvider := flag.String("llm-provider", "", "Default LLM provider (openai, ollama)")
	llmUrl := flag.String("llm-url", "", "Default LLM API URL")
	llmModel := flag.String("llm-model", "", "Default LLM model name")
	llmKey := flag.String("llm-key", "", "Default LLM API key")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", name, version)
		os.Exit(0)
	}

	// 1. Initialize Artifact Client
	artifactCli, err := client.NewClient()
	if err != nil {
		log.Fatalf("Failed to connect to artifact server: %v", err)
	}
	defer artifactCli.Close()

	// 2. Initialize UseCase
	convertUC := usecase.NewConvertUseCase(artifactCli, *threshold, usecase.LlmConfig{
		Provider: *llmProvider,
		URL:      *llmUrl,
		Model:    *llmModel,
		AuthKey:  *llmKey,
	})

	// 3. Initialize MCP Server
	mcpServer := server.NewMCPServer(name, version)

	// 4. Register Tools
	convertHandler := handler.NewConvertHandler(convertUC, mcpServer)
	convertArtifactHandler := handler.NewConvertArtifactHandler(convertUC, artifactCli, mcpServer)
	quickInspectHandler := handler.NewQuickInspectHandler()
	promptHandler := handler.NewPromptHandler()
	llmUtilsHandler := handler.NewLlmUtilsHandler(convertUC)

	mcpServer.AddTool(convertHandler.GetTool(), convertHandler.Handle)
	mcpServer.AddTool(convertArtifactHandler.GetTool(), convertArtifactHandler.Handle)
	mcpServer.AddTool(quickInspectHandler.GetTool(), quickInspectHandler.Handle)
	mcpServer.AddTool(llmUtilsHandler.GetTool(), llmUtilsHandler.Handle)
	mcpServer.AddPrompt(promptHandler.GetPrompt(), promptHandler.Handle)

	// 5. Start Server
	fmt.Fprintf(os.Stderr, "MLC MarkItDown MCP Server starting (version %s)\n", version)
	if err := server.ServeStdio(mcpServer); err != nil {
		log.Fatalf("MCP server error: %v", err)
	}
}
