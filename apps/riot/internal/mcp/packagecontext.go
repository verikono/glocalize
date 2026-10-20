package mcp

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
)

// packageContextUnavailable is returned when no guidance file governs the path.
const packageContextUnavailable = "no information available"

// packageMarkers identify a directory as a package root.
var packageMarkers = []string{"package.json", "pyproject.toml", "go.mod", "cargo.toml"}

// agentFileNames are the guidance files searched for, in preference order.
var agentFileNames = []string{".agents.md", "agent.md"}

// onLoadPackageContext returns the package guidance that governs the file or
// directory the caller is working in.
func (m *RiotMcp) onLoadPackageContext(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	raw := stringArg(args, "path")
	if raw == "" {
		return errorResult(fmt.Errorf("load-package-context: path is required")), nil
	}

	path, err := filepath.Abs(raw)
	if err != nil {
		return errorResult(fmt.Errorf("load-package-context: resolve %s: %w", raw, err)), nil
	}

	info, err := os.Stat(path)
	if err != nil {
		return errorResult(fmt.Errorf("load-package-context: %w", err)), nil
	}

	dir := path
	if !info.IsDir() {
		dir = filepath.Dir(path)
	}

	file := findPackageContext(dir)
	if file == "" {
		return textResult(packageContextUnavailable), nil
	}

	content, err := os.ReadFile(file)
	if err != nil {
		return errorResult(fmt.Errorf("load-package-context: read %s: %w", file, err)), nil
	}
	return textResult(string(content)), nil
}

// findPackageContext returns the guidance file that governs dir, or "" when
// none is found. dir and each of its parents are checked in turn; a directory
// is accepted only when it holds BOTH a guidance file and a package marker, so
// a stray guidance file in a non-package directory is skipped rather than
// returned.
func findPackageContext(dir string) string {
	for current := dir; ; current = filepath.Dir(current) {
		if file := agentFileIn(current); file != "" && hasPackageMarker(current) {
			return file
		}
		if current == filepath.Dir(current) {
			return ""
		}
	}
}

// agentFileIn returns the guidance file in dir, or "" when there is none.
func agentFileIn(dir string) string {
	for _, name := range agentFileNames {
		path := filepath.Join(dir, name)
		if info, err := os.Stat(path); err == nil && !info.IsDir() {
			return path
		}
	}
	return ""
}

// hasPackageMarker reports whether dir holds one of the package markers.
func hasPackageMarker(dir string) bool {
	for _, name := range packageMarkers {
		if info, err := os.Stat(filepath.Join(dir, name)); err == nil && !info.IsDir() {
			return true
		}
	}
	return false
}

// textResult wraps plain text in a successful tool result.
func textResult(text string) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: text}},
	}
}
