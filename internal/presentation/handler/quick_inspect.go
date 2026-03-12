package handler

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
)

type QuickInspectHandler struct{}

func NewQuickInspectHandler() *QuickInspectHandler {
	return &QuickInspectHandler{}
}

func (h *QuickInspectHandler) GetTool() mcp.Tool {
	return mcp.NewTool(
		"markitdown__quick_inspect__mlc",
		mcp.WithDescription("Quickly retrieves metadata about a document (file size, extension) without full conversion."),
		mcp.WithString("uri", mcp.Description("Path to the file to inspect"), mcp.Required()),
	)
}

func (h *QuickInspectHandler) Handle(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	uri := mcp.ParseString(request, "uri", "")
	if uri == "" {
		return mcp.NewToolResultError("uri is required"), nil
	}

	info, err := os.Stat(uri)
	if err != nil {
		return mcp.NewToolResultErrorFromErr("Failed to inspect file", err), nil
	}

	report := fmt.Sprintf("## File Inspection: %s\n", filepath.Base(uri))
	report += fmt.Sprintf("- **Size**: %d bytes\n", info.Size())
	report += fmt.Sprintf("- **Extension**: %s\n", filepath.Ext(uri))
	report += fmt.Sprintf("- **Last Modified**: %s\n", info.ModTime().Format("2006-01-02 15:04:05"))

	return mcp.NewToolResultText(report), nil
}
