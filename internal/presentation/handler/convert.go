package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"

	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
)

type ConvertHandler struct {
	useCase *usecase.ConvertUseCase
	server  *server.MCPServer
}

type ArtifactInfo struct {
	ID        string `json:"id" jsonschema:"description=The unique ID of the artifact"`
	Filename  string `json:"filename" jsonschema:"description=Original filename"`
	Source    string `json:"source" jsonschema:"description=Source system"`
	ExpiresAt string `json:"expires_at" jsonschema:"description=Expiration timestamp"`
}

type ConvertResponse struct {
	Markdown string        `json:"markdown" jsonschema:"description=The converted markdown content (full or preview)"`
	Artifact *ArtifactInfo `json:"artifact,omitempty" jsonschema:"description=Information about the saved artifact if applicable"`
	IsFull   bool          `json:"is_full" jsonschema:"description=True if the markdown field contains the full document"`
}

func NewConvertHandler(useCase *usecase.ConvertUseCase, srv *server.MCPServer) *ConvertHandler {
	return &ConvertHandler{
		useCase: useCase,
		server:  srv,
	}
}

func (h *ConvertHandler) GetTool() mcp.Tool {
	return mcp.NewTool(
		"markitdown__convert__mlc",
		mcp.WithDescription("Converts a file or URL to Markdown. Smart auto-archiving is applied for large outputs. Supports vision/audio descriptions if enable_vision is true."),
		mcp.WithString("uri", mcp.Description("The source path or URL to convert"), mcp.Required()),
		mcp.WithBoolean("force_artifact", mcp.Description("If true, always save as artifact storage and return a notice.")),
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

func (h *ConvertHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uri := mcp.ParseString(request, "uri", "")
	if uri == "" {
		return mcp.NewToolResultError("uri is required"), nil
	}

	forceArtifact := mcp.ParseBoolean(request, "force_artifact", false)
	enableVision := mcp.ParseBoolean(request, "enable_vision", false)
	progressToken := request.Params.Meta.ProgressToken

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

	// 1. Define progress monitor
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

	// 2. Call Usecase
	content, artifact, err := h.useCase.Convert(ctx, uri, forceArtifact, progress, llmModel, openaiKey, llmBaseUrl)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Conversion failed", err), nil
	}

	// 3. Prepare response
	res := &mcp.CallToolResult{
		Content: []mcp.Content{},
	}

	if artifact != nil || forceArtifact {
		// Artifact saved - provide preview + ID notice
		preview := content
		if len(preview) > 500 {
			preview = preview[:500] + "..."
		}

		notice := fmt.Sprintf("## Document Converted\n\nPreview:\n%s\n\n**Notice**: The complete file is available in the artifact server under id = %s", preview, artifact.Id)
		res.Content = append(res.Content, mcp.NewTextContent(notice))

		// Structured JSON
		meta := map[string]interface{}{
			"artifact": map[string]interface{}{
				"id":         artifact.Id,
				"filename":   artifact.Filename,
				"source":     "mlc-markitdown",
				"expires_at": artifact.ExpiresAt,
			},
		}
		jsonBytes, _ := json.MarshalIndent(meta, "", "  ")
		res.Content = append(res.Content, mcp.TextContent{
			Type: "text",
			Text: string(jsonBytes),
		})
	} else {
		// Tiny document - return full content
		res.Content = append(res.Content, mcp.NewTextContent(content))
	}

	return res, nil
}
