package mcp

import (
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/server"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
	T "github.com/verikono/glocalize/types/go"
)

const (
	// ServerName is the MCP server's advertised name.
	ServerName = "riot-mcp"

	// ServerVersion is the MCP server's advertised version.
	ServerVersion = "0.1.0"
)

// RiotMcp is the standalone MCP server: the full mcpikon tool surface served
// over streamable HTTP, backed by the riot orm client.
type RiotMcp struct {
	db         *orm.DB
	provider   string
	server     *server.MCPServer
	httpServer *server.StreamableHTTPServer
}

// NewRiotMcp builds the server, registers its tools and returns it serving-ready.
func NewRiotMcp(options ...T.OptionerFn) (*RiotMcp, error) {
	inst := &RiotMcp{}
	for _, option := range options {
		if err := option(inst); err != nil {
			return nil, err
		}
	}
	if inst.db == nil {
		return nil, fmt.Errorf("mcp: an orm client is required (WithDB)")
	}

	inst.server = server.NewMCPServer(ServerName, ServerVersion)
	inst.httpServer = server.NewStreamableHTTPServer(inst.server)
	inst.registerTools()
	return inst, nil
}

// WithDB injects the orm client the tool handlers read and write.
func WithDB(db *orm.DB) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if inst, ok := arg.(*RiotMcp); ok {
				inst.db = db
			}
		}
		return nil
	}
}

// WithProvider records the sql dialect for join-table writes.
func WithProvider(provider string) T.OptionerFn {
	return func(args ...any) error {
		for _, arg := range args {
			if inst, ok := arg.(*RiotMcp); ok {
				inst.provider = provider
			}
		}
		return nil
	}
}

// Handler returns the streamable HTTP handler for the MCP endpoint.
func (m *RiotMcp) Handler() http.Handler {
	return m.httpServer
}

// registerTools registers the full tool surface with the underlying server.
func (m *RiotMcp) registerTools() {
	m.server.AddTool(agentCheckin, m.onAgentCheckin)
	m.server.AddTool(getAgentManifest, m.onGetAgentManifest)
	m.server.AddTool(getProjects, m.onGetProjects)
	m.server.AddTool(listSkills, m.onListSkills)
	m.server.AddTool(getSkill, m.onGetSkill)
	m.server.AddTool(addSkill, m.onAddSkill)
	m.server.AddTool(updateSkill, m.onUpdateSkill)
	m.server.AddTool(deleteSkill, m.onDeleteSkill)
	m.server.AddTool(createAgent, m.onCreateAgent)
	m.server.AddTool(updateAgent, m.onUpdateAgent)
	m.server.AddTool(deleteAgent, m.onDeleteAgent)
	m.server.AddTool(associateAgentSkill, m.onAssociateAgentSkill)
	m.server.AddTool(associateAppSkill, m.onAssociateAppSkill)
	m.server.AddTool(associateAgentProject, m.onAssociateAgentProject)
	m.server.AddTool(addReference, m.onAddReference)
	m.server.AddTool(dropReference, m.onDropReference)
	m.server.AddTool(associateSkillReference, m.onAssociateSkillReference)
	m.server.AddTool(dropSkillReference, m.onDropSkillReference)
}
