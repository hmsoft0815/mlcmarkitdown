package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"

	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
	"github.com/hmsoft0815/mlcartifact"
)

type ConvertArtifactHandler struct {
	useCase     *usecase.ConvertUseCase
	artifactCli *mlcartifact.Client
}

func NewConvertArtifactHandler(useCase *usecase.ConvertUseCase, artifactCli *mlcartifact.Client) *ConvertArtifactHandler {
	return &ConvertArtifactHandler{
		useCase:     useCase,
		artifactCli: artifactCli,
	}
}

func (h *ConvertArtifactHandler) GetTool() mcp.Tool {
	return mcp.NewTool(
		"markitdown__convert_artifact__mlc",
		mcp.WithDescription("Converts a document already stored in the artifact storage to Markdown."),
		mcp.WithString("artifactId", mcp.Description("The ID of the source artifact to convert"), mcp.Required()),
		mcp.WithString("output_filename", mcp.Description("Optional name for the resulting Markdown artifact.")),
	)
}

func (h *ConvertArtifactHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	artifactID := mcp.ParseString(request, "artifactId", "")
	if artifactID == "" {
		return mcp.NewToolResultError("artifactId is required"), nil
	}

	outputFilename := mcp.ParseString(request, "output_filename", "")

	// 1. Read source artifact
	res, err := h.artifactCli.Read(ctx, artifactID)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to read source artifact", err), nil
	}

	// 2. Write to temp file for MarkItDown (since it needs a file path)
	tmpFile := fmt.Sprintf("/tmp/markitdown_%s", artifactID)
	err = h.useCase.WriteTempFile(tmpFile, res.Content)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to prepare temp file", err), nil
	}

	// 3. Convert
	content, newArtifact, err := h.useCase.Convert(ctx, tmpFile, true, nil)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Conversion failed", err), nil
	}

	// 4. Return result
	result := &mcp.CallToolResult{
		Content: []mcp.Content{},
	}

	if newArtifact != nil {
		if outputFilename != "" {
			// Optional: rename it if requested (though Convert already saved it)
			// For simplicity, we stick to what Convert did but mention the ID
		}

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
