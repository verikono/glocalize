package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// packageContextUnavailable is returned when no guidance file governs the path.
const packageContextUnavailable = "no information available"

// packageMarkers identify a directory as a package root.
var packageMarkers = []string{"package.json", "pyproject.toml", "go.mod", "cargo.toml"}

// agentFileNames are the guidance files searched for, in preference order.
var agentFileNames = []string{".agents.md", "agent.md"}

// packageGuidance is one resolved guidance file and its contents.
type packageGuidance struct {
	file    string
	content string
}

// onLoadPackageContext returns the package guidance that governs the file or
// directory the caller is working in. It accepts one path or many; when several
// paths resolve to different guidance files their contents are stitched into a
// single result, each preceded by the path of the file it came from.
func (m *RiotMcp) onLoadPackageContext(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	raws := packageContextPaths(args)
	if len(raws) == 0 {
		return errorResult(fmt.Errorf("load-package-context: path or paths is required")), nil
	}

	found := make([]packageGuidance, 0, len(raws))
	seen := make(map[string]bool, len(raws))

	for _, raw := range raws {
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
		if file == "" || seen[file] {
			continue
		}
		seen[file] = true

		content, err := os.ReadFile(file)
		if err != nil {
			return errorResult(fmt.Errorf("load-package-context: read %s: %w", file, err)), nil
		}
		found = append(found, packageGuidance{file: file, content: string(content)})
	}

	switch len(found) {
	case 0:
		return textResult(packageContextUnavailable), nil
	case 1:
		return textResult(found[0].content), nil
	}
	return textResult(stitchGuidance(found)), nil
}

// packageContextPaths collects the paths a load-package-context call names,
// accepting either the single `path` argument or the `paths` array.
func packageContextPaths(args map[string]any) []string {
	var out []string
	if single := stringArg(args, "path"); single != "" {
		out = append(out, single)
	}
	out = append(out, stringSlice(args["paths"])...)
	return out
}

// stitchGuidance joins several guidance files into one document, each section
// headed by the path of the file it came from.
func stitchGuidance(found []packageGuidance) string {
	var b strings.Builder
	for i, g := range found {
		if i > 0 {
			b.WriteString("\n\n---\n\n")
		}
		b.WriteString("# ")
		b.WriteString(g.file)
		b.WriteString("\n\n")
		b.WriteString(g.content)
	}
	return b.String()
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

// onGetProjectContext returns the guidance file at the root of the package
// named by, or containing, the caller's path.
func (m *RiotMcp) onGetProjectContext(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := stringArg(args, "package-name")
	dir := stringArg(args, "directory")
	file := stringArg(args, "filepath")

	supplied := 0
	for _, value := range []string{name, dir, file} {
		if value != "" {
			supplied++
		}
	}
	if supplied == 0 {
		return errorResult(fmt.Errorf("get-project-context: one of package-name, directory or filepath is required")), nil
	}
	if supplied > 1 {
		return errorResult(fmt.Errorf("get-project-context: supply only one of package-name, directory or filepath")), nil
	}

	var root string
	switch {
	case name != "":
		resolved, err := resolvePackageByName(name)
		if err != nil {
			return errorResult(fmt.Errorf("get-project-context: %w", err)), nil
		}
		root = resolved
	case dir != "":
		resolved, err := resolvePackageRoot(dir)
		if err != nil {
			return errorResult(fmt.Errorf("get-project-context: %w", err)), nil
		}
		root = resolved
	default:
		resolved, err := resolvePackageRoot(file)
		if err != nil {
			return errorResult(fmt.Errorf("get-project-context: %w", err)), nil
		}
		root = resolved
	}

	if root == "" {
		return textResult(packageContextUnavailable), nil
	}

	guidance := agentFileIn(root)
	if guidance == "" {
		return textResult(packageContextUnavailable), nil
	}

	content, err := os.ReadFile(guidance)
	if err != nil {
		return errorResult(fmt.Errorf("get-project-context: read %s: %w", guidance, err)), nil
	}
	return textResult(string(content)), nil
}

// resolvePackageRoot returns the closest package root at or above path, or ""
// when the walk reaches the file system root without finding one. A file path
// starts the walk at its directory.
func resolvePackageRoot(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve %s: %w", path, err)
	}
	info, err := os.Stat(abs)
	if err != nil {
		return "", err
	}
	dir := abs
	if !info.IsDir() {
		dir = filepath.Dir(abs)
	}
	for current := dir; ; current = filepath.Dir(current) {
		if hasPackageMarker(current) {
			return current, nil
		}
		if current == filepath.Dir(current) {
			return "", nil
		}
	}
}

// resolvePackageByName finds the package root whose marker declares name, or
// whose directory is named name. Each candidate repository root is walked
// breadth-first, so the shallowest match wins. It returns "" when no package
// matches.
func resolvePackageByName(name string) (string, error) {
	wanted := strings.ToLower(strings.TrimSpace(name))
	for _, root := range repoRoots() {
		if found := searchPackage(root, wanted); found != "" {
			return found, nil
		}
	}
	return "", nil
}

// searchPackage walks root breadth-first looking for a package whose directory
// name or declared name matches wanted.
func searchPackage(root, wanted string) string {
	queue := []string{root}
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]

		if hasPackageMarker(current) {
			if strings.ToLower(filepath.Base(current)) == wanted || strings.ToLower(declaredPackageName(current)) == wanted {
				return current
			}
		}

		entries, err := os.ReadDir(current)
		if err != nil {
			continue
		}
		for _, entry := range entries {
			if !entry.IsDir() || skipDir(entry.Name()) {
				continue
			}
			queue = append(queue, filepath.Join(current, entry.Name()))
		}
	}
	return ""
}

// skipDir reports whether a directory is excluded from the package search.
func skipDir(name string) bool {
	switch name {
	case "node_modules", ".git", "vendor", "dist", "build", "tmp", ".cache", "target":
		return true
	}
	return strings.HasPrefix(name, ".")
}

// repoRoots returns the repository roots to search for a named package: the
// root above the working directory and the root above the running executable.
// A location with no repository marker contributes nothing, so a daemon
// started outside the monorepo still resolves packages from its own binary.
func repoRoots() []string {
	var roots []string
	seen := map[string]bool{}
	add := func(dir string) {
		if dir == "" || seen[dir] {
			return
		}
		seen[dir] = true
		roots = append(roots, dir)
	}

	if dir, err := os.Getwd(); err == nil {
		add(findRepoRoot(dir))
	}
	if exe, err := os.Executable(); err == nil {
		add(findRepoRoot(filepath.Dir(exe)))
	}
	return roots
}

// findRepoRoot walks up from dir to the closest directory holding a repository
// marker, returning "" when none is found.
func findRepoRoot(dir string) string {
	for current := dir; ; current = filepath.Dir(current) {
		for _, marker := range []string{"go.work", ".git"} {
			if info, err := os.Stat(filepath.Join(current, marker)); err == nil && !info.IsDir() {
				return current
			}
		}
		if current == filepath.Dir(current) {
			return ""
		}
	}
}

// declaredPackageName returns the name a package marker declares, or "" when
// no marker declares one.
func declaredPackageName(dir string) string {
	if name := jsonName(filepath.Join(dir, "package.json")); name != "" {
		return name
	}
	if name := tomlName(filepath.Join(dir, "pyproject.toml")); name != "" {
		return name
	}
	if name := tomlName(filepath.Join(dir, "cargo.toml")); name != "" {
		return name
	}
	return goModuleName(filepath.Join(dir, "go.mod"))
}

// jsonName reads the `name` field of a JSON file.
func jsonName(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var doc struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &doc); err != nil {
		return ""
	}
	return doc.Name
}

// tomlName reads the first `name = "..."` line of a TOML file.
func tomlName(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "name") {
			continue
		}
		rest := strings.TrimSpace(strings.TrimPrefix(line, "name"))
		if !strings.HasPrefix(rest, "=") {
			continue
		}
		return strings.Trim(strings.TrimSpace(strings.TrimPrefix(rest, "=")), `"'`)
	}
	return ""
}

// goModuleName reads the module path of a go.mod file.
func goModuleName(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[0] == "module" {
			return fields[1]
		}
	}
	return ""
}
