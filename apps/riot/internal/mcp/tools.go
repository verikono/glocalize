package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
)

const (
	// AgentCheckinToolName is the tool an agent calls first in a session.
	AgentCheckinToolName = "agent-checkin"
	// GetAgentManifestToolName summarises the agent roster.
	GetAgentManifestToolName = "get-agent-manifest"
	// GetProjectsToolName lists the project allow-list.
	GetProjectsToolName = "get-projects"
	// ListSkillsToolName lists the skill catalogue.
	ListSkillsToolName = "list-skills"
	// GetSkillToolName reads a skill's body.
	GetSkillToolName = "get-skill"
	// AddSkillToolName creates a skill.
	AddSkillToolName = "add-skill"
	// UpdateSkillToolName updates a skill.
	UpdateSkillToolName = "update-skill"
	// DeleteSkillToolName deletes a skill.
	DeleteSkillToolName = "delete-skill"
	// MutateSkillToolName creates, updates or deletes a skill and its
	// dependency/dependent associations in one call.
	MutateSkillToolName = "mutate-skill"
	// CreateAgentToolName creates an agent.
	CreateAgentToolName = "create-agent"
	// UpdateAgentToolName updates an agent.
	UpdateAgentToolName = "update-agent"
	// DeleteAgentToolName deletes an agent.
	DeleteAgentToolName = "delete-agent"
	// AssociateAgentSkillToolName binds a skill to an agent.
	AssociateAgentSkillToolName = "associate-agent-skill"
	// AssociateAppSkillToolName binds a skill to an app.
	AssociateAppSkillToolName = "associate-app-skill"
	// AssociateAgentProjectToolName binds an agent to a project.
	AssociateAgentProjectToolName = "associate-agent-project"
	// AddReferenceToolName creates a reference.
	AddReferenceToolName = "add-reference"
	// DropReferenceToolName deletes a reference.
	DropReferenceToolName = "drop-reference"
	// AssociateSkillReferenceToolName attaches a reference to a skill.
	AssociateSkillReferenceToolName = "associate-skill-reference"
	// DropSkillReferenceToolName detaches a reference from a skill.
	DropSkillReferenceToolName = "drop-skill-reference"
)

const priorityVocabulary = "LOW, NORMAL, HIGH, ESSENTIAL (least to most essential)"

// tool definitions

var agentCheckin = mcp.NewTool(
	AgentCheckinToolName,
	mcp.WithDescription("Check in as an agent and receive a stable agent id, returned as an AGENTS.md document listing the skills bound to you by priority. Call this ONCE as the first action of every session, passing your own agent name and the project you are working in; the server returns the same id on every later check-in for the same project+agent. Read the document: it tells you which skills are ESSENTIAL (load them immediately), HIGH, NORMAL and LOW. Consume every skill it lists with the get-skill tool. Thereafter send the returned id in the `agentid` header of every MCP request."),
	mcp.WithString("agent", mcp.Required(), mcp.Description("The calling agent's name, e.g. \"McpikonArchitect\".")),
	mcp.WithString("project", mcp.Required(), mcp.Description("The project the agent is checking in for; must be one of the projects listed by get-projects.")),
)

var getAgentManifest = mcp.NewTool(
	GetAgentManifestToolName,
	mcp.WithDescription("retrieve a summary of every available agent"),
)

var getProjects = mcp.NewTool(
	GetProjectsToolName,
	mcp.WithDescription("List the projects riot knows about, as a JSON array of names. Address a project by name in every other tool — get-skill, add-skill and the rest take a `project` argument drawn from this list."),
)

var listSkills = mcp.NewTool(
	ListSkillsToolName,
	mcp.WithDescription("List every skill you can consume as name, description and url, newest first — call this to see what is on offer before you read one. Bodies are not included: call get-skill with a project and a name to read a skill and follow it."),
)

var getSkill = mcp.NewTool(
	GetSkillToolName,
	mcp.WithDescription("Read a skill and follow it. The result carries TWO content blocks: the FIRST is the skill rendered as markdown — frontmatter generated from the record (name, description, url, project) followed by the skill's body — READ THAT BLOCK AND OBEY IT; the SECOND repeats the stored record as JSON, verbatim, for machine consumers."),
	mcp.WithString("project", mcp.Required(), mcp.Description("The project that owns the skill, e.g. \"all\".")),
	mcp.WithString("name", mcp.Required(), mcp.Description("The skill's name, e.g. \"participating-agent\".")),
)

var addSkill = mcp.NewTool(
	AddSkillToolName,
	mcp.WithDescription("Create a skill in a project. The arguments are the skill's properties: `name`, `description`, `project`, `content` and an optional `url` — a sandbox url of the form skill://<project>/<name> is synthesized when `url` is omitted. The name must be unique within the project, and the owning project must already exist. Returns the stored record as JSON."),
	mcp.WithString("project", mcp.Required(), mcp.Description("The project that will own the skill; the project must already exist.")),
	mcp.WithString("name", mcp.Required(), mcp.Description("The skill's name, unique within the project, e.g. \"font-acquisition\".")),
	mcp.WithString("description", mcp.Required(), mcp.Description("A one-line summary of what the skill does.")),
	mcp.WithString("content", mcp.Required(), mcp.Description("The skill BODY in Markdown: frontmatter removed, newlines left unescaped.")),
	mcp.WithString("keywords", mcp.Description("Optional searchable keyword string for the skill, e.g. \"zustand state react store\".")),
	mcp.WithString("url", mcp.Description("Optional source url. When omitted a sandbox url skill://<project>/<name> is synthesized.")),
)

var updateSkill = mcp.NewTool(
	UpdateSkillToolName,
	mcp.WithDescription("Update an existing skill owned by a project. Only the fields you supply are changed; omitted fields keep their current value. `content` REPLACES the stored body outright. Returns the updated record as JSON."),
	mcp.WithString("project", mcp.Required(), mcp.Description("The project that owns the skill.")),
	mcp.WithString("name", mcp.Required(), mcp.Description("The name of the skill to update, as it exists in the project.")),
	mcp.WithString("description", mcp.Description("A new one-line summary. Omit to leave unchanged.")),
	mcp.WithString("content", mcp.Description("A new body. Omit to leave unchanged.")),
	mcp.WithString("keywords", mcp.Description("A new searchable keyword string. Omit to leave unchanged.")),
	mcp.WithString("url", mcp.Description("A new source url. Omit to leave unchanged.")),
)

var deleteSkill = mcp.NewTool(
	DeleteSkillToolName,
	mcp.WithDescription("Delete a skill from a project. A name matching no skill is reported as an error result."),
	mcp.WithString("project", mcp.Required(), mcp.Description("The project that owns the skill.")),
	mcp.WithString("name", mcp.Required(), mcp.Description("The name of the skill to delete, as it exists in the project.")),
)

// associationMutationSchema is the JSON schema for one side of a skill's
// associations (dependencies or dependents).
func associationMutationSchema() map[string]any {
	return map[string]any{
		"type":        "object",
		"description": "Add, remove or replace one side of a skill's associations. Each key is a skill id or name.",
		"properties": map[string]any{
			"add": map[string]any{
				"type":        "array",
				"description": "Association keys to add. Each key is a skill id or name.",
				"items":       map[string]any{"type": "string"},
			},
			"remove": map[string]any{
				"type":        "array",
				"description": "Association keys to remove. Each key is a skill id or name.",
				"items":       map[string]any{"type": "string"},
			},
			"replace": map[string]any{
				"type":        "boolean",
				"description": "When true the add/remove arrays are treated as the COMPLETE association set: every association not listed in `add` is voided first. Defaults to false.",
			},
		},
	}
}

var mutateSkill = mcp.NewTool(
	MutateSkillToolName,
	mcp.WithDescription("Create, update or delete a skill and its relationships (dependencies, dependents) in one call. The target skill is identified by `id` or `name` — one of them is mandatory so the subject can be pulled back (for `op: \"create\"` the name is the new skill's name). `mutation.op` selects the operation: `create` makes a new skill from `mutation.data`, `update` patches the fields present in `mutation.data`, `delete` removes the skill and every association pointing at it. `mutation.dependencies` are the skills this skill relies on; `mutation.dependents` are the skills that rely on this one. Each side accepts `add` and `remove` arrays of skill ids or names, plus `replace: true` to treat `add` as the complete set (everything else is voided). The whole mutation runs in one transaction."),
	mcp.WithString("id", mcp.Description("The target skill's id. Provide this or `name`.")),
	mcp.WithString("name", mcp.Description("The target skill's name. Provide this or `id`.")),
	mcp.WithObject("mutation", mcp.Required(),
		mcp.Description("The mutation to apply to the target skill."),
		mcp.Properties(map[string]any{
			"id":   map[string]any{"type": "string", "description": "The target skill's id (alternative to the top-level id)."},
			"name": map[string]any{"type": "string", "description": "The target skill's name (alternative to the top-level name)."},
			"op": map[string]any{
				"type":        "string",
				"enum":        []string{"create", "update", "delete"},
				"description": "The operation: create, update or delete.",
			},
			"data": map[string]any{
				"type":                 "object",
				"description":          "Key/value pairs of the skill fields being written: name, description, url, content, keywords.",
				"additionalProperties": true,
			},
			"dependencies": associationMutationSchema(),
			"dependents":   associationMutationSchema(),
		}),
	),
)

var createAgent = mcp.NewTool(
	CreateAgentToolName,
	mcp.WithDescription("Create a new agent. The name is the identity and must be unique; creating a second agent with the same name is rejected. allowed-tools is a JSON-encoded array and metadata a JSON-encoded object; both default to [] and {} when omitted."),
	mcp.WithString("name", mcp.Required(), mcp.Description("The agent's name, unique across agents, e.g. \"McpikonArchitect\".")),
	mcp.WithString("description", mcp.Required(), mcp.Description("A one-line summary of what the agent does.")),
	mcp.WithString("allowed-tools", mcp.Description("Optional JSON-encoded array of allowed tool names. Defaults to [].")),
	mcp.WithString("metadata", mcp.Description("Optional JSON-encoded object of arbitrary metadata. Defaults to {}.")),
)

var updateAgent = mcp.NewTool(
	UpdateAgentToolName,
	mcp.WithDescription("Update an existing agent. Only the fields you supply are changed; omitted fields keep their current value. The name is the identity and cannot be changed."),
	mcp.WithString("name", mcp.Required(), mcp.Description("The name of the agent to update, as it exists.")),
	mcp.WithString("description", mcp.Description("A new one-line summary. Omit to leave unchanged.")),
	mcp.WithString("allowed-tools", mcp.Description("A new JSON-encoded array of allowed tool names. Omit to leave unchanged.")),
	mcp.WithString("metadata", mcp.Description("A new JSON-encoded object of metadata. Omit to leave unchanged.")),
)

var deleteAgent = mcp.NewTool(
	DeleteAgentToolName,
	mcp.WithDescription("Delete another agent. Deleting a name that does not exist is reported as an error result."),
	mcp.WithString("name", mcp.Required(), mcp.Description("The name of the agent to delete, as it exists.")),
)

var associateAgentSkill = mcp.NewTool(
	AssociateAgentSkillToolName,
	mcp.WithDescription("Associate a skill with an agent, recording how essential that skill is to the agent. The pair (agent, skill) is the identity, so re-associating a skill updates its priority instead of erroring. The agent and the skill must both already exist."),
	mcp.WithString("agent", mcp.Required(), mcp.Description("The agent to bind the skill to, by name, e.g. \"golang-developer\".")),
	mcp.WithString("skill", mcp.Required(), mcp.Description("The skill to bind, by name, e.g. \"font-acquisition\".")),
	mcp.WithString("priority", mcp.Description("How essential the skill is to the agent: "+priorityVocabulary+". Defaults to NORMAL when omitted.")),
)

var associateAppSkill = mcp.NewTool(
	AssociateAppSkillToolName,
	mcp.WithDescription("Associate a skill with an app, recording how essential that skill is to the app. The pair (app, skill) is the identity, so re-associating a skill updates its priority instead of erroring. The app and the skill must both already exist."),
	mcp.WithString("app", mcp.Required(), mcp.Description("The app to bind the skill to, by name, e.g. \"mcpikon\".")),
	mcp.WithString("skill", mcp.Required(), mcp.Description("The skill to bind, by name, e.g. \"font-acquisition\".")),
	mcp.WithString("priority", mcp.Description("How essential the skill is to the app: "+priorityVocabulary+". Defaults to NORMAL when omitted.")),
)

var associateAgentProject = mcp.NewTool(
	AssociateAgentProjectToolName,
	mcp.WithDescription("Associate an agent with a project, recording that the agent works in that project. The pair (agent, project) is the identity, so re-associating is a no-op instead of an error. The agent and the project must both already exist."),
	mcp.WithString("agent", mcp.Required(), mcp.Description("The agent to bind to the project, by name, e.g. \"golang-developer\".")),
	mcp.WithString("project", mcp.Required(), mcp.Description("The project to bind the agent to, by name, e.g. \"kaltra\".")),
)

var addReference = mcp.NewTool(
	AddReferenceToolName,
	mcp.WithDescription("Create a standalone reference: a resource (documentation, an image, a video, commands, a script, a configuration) that skills can draw on. The reference is created on its own and is NOT attached to any skill yet — use associate-skill-reference to attach it. The name is the identity and must be free."),
	mcp.WithString("name", mcp.Required(), mcp.Description("The reference's unique name, e.g. \"go-testing-style-guide\".")),
	mcp.WithString("description", mcp.Required(), mcp.Description("One line on what the reference is and when a skill should reach for it.")),
	mcp.WithString("accessType", mcp.Required(), mcp.Description("How the reference is addressed: static or url. A value outside those two is rejected.")),
	mcp.WithString("resourceType", mcp.Required(), mcp.Description("What the reference is: literature, image, video, commands, script or configuration. A value outside those six is rejected.")),
	mcp.WithString("content", mcp.Description("The payload when accessType is static, or the address to fetch when it is url. Optional.")),
)

var dropReference = mcp.NewTool(
	DropReferenceToolName,
	mcp.WithDescription("Delete a reference and every skill link that points at it, in one transaction. Use drop-skill-reference instead to detach it from one skill while keeping the reference available to the others. A name matching no reference is reported as an error."),
	mcp.WithString("name", mcp.Required(), mcp.Description("The reference to delete, by name.")),
)

var associateSkillReference = mcp.NewTool(
	AssociateSkillReferenceToolName,
	mcp.WithDescription("Attach an existing reference to a skill, so the reference is one of the resources that skill draws on. The reference must already exist — create it with add-reference first. The pair (skill, reference) is the identity, so re-associating is a no-op instead of an error."),
	mcp.WithString("skill", mcp.Required(), mcp.Description("The skill to attach the reference to, by name.")),
	mcp.WithString("reference", mcp.Required(), mcp.Description("The reference to attach, by name.")),
)

var dropSkillReference = mcp.NewTool(
	DropSkillReferenceToolName,
	mcp.WithDescription("Detach a reference from a skill WITHOUT deleting the reference: only the link is removed. Detaching a pair that was never linked is a no-op. Use drop-reference to remove the reference itself."),
	mcp.WithString("skill", mcp.Required(), mcp.Description("The skill to detach the reference from, by name.")),
	mcp.WithString("reference", mcp.Required(), mcp.Description("The reference to detach, by name.")),
)

// helpers

func stringArg(args map[string]any, key string) string {
	v, _ := args[key].(string)
	return strings.TrimSpace(v)
}

func jsonResult(v any) (*mcp.CallToolResult, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: string(b)}},
	}, nil
}

func errorResult(err error) *mcp.CallToolResult {
	return &mcp.CallToolResult{
		Content: []mcp.Content{mcp.TextContent{Type: "text", Text: err.Error()}},
		IsError: true,
	}
}

func eq[M any, T any](column string, value T) orm.Predicate[M] {
	return orm.Field[M, T]{Column: column}.EQ(value)
}

func unique[M any](column string, value any) orm.UniquePredicate[M] {
	return orm.UniquePredicate[M]{
		Data: orm.PredicateData{Column: column, Operator: "=", Value: value},
	}
}

func priorityArg(args map[string]any, key string) (orm.PriorityType, error) {
	raw := stringArg(args, key)
	if raw == "" {
		return orm.Priority_NORMAL, nil
	}
	switch strings.ToUpper(raw) {
	case "LOW":
		return orm.Priority_LOW, nil
	case "NORMAL":
		return orm.Priority_NORMAL, nil
	case "HIGH":
		return orm.Priority_HIGH, nil
	case "ESSENTIAL":
		return orm.Priority_ESSENTIAL, nil
	default:
		return "", fmt.Errorf("priority must be one of %s", priorityVocabulary)
	}
}

func accessTypeArg(raw string) (orm.URIType, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "static":
		return orm.URI_STATIC, nil
	case "url":
		return orm.URI_URL, nil
	default:
		return "", fmt.Errorf("accessType must be one of: static, url")
	}
}

func resourceTypeArg(raw string) (orm.ReferenceTypeType, error) {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "literature":
		return orm.ReferenceType_LITERATURE, nil
	case "image":
		return orm.ReferenceType_IMAGE, nil
	case "video":
		return orm.ReferenceType_VIDEO, nil
	case "commands":
		return orm.ReferenceType_COMMANDS, nil
	case "script":
		return orm.ReferenceType_SCRIPT, nil
	case "configuration":
		return orm.ReferenceType_CONFIGURATION, nil
	default:
		return "", fmt.Errorf("resourceType must be one of: literature, image, video, commands, script, configuration")
	}
}

// skillByID finds a skill by id.
func (m *RiotMcp) skillByID(ctx context.Context, id string) (*orm.Skill, error) {
	return m.db.Skill.FindFirst(eq[orm.Skill, string]("id", id)).Exec(ctx)
}

// skillByName finds a skill by name (best-effort sandbox identity).
func (m *RiotMcp) skillByName(ctx context.Context, name string) (*orm.Skill, error) {
	return m.db.Skill.FindFirst(eq[orm.Skill, string]("name", name)).Exec(ctx)
}

// agentByName finds an agent by name.
func (m *RiotMcp) agentByName(ctx context.Context, name string) (*orm.Agent, error) {
	return m.db.Agent.FindFirst(eq[orm.Agent, string]("name", name)).Exec(ctx)
}

// projectByName finds a project by name.
func (m *RiotMcp) projectByName(ctx context.Context, name string) (*orm.Project, error) {
	return m.db.Project.FindFirst(eq[orm.Project, string]("name", name)).Exec(ctx)
}
