# Testing mlc-markitdown

This project uses a two-tier testing strategy: standard Go unit tests and automated integration tests using `mcp-tester`.

## 1. Unit Tests
Standard Go unit tests are located in `internal/usecase/` and `internal/presentation/handler/`.

To run all unit tests:
```bash
make test
```

## 2. Integration Tests (mcp-tester)
We use `mlc_mcptester` (https://github.com/hmsoft0815/mlc_mcptester) to run automated scripts against the compiled MCP server. This ensures that the tool registration, argument parsing, and core conversion logic (including the Python shim) work correctly.

The integration test script is located at `tests/integration.mcp`.

### Running Integration Tests
The `Makefile` automatically checks for `mcp-tester` and installs it if missing:

```bash
make test-integration
```

### Script Syntax
The script uses the `.mcp` DSL:
- `call_tool <name> <key>:<value>`: Invokes an MCP tool.
- `assert_contains "<substring>"`: Verifies the response content.

## 3. Progress Reporting
The server supports real-time progress notifications during document conversion. You can observe these when running `mcp-tester` with the verbose flag:

```bash
make build
mcp-tester test -c "bin/mlc-markitdown" --script tests/integration.mcp -v
```

Look for `[PROGRESS]` entries in the output.

## 4. Manual Testing
You can manually interact with the server using `mcp-tester` CLI:

```bash
# List available tools
mcp-tester list -c "bin/mlc-markitdown"

# Call a tool manually
mcp-tester call markitdown__convert__mlc -a '{"uri": "test.txt"}' -c "bin/mlc-markitdown"
```
