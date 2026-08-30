package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
	"github.com/hmsoft0815/mlcartifact/client"
)

type ConvertArtifactHandler struct {
	useCase     *usecase.ConvertUseCase
	artifactCli *client.Client
	server      *server.MCPServer
}

func NewConvertArtifactHandler(useCase *usecase.ConvertUseCase, artifactCli *client.Client, srv *server.MCPServer) *ConvertArtifactHandler {
	return &ConvertArtifactHandler{
		useCase:     useCase,
		artifactCli: artifactCli,
		server:      srv,
	}
}

func (h *ConvertArtifactHandler) GetTool() mcp.Tool {
	return mcp.NewTool(
		"markitdown__convert_artifact__mlc",
		mcp.WithDescription("Converts a document already stored in the artifact storage to Markdown. Supports vision if enable_vision is true."),
		mcp.WithString("artifactId", mcp.Description("The ID of the source artifact to convert"), mcp.Required()),
		mcp.WithString("output_filename", mcp.Description("Optional name for the resulting Markdown artifact (e.g. report.md).")),
		mcp.WithBoolean("enable_vision", mcp.Description("If true, use an LLM for vision/audio descriptions.")),
		mcp.WithString("llm_provider", mcp.Description("LLM provider to use: 'openai' or 'ollama'")),
		mcp.WithString("openai_key", mcp.Description("OpenAI API Key (if provider is openai)")),
		mcp.WithString("openai_model", mcp.Description("OpenAI model to use (default: gpt-4o)")),
		mcp.WithString("openai_url", mcp.Description("Custom OpenAI compatible URL")),
		mcp.WithString("ollama_model", mcp.Description("The model to use with Ollama")),
		mcp.WithString("ollama_url", mcp.Description("The URL of the Ollama server")),
		mcp.WithOutputSchema[ConvertResponse](),
	)
}

func (h *ConvertArtifactHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	artifactID := mcp.ParseString(request, "artifactId", "")
	if artifactID == "" {
		return mcp.NewToolResultError("artifactId is required"), nil
	}

	outputFilename := mcp.ParseString(request, "output_filename", "")
	enableVision := mcp.ParseBoolean(request, "enable_vision", false)
	var progressToken any
	if request.Params.Meta != nil {
		progressToken = request.Params.Meta.ProgressToken
	}

	var llmModel, openaiKey, llmBaseUrl string
	if enableVision {
		defProvider, defModel, defUrl, defKey := h.useCase.GetLlmDefaults()

		llmProvider := mcp.ParseString(request, "llm_provider", defProvider)

		if llmProvider == "ollama" {
			llmModel = mcp.ParseString(request, "ollama_model", defModel)
			llmBaseUrl = mcp.ParseString(request, "ollama_url", defUrl)
			if llmBaseUrl == "" {
				llmBaseUrl = "http://localhost:11434/v1"
			}
			openaiKey = "ollama" // Dummy key for shim
		} else {
			openaiKey = mcp.ParseString(request, "openai_key", defKey)
			if openaiKey == "" {
				return mcp.NewToolResultError("OpenAI API key is required but not provided or set in environment"), nil
			}
			llmModel = mcp.ParseString(request, "openai_model", defModel)
			if llmModel == "" {
				llmModel = "gpt-4o"
			}
			llmBaseUrl = mcp.ParseString(request, "openai_url", defUrl)
		}
	}

	// 0. Define progress monitor
	progress := func(percent int, status string) {
		if progressToken != "" && h.server != nil {
			_ = h.server.SendNotificationToClient(ctx, "notifications/progress", map[string]interface{}{
				"progress":      float64(percent),
				"total":         100.0,
				"progressToken": progressToken,
				"message":       status,
			})
		}
	}

	// 1. Read source artifact
	progress(5, "Reading source artifact...")
	res, err := h.artifactCli.Read(ctx, artifactID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to read source artifact", err), nil
	}

	if outputFilename == "" && res.Filename != "" {
		outputFilename = res.Filename + ".md"
	}

	// 2. Write to temp file for MarkItDown (since it needs a file path)
	progress(20, "Preparing temp file...")
	tmpFile := fmt.Sprintf("/tmp/markitdown_%s", artifactID)
	err = h.useCase.WriteTempFile(tmpFile, res.Content)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to prepare temp file", err), nil
	}

	// 3. Convert
	req := usecase.ConvertRequest{
		URI:            tmpFile,
		ForceArtifact:  true,
		OutputFilename: outputFilename,
		Progress:       progress,
		LlmModel:       llmModel,
		OpenaiKey:      openaiKey,
		LlmBaseUrl:     llmBaseUrl,
	}
	content, newArtifact, err := h.useCase.Convert(ctx, req)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Conversion failed", err), nil
	}

	// 4. Return result
	result := &mcp.CallToolResult{
		Content: []mcp.Content{},
	}

	if newArtifact != nil {
		notice := fmt.Sprintf("## Artifact Converted\n\n**Notice**: The complete file is available in the artifact server under id = %s", newArtifact.Id)
		result.Content = append(result.Content, mcp.NewTextContent(notice))

		// Structured JSON
		meta := map[string]interface{}{
			"artifact": map[string]interface{}{
				"id":         newArtifact.Id,
				"filename":   newArtifact.Filename,
				"source":     "mlc-markitdown",
				"expires_at": newArtifact.ExpiresAt,
			},
		}
		jsonBytes, _ := json.MarshalIndent(meta, "", "  ")
		result.Content = append(result.Content, mcp.TextContent{
			Type: "text",
			Text: string(jsonBytes),
		})
	} else {
		result.Content = append(result.Content, mcp.NewTextContent(content))
	}

	return result, nil
}
