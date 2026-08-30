package handler

import (
	"context"
	"testing"

	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
	"github.com/mark3labs/mcp-go/mcp"
)

func TestNewConvertArtifactHandler(t *testing.T) {
	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewConvertArtifactHandler(uc, nil, nil)
	if h == nil {
		t.Fatal("Expected NewConvertArtifactHandler to return a non-nil object")
	}
}

func TestConvertArtifactHandler_GetTool(t *testing.T) {
	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewConvertArtifactHandler(uc, nil, nil)
	tool := h.GetTool()

	if tool.Name != "markitdown__convert_artifact__mlc" {
		t.Errorf("Expected tool name 'markitdown__convert_artifact__mlc', got '%s'", tool.Name)
	}

	expectedProps := []string{
		"artifactId",
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

func TestConvertArtifactHandler_Handle_MissingArtifactId(t *testing.T) {
	uc := usecase.NewConvertUseCase(nil, 1000, usecase.LlmConfig{})
	h := NewConvertArtifactHandler(uc, nil, nil)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "markitdown__convert_artifact__mlc",
			Arguments: map[string]interface{}{},
		},
	}

	res, err := h.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("Unexpected Go error: %v", err)
	}
	if !res.IsError {
		t.Error("Expected error result when artifactId is missing")
	}
}
