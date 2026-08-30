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

type LlmConfig struct {
	Provider string
	URL      string
	Model    string
	AuthKey  string
}

func NewConvertUseCase(artifactCli *client.Client, threshold int, llmCfg LlmConfig) *ConvertUseCase {
	pythonCmd := os.Getenv("PYTHON_CMD")
	if pythonCmd == "" {
		pythonCmd = "python3"
	}

	provider := llmCfg.Provider
	if provider == "" {
		provider = os.Getenv("DEFAULT_LLM_PROVIDER")
	}
	if provider == "" {
		provider = "openai"
	}

	llmUrl := llmCfg.URL
	if llmUrl == "" {
		llmUrl = os.Getenv("DEFAULT_LLM_URL")
	}

	llmModel := llmCfg.Model
	if llmModel == "" {
		llmModel = os.Getenv("DEFAULT_LLM_MODEL")
	}
	if llmModel == "" {
		if provider == "openai" {
			llmModel = "gpt-4o"
		} else {
			llmModel = "llama3.2-vision"
		}
	}

	llmAuthKey := llmCfg.AuthKey
	if llmAuthKey == "" {
		llmAuthKey = os.Getenv("DEFAULT_LLM_AUTH_KEY")
	}
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

// ConvertRequest holds parameters for document conversion.
type ConvertRequest struct {
	URI            string
	ForceArtifact  bool
	OutputFilename string
	Progress       ProgressFunc
	LlmModel       string
	OpenaiKey      string
	LlmBaseUrl     string
}

func (uc *ConvertUseCase) Convert(ctx context.Context, req ConvertRequest) (string, *pb.WriteResponse, error) {
	if req.Progress != nil {
		req.Progress(10, "Initializing conversion engine...")
	}

	// 1. Call Python shim
	if req.Progress != nil {
		req.Progress(30, "Parsing document...")
	}

	// Assuming the shim is in the fixed location for now
	args := []string{"internal/infrastructure/python/shim.py", req.URI}
	if req.LlmModel != "" {
		args = append(args, "--llm-model", req.LlmModel)
		if req.OpenaiKey != "" {
			args = append(args, "--openai-key", req.OpenaiKey)
		}
		if req.LlmBaseUrl != "" {
			args = append(args, "--llm-base-url", req.LlmBaseUrl)
		}
	}

	cmd := exec.CommandContext(ctx, uc.pythonCmd, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", nil, fmt.Errorf("markitdown failed: %w (output: %s)", err, string(output))
	}

	content := string(output)
	if req.Progress != nil {
		req.Progress(70, "Markdown generated.")
	}

	// 2. Check threshold or force
	if len(content) > uc.threshold || req.ForceArtifact {
		if req.Progress != nil {
			req.Progress(90, "Saving document to artifact storage...")
		}

		filename := "converted_document.md"
		if req.OutputFilename != "" {
			filename = req.OutputFilename
		} else {
			// Try to derive filename from URI
			parts := strings.Split(req.URI, "/")
			if len(parts) > 0 && parts[len(parts)-1] != "" {
				filename = parts[len(parts)-1] + ".md"
			}
		}

		res, err := uc.artifactCli.Write(ctx, filename, []byte(content),
			client.WithSource("mlc-markitdown"),
			client.WithDescription("Auto-archived MarkItDown conversion result"))
		if err != nil {
			return content, nil, fmt.Errorf("failed to auto-archive: %w", err)
		}

		if req.Progress != nil {
			req.Progress(100, "Done.")
		}
		return content, res, nil
	}

	if req.Progress != nil {
		req.Progress(100, "Done.")
	}
	return content, nil, nil
}

func (uc *ConvertUseCase) WriteTempFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
