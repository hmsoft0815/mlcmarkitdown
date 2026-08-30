package usecase

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewConvertUseCase_Defaults(t *testing.T) {
	// Clear any ambient env vars during test
	t.Setenv("DEFAULT_LLM_PROVIDER", "")
	t.Setenv("DEFAULT_LLM_URL", "")
	t.Setenv("DEFAULT_LLM_MODEL", "")
	t.Setenv("DEFAULT_LLM_AUTH_KEY", "")
	t.Setenv("OPENAI_API_KEY", "")

	uc := NewConvertUseCase(nil, 1000, LlmConfig{})
	if uc == nil {
		t.Fatal("Expected NewConvertUseCase to return a non-nil object")
		return
	}
	if uc.threshold != 1000 {
		t.Errorf("Expected threshold to be 1000, got %d", uc.threshold)
	}

	provider, model, url, authKey := uc.GetLlmDefaults()
	if provider != "openai" {
		t.Errorf("Expected default provider 'openai', got '%s'", provider)
	}
	if model != "gpt-4o" {
		t.Errorf("Expected default model 'gpt-4o', got '%s'", model)
	}
	if url != "" {
		t.Errorf("Expected empty default url, got '%s'", url)
	}
	if authKey != "" {
		t.Errorf("Expected empty default authKey, got '%s'", authKey)
	}
}

func TestNewConvertUseCase_OllamaConfig(t *testing.T) {
	cfg := LlmConfig{
		Provider: "ollama",
		URL:      "http://macstudio.fritz.box:11434/v1",
		Model:    "llama3.2-vision",
		AuthKey:  "ollama-dummy",
	}

	uc := NewConvertUseCase(nil, 5000, cfg)
	if uc == nil {
		t.Fatal("Expected non-nil ConvertUseCase")
		return
	}

	provider, model, url, authKey := uc.GetLlmDefaults()
	if provider != "ollama" {
		t.Errorf("Expected provider 'ollama', got '%s'", provider)
	}
	if model != "llama3.2-vision" {
		t.Errorf("Expected model 'llama3.2-vision', got '%s'", model)
	}
	if url != "http://macstudio.fritz.box:11434/v1" {
		t.Errorf("Expected url 'http://macstudio.fritz.box:11434/v1', got '%s'", url)
	}
	if authKey != "ollama-dummy" {
		t.Errorf("Expected authKey 'ollama-dummy', got '%s'", authKey)
	}
}

func TestNewConvertUseCase_EnvFallbacks(t *testing.T) {
	t.Setenv("DEFAULT_LLM_PROVIDER", "ollama")
	t.Setenv("DEFAULT_LLM_URL", "http://localhost:11434/v1")
	t.Setenv("DEFAULT_LLM_MODEL", "")
	t.Setenv("DEFAULT_LLM_AUTH_KEY", "env-key")

	uc := NewConvertUseCase(nil, 2000, LlmConfig{})
	provider, model, url, authKey := uc.GetLlmDefaults()

	if provider != "ollama" {
		t.Errorf("Expected provider from env 'ollama', got '%s'", provider)
	}
	if model != "llama3.2-vision" {
		t.Errorf("Expected model fallback 'llama3.2-vision', got '%s'", model)
	}
	if url != "http://localhost:11434/v1" {
		t.Errorf("Expected url from env 'http://localhost:11434/v1', got '%s'", url)
	}
	if authKey != "env-key" {
		t.Errorf("Expected authKey from env 'env-key', got '%s'", authKey)
	}
}

func TestConvertUseCase_WriteTempFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test_artifact.txt")
	testData := []byte("temporary file content")

	uc := NewConvertUseCase(nil, 1000, LlmConfig{})
	err := uc.WriteTempFile(filePath, testData)
	if err != nil {
		t.Fatalf("WriteTempFile failed: %v", err)
	}

	readData, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("Failed to read created temp file: %v", err)
	}
	if string(readData) != string(testData) {
		t.Errorf("Expected '%s', got '%s'", string(testData), string(readData))
	}
}
