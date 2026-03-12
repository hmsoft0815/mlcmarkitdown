package handler

import (
	"testing"

	"github.com/hmsoft0815/mlc-markitdown/internal/usecase"
)

func TestNewConvertHandler(t *testing.T) {
	uc := usecase.NewConvertUseCase(nil, 1000)
	h := NewConvertHandler(uc)
	if h == nil {
		t.Fatal("Expected NewConvertHandler to return a non-nil object")
	}
}

func TestConvertHandler_GetTool(t *testing.T) {
	uc := usecase.NewConvertUseCase(nil, 1000)
	h := NewConvertHandler(uc)
	tool := h.GetTool()

	if tool.Name != "markitdown__convert__mlc" {
		t.Errorf("Expected tool name to be markitdown__convert__mlc, got %s", tool.Name)
	}
}
