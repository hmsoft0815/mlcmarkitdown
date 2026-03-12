package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
)

type ConvertHandler struct {
	useCase *usecase.ConvertUseCase
}

func NewConvertHandler(useCase *usecase.ConvertUseCase) *ConvertHandler {
	return &ConvertHandler{
		useCase: useCase,
	}
}

func (h *ConvertHandler) GetTool() mcp.Tool {
	return mcp.NewTool(
		"markitdown__convert__mlc",
		mcp.WithDescription("Converts a file or URL to Markdown. Smart auto-archiving is applied for large outputs."),
		mcp.WithString("uri", mcp.Description("The source path or URL to convert"), mcp.Required()),
		mcp.WithBoolean("force_artifact", mcp.Description("If true, always save as artifact storage and return a notice.")),
	)
}

func (h *ConvertHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uri := mcp.ParseString(request, "uri", "")
	if uri == "" {
		return mcp.NewToolResultError("uri is required"), nil
	}

	forceArtifact := mcp.ParseBoolean(request, "force_artifact", false)

	// 1. Define progress monitor
	progress := func(percent int, status string) {
		// Emit progress via MCP if possible
		// (Need to pass the progress callback context if mcp-go supports it)
		// For now, we are using the simple approach or server-side logging
	}

	// 2. Call Usecase
	content, artifact, err := h.useCase.Convert(ctx, uri, forceArtifact, progress)
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
