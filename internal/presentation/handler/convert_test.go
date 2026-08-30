package handler

import (
	"context"
	"testing"

	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewConvertHandler(t *testing.T) {
	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewConvertHandler(uc, nil)
	if h == nil {
		t.Fatal("Expected NewConvertHandler to return a non-nil object")
	}
}

func TestConvertHandler_GetTool(t *testing.T) {
	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewConvertHandler(uc, nil)
	tool := h.GetTool()

	if tool.Name != "markitdown__convert__mlc" {
		t.Errorf("Expected tool name to be markitdown__convert__mlc, got %s", tool.Name)
	}

	expectedProps := []string{
		"uri",
		"force_artifact",
		"output_filename",
		"enable_vision",
		"llm_provider",
		"openai_key",
		"openai_model",
		"openai_url",
		"ollama_model",
		"ollama_url",
	}

	for _, prop := range expectedProps {
		if _, exists := tool.InputSchema.Properties[prop]; !exists {
			t.Errorf("Expected InputSchema to contain property '%s'", prop)
		}
	}
}

func TestConvertHandler_Handle_MissingUri(t *testing.T) {
	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewConvertHandler(uc, nil)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "markitdown__convert__mlc",
			Arguments: map[string]interface{}{},
		},
	}

	res, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected Go error: %v", err)
	}
	if !res.IsError {
		t.Error("Expected error result when uri is missing")
	}
}

func TestConvertHandler_Handle_OpenAiMissingKey(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("DEFAULT_LLM_AUTH_KEY", "")

	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{Provider: "openai"})
	h := NewConvertHandler(uc, nil)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name: "markitdown__convert__mlc",
			Arguments: map[string]interface{}{
				"uri":           "https://example.com/test.jpg",
				"enable_vision": true,
				"llm_provider":  "openai",
			},
		},
	}

	res, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected Go error: %v", err)
	}
	if !res.IsError {
		t.Error("Expected error result when OpenAI API key is missing")
	}
}
