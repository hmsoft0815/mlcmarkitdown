package handler

import (
	"context"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

func TestPromptHandler_GetPrompt(t *testing.T) {
	h := NewPromptHandler()
	prompt := h.GetPrompt()

	if prompt.Name != "markitdown_instruction" {
		t.Errorf("Expected prompt name 'markitdown_instruction', got '%s'", prompt.Name)
	}
	if prompt.Description == "" {
		t.Error("Expected non-empty prompt description")
	}
}

func TestPromptHandler_Handle(t *testing.T) {
	h := NewPromptHandler()
	res, err := h.Handle(context.Background(), mcp.GetPromptRequest{})
	if err != nil {
		t.Fatalf("Unexpected error from PromptHandler.Handle: %v", err)
	}

	if len(res.Messages) == 0 {
		t.Fatal("Expected at least one prompt message")
	}

	msgContent, ok := res.Messages[0].Content.(mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent message, got %T", res.Messages[0].Content)
	}

	if !strings.Contains(msgContent.Text, "ollama") {
		t.Errorf("Expected prompt instruction to mention ollama, got: %s", msgContent.Text)
	}
}
