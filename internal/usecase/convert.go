package usecase

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/hmsoft0815/mlcartifact/client"
	pb "github.com/hmsoft0815/mlcartifact/proto"
)

type ProgressFunc func(int, string)

type ConvertUseCase struct {
	artifactCli       *client.Client
	threshold         int
	pythonCmd         string
	defaultProvider   string
	defaultLlmUrl     string
	defaultLlmModel   string
	defaultLlmAuthKey string
}

func NewConvertUseCase(artifactCli *client.Client, threshold int) *ConvertUseCase {
	pythonCmd := os.Getenv("PYTHON_CMD")
	if pythonCmd == "" {
		pythonCmd = "python3"
	}
	provider := os.Getenv("DEFAULT_LLM_PROVIDER")
	if provider == "" {
		provider = "openai"
	}
	llmUrl := os.Getenv("DEFAULT_LLM_URL")
	llmModel := os.Getenv("DEFAULT_LLM_MODEL")
	if llmModel == "" {
		if provider == "openai" {
			llmModel = "gpt-4o"
		} else {
			llmModel = "llama3.2-vision"
		}
	}
	llmAuthKey := os.Getenv("DEFAULT_LLM_AUTH_KEY")
	if llmAuthKey == "" && provider == "openai" {
		llmAuthKey = os.Getenv("OPENAI_API_KEY")
	}

	return &ConvertUseCase{
		artifactCli:       artifactCli,
		threshold:         threshold,
		pythonCmd:         pythonCmd,
		defaultProvider:   provider,
		defaultLlmUrl:     llmUrl,
		defaultLlmModel:   llmModel,
		defaultLlmAuthKey: llmAuthKey,
	}
}

func (uc *ConvertUseCase) GetLlmDefaults() (string, string, string, string) {
	return uc.defaultProvider, uc.defaultLlmModel, uc.defaultLlmUrl, uc.defaultLlmAuthKey
}

func (uc *ConvertUseCase) Convert(ctx context.Context, uri string, force bool, progress ProgressFunc, llmModel string, openaiKey string, llmBaseUrl string) (string, *pb.WriteResponse, error) {
	if progress != nil {
		progress(10, "Initializing conversion engine...")
	}

	// 1. Call Python shim
	if progress != nil {
		progress(30, "Parsing document...")
	}

	// Assuming the shim is in the fixed location for now
	args := []string{"internal/infrastructure/python/shim.py", uri}
	if llmModel != "" {
		args = append(args, "--llm-model", llmModel)
		if openaiKey != "" {
			args = append(args, "--openai-key", openaiKey)
		}
		if llmBaseUrl != "" {
			args = append(args, "--llm-base-url", llmBaseUrl)
		}
	}

	cmd := exec.CommandContext(ctx, uc.pythonCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil, fmt.Errorf("markitdown failed: %w (output: %s)", err, string(output))
	}

	content := string(output)
	if progress != nil {
		progress(70, "Markdown generated.")
	}

	// 2. Check threshold or force
	if len(content) > uc.threshold || force {
		if progress != nil {
			progress(90, "Saving document to artifact storage...")
		}

		filename := "converted_document.md"
		// Try to derive filename from URI
		parts := strings.Split(uri, "/")
		if len(parts) > 0 {
			filename = parts[len(parts)-1] + ".md"
		}

		res, err := uc.artifactCli.Write(ctx, filename, []byte(content),
			client.WithSource("mlc-markitdown"),
			client.WithDescription("Auto-archived MarkItDown conversion result"))
		if err != nil {
			return content, nil, fmt.Errorf("failed to auto-archive: %w", err)
		}

		if progress != nil {
			progress(100, "Done.")
		}
		return content, res, nil
	}

	if progress != nil {
		progress(100, "Done.")
	}
	return content, nil, nil
}

func (uc *ConvertUseCase) WriteTempFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
