package mcp

import (
	"fmt"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
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

// toolDefinitions returns every tool definition the server exposes, in
// registration order. It is static: no database is needed to describe the
// surface, which is what `riot tools` relies on.
func toolDefinitions() []mcp.Tool {
	return []mcp.Tool{
		agentCheckin,
		getAgentManifest,
		listAgents,
		getProjects,
		listSkills,
		getSkill,
		addSkill,
		updateSkill,
		deleteSkill,
		mutateSkill,
		createAgent,
		updateAgent,
		deleteAgent,
		associateAgentSkill,
		associateAppSkill,
		associateAgentProject,
		addReference,
		getReference,
		dropReference,
		associateSkillReference,
		dropSkillReference,
		loadPackageContext,
		getProjectContext,
	}
}

// Tools returns every tool definition the server exposes, in registration
// order, for callers that only need to describe the surface.
func Tools() []mcp.Tool {
	return toolDefinitions()
}

// handlerFor returns the handler bound to a tool name, or nil when the name is
// not part of the surface.
func (m *RiotMcp) handlerFor(name string) server.ToolHandlerFunc {
	switch name {
	case AgentCheckinToolName:
		return m.onAgentCheckin
	case GetAgentManifestToolName:
		return m.onGetAgentManifest
	case ListAgentsToolName:
		return m.onListAgents
	case GetProjectsToolName:
		return m.onGetProjects
	case ListSkillsToolName:
		return m.onListSkills
	case GetSkillToolName:
		return m.onGetSkill
	case AddSkillToolName:
		return m.onAddSkill
	case UpdateSkillToolName:
		return m.onUpdateSkill
	case DeleteSkillToolName:
		return m.onDeleteSkill
	case MutateSkillToolName:
		return m.onMutateSkill
	case CreateAgentToolName:
		return m.onCreateAgent
	case UpdateAgentToolName:
		return m.onUpdateAgent
	case DeleteAgentToolName:
		return m.onDeleteAgent
	case AssociateAgentSkillToolName:
		return m.onAssociateAgentSkill
	case AssociateAppSkillToolName:
		return m.onAssociateAppSkill
	case AssociateAgentProjectToolName:
		return m.onAssociateAgentProject
	case AddReferenceToolName:
		return m.onAddReference
	case GetReferenceToolName:
		return m.onGetReference
	case DropReferenceToolName:
		return m.onDropReference
	case AssociateSkillReferenceToolName:
		return m.onAssociateSkillReference
	case DropSkillReferenceToolName:
		return m.onDropSkillReference
	case LoadPackageContextToolName:
		return m.onLoadPackageContext
	case GetProjectContextToolName:
		return m.onGetProjectContext
	}
	return nil
}

// registerTools registers the full tool surface with the underlying server.
func (m *RiotMcp) registerTools() {
	for _, tool := range toolDefinitions() {
		m.server.AddTool(tool, m.handlerFor(tool.Name))
	}
}
