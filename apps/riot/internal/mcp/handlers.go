package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/mark3labs/mcp-go/mcp"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
)

// ph returns the SQL placeholder for index i, dialect-aware.
func (m *RiotMcp) ph(i int) string {
	if m.provider == "postgres" {
		return fmt.Sprintf("$%d", i)
	}
	return "?"
}

func (m *RiotMcp) onAgentCheckin(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	agentName := stringArg(args, "agent")
	project := stringArg(args, "project")
	if agentName == "" {
		return errorResult(fmt.Errorf("agent-checkin: agent is required")), nil
	}
	if project == "" {
		return errorResult(fmt.Errorf("agent-checkin: project is required")), nil
	}

	proj, err := m.projectByName(ctx, project)
	if err != nil {
		return errorResult(fmt.Errorf("agent-checkin: %w", err)), nil
	}
	if proj == nil {
		return errorResult(fmt.Errorf("agent-checkin: project %q is not accepted", project)), nil
	}

	agent, err := m.agentByName(ctx, agentName)
	if err != nil {
		return errorResult(fmt.Errorf("agent-checkin: %w", err)), nil
	}
	if agent == nil {
		agent, err = m.db.Agent.Create().SetName(agentName).SetDescription(agentName).SetAllowedTools("[]").SetMetadata("{}").Exec(ctx)
		if err != nil {
			return errorResult(fmt.Errorf("agent-checkin: %w", err)), nil
		}
	}

	skills, err := m.boundSkills(ctx, agent.Id)
	if err != nil {
		return errorResult(fmt.Errorf("agent-checkin: %w", err)), nil
	}

	document := renderAgentDocument(agent.Name, agent.Id, project, skills)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: document},
			mcp.TextContent{Type: "text", Text: fmt.Sprintf(`{"agentId":%q,"project":%q}`, agent.Id, project)},
		},
	}, nil
}

type boundSkill struct {
	Name        string
	Description string
	Priority    orm.PriorityType
}

func (m *RiotMcp) boundSkills(ctx context.Context, agentID string) ([]boundSkill, error) {
	rows, err := m.db.AgentSkill.FindMany(eq[orm.AgentSkill, string]("agentId", agentID)).Exec(ctx)
	if err != nil {
		return nil, err
	}
	skills := make([]boundSkill, 0, len(rows))
	for _, row := range rows {
		skill, err := m.skillByID(ctx, row.SkillId)
		if err != nil {
			return nil, err
		}
		if skill == nil {
			continue
		}
		skills = append(skills, boundSkill{Name: skill.Name, Description: skill.Description, Priority: row.Priority})
	}
	return skills, nil
}

func renderAgentDocument(name, id, project string, skills []boundSkill) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# %s\n\n", name)
	fmt.Fprintf(&b, "You are this agent. You have just checked in and been given a stable identity.\n\n")
	fmt.Fprintf(&b, "## Rules\n")
	fmt.Fprintf(&b, "- You are **%s**, agent id `%s`, working in project **%s**.\n", name, id, project)
	fmt.Fprintf(&b, "- add your agentid as a request header on every MCP call to this mcp service.\n")
	fmt.Fprintf(&b, "- consume a skill with the get-skill tool — it takes priority over every other mechanism.\n")
	fmt.Fprintf(&b, "- load every ESSENTIAL skill bound to you before any other work.\n\n")

	buckets := map[orm.PriorityType][]boundSkill{}
	var order []orm.PriorityType
	for _, s := range skills {
		if _, ok := buckets[s.Priority]; !ok {
			order = append(order, s.Priority)
		}
		buckets[s.Priority] = append(buckets[s.Priority], s)
	}

	if len(skills) == 0 {
		b.WriteString("## Skills\n\n_no skills bound yet_\n")
		return b.String()
	}

	b.WriteString("## Skills\n\n")
	labels := map[orm.PriorityType]string{
		orm.Priority_ESSENTIAL: "ESSENTIAL — mandatory — load these immediately",
		orm.Priority_HIGH:      "HIGH — you will almost certainly require these",
		orm.Priority_NORMAL:    "NORMAL — you will likely need these",
		orm.Priority_LOW:       "LOW — be aware they exist",
	}
	for _, p := range order {
		fmt.Fprintf(&b, "### %s\n", labels[p])
		for _, s := range buckets[p] {
			fmt.Fprintf(&b, "- **%s** — %s\n", s.Name, s.Description)
		}
		b.WriteString("\n")
	}
	return b.String()
}

func (m *RiotMcp) onGetAgentManifest(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	records, err := m.db.Agent.FindMany().OrderBy(orm.Field[orm.Agent, string]{Column: "name"}.Asc()).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("get-agent-manifest: %w", err)), nil
	}
	type summary struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	out := make([]summary, 0, len(records))
	for _, r := range records {
		out = append(out, summary{Name: r.Name, Description: r.Description})
	}
	return jsonResult(out)
}

func (m *RiotMcp) onGetProjects(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	records, err := m.db.Project.FindMany().OrderBy(orm.Field[orm.Project, string]{Column: "name"}.Asc()).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("get-projects: %w", err)), nil
	}
	names := make([]string, 0, len(records))
	for _, r := range records {
		names = append(names, r.Name)
	}
	return jsonResult(names)
}

func (m *RiotMcp) onListSkills(ctx context.Context, _ mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	records, err := m.db.Skill.FindMany().OrderBy(orm.Field[orm.Skill, time.Time]{Column: "createdAt"}.Desc()).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("list-skills: %w", err)), nil
	}
	type summary struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		URL         string `json:"url"`
		Keywords    string `json:"keywords"`
	}
	out := make([]summary, 0, len(records))
	for _, r := range records {
		out = append(out, summary{Name: r.Name, Description: r.Description, URL: r.Url, Keywords: r.Keywords})
	}
	return jsonResult(out)
}

func (m *RiotMcp) onGetSkill(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	project := stringArg(args, "project")
	name := stringArg(args, "name")
	if project == "" {
		return errorResult(fmt.Errorf("get-skill: project is required")), nil
	}
	if name == "" {
		return errorResult(fmt.Errorf("get-skill: name is required")), nil
	}
	skill, err := m.skillByName(ctx, name)
	if err != nil {
		return errorResult(fmt.Errorf("get-skill: %w", err)), nil
	}
	if skill == nil {
		return errorResult(fmt.Errorf("get-skill: skill %q not found in project %q", name, project)), nil
	}
	document := fmt.Sprintf("---\nname: %s\ndescription: %s\nurl: %s\nkeywords: %s\nproject: %s\n---\n%s\n", skill.Name, skill.Description, skill.Url, skill.Keywords, project, skill.Content)
	record, err := jsonMarshal(skill)
	if err != nil {
		return errorResult(err), nil
	}
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{Type: "text", Text: document},
			mcp.TextContent{Type: "text", Text: record},
		},
	}, nil
}

func (m *RiotMcp) onAddSkill(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	project := stringArg(args, "project")
	name := stringArg(args, "name")
	description := stringArg(args, "description")
	content := stringArg(args, "content")
	keywords := stringArg(args, "keywords")
	url := stringArg(args, "url")
	if project == "" || name == "" || description == "" {
		return errorResult(fmt.Errorf("add-skill: project, name and description are required")), nil
	}
	if proj, err := m.projectByName(ctx, project); err != nil || proj == nil {
		if err != nil {
			return errorResult(fmt.Errorf("add-skill: %w", err)), nil
		}
		return errorResult(fmt.Errorf("add-skill: project %q does not exist", project)), nil
	}
	if existing, err := m.skillByName(ctx, name); err != nil {
		return errorResult(fmt.Errorf("add-skill: %w", err)), nil
	} else if existing != nil {
		return errorResult(fmt.Errorf("add-skill: skill %q already exists in project %q", name, project)), nil
	}
	if url == "" {
		url = fmt.Sprintf("skill://%s/%s", project, name)
	}
	skill, err := m.db.Skill.Create().SetName(name).SetDescription(description).SetContent(content).SetKeywords(keywords).SetUrl(url).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("add-skill: %w", err)), nil
	}
	if _, err := m.db.ProjectSkill.Create().SetProjectName(project).SetSkillId(skill.Id).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("add-skill: link to project: %w", err)), nil
	}
	return jsonResult(skill)
}

func (m *RiotMcp) onUpdateSkill(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := stringArg(args, "name")
	if name == "" {
		return errorResult(fmt.Errorf("update-skill: name is required")), nil
	}
	skill, err := m.skillByName(ctx, name)
	if err != nil {
		return errorResult(fmt.Errorf("update-skill: %w", err)), nil
	}
	if skill == nil {
		return errorResult(fmt.Errorf("update-skill: skill %q not found", name)), nil
	}
	builder := m.db.Skill.Update(unique[orm.Skill]("id", skill.Id))
	if v := stringArg(args, "description"); v != "" {
		builder.SetDescription(v)
	}
	if v := stringArg(args, "url"); v != "" {
		builder.SetUrl(v)
	}
	if v, ok := args["content"].(string); ok && v != "" {
		builder.SetContent(v)
	}
	if v, ok := args["keywords"].(string); ok {
		builder.SetKeywords(v)
	}
	updated, err := builder.Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("update-skill: %w", err)), nil
	}
	return jsonResult(updated)
}

func (m *RiotMcp) onDeleteSkill(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := stringArg(args, "name")
	if name == "" {
		return errorResult(fmt.Errorf("delete-skill: name is required")), nil
	}
	skill, err := m.skillByName(ctx, name)
	if err != nil {
		return errorResult(fmt.Errorf("delete-skill: %w", err)), nil
	}
	if skill == nil {
		return errorResult(fmt.Errorf("delete-skill: skill %q not found", name)), nil
	}
	if _, err := m.db.ProjectSkill.DeleteMany(eq[orm.ProjectSkill, string]("skillId", skill.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-skill: %w", err)), nil
	}
	if _, err := m.db.AgentSkill.DeleteMany(eq[orm.AgentSkill, string]("skillId", skill.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-skill: %w", err)), nil
	}
	if _, err := m.db.AppSkill.DeleteMany(eq[orm.AppSkill, string]("skillId", skill.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-skill: %w", err)), nil
	}
	if _, err := m.db.SkillReference.DeleteMany(eq[orm.SkillReference, string]("skillId", skill.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-skill: %w", err)), nil
	}
	if _, err := m.db.SkillDependency.DeleteMany(eq[orm.SkillDependency, string]("dependencyId", skill.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-skill: %w", err)), nil
	}
	if _, err := m.db.SkillDependency.DeleteMany(eq[orm.SkillDependency, string]("dependentId", skill.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-skill: %w", err)), nil
	}
	if _, err := m.db.Skill.Delete(unique[orm.Skill]("id", skill.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-skill: %w", err)), nil
	}
	return jsonResult(map[string]any{"deleted": true, "name": skill.Name})
}

// associationMutation is one side of a skill's associations, as supplied to
// mutate-skill. Add and Remove hold skill ids or names; Replace voids every
// association on that side before Add is applied.
type associationMutation struct {
	Add     []string
	Remove  []string
	Replace bool
}

// parseAssociationMutation reads an association object from the mutation.
func parseAssociationMutation(raw any) (associationMutation, error) {
	out := associationMutation{}
	obj, ok := raw.(map[string]any)
	if !ok || obj == nil {
		return out, nil
	}
	out.Add = stringSlice(obj["add"])
	out.Remove = stringSlice(obj["remove"])
	if v, ok := obj["replace"].(bool); ok {
		out.Replace = v
	}
	return out, nil
}

// stringSlice coerces a JSON array of strings into a []string, dropping blanks.
func stringSlice(raw any) []string {
	items, ok := raw.([]any)
	if !ok {
		if direct, ok := raw.([]string); ok {
			items = make([]any, len(direct))
			for i, v := range direct {
				items[i] = v
			}
		} else {
			return nil
		}
	}
	out := make([]string, 0, len(items))
	for _, item := range items {
		if s, ok := item.(string); ok && strings.TrimSpace(s) != "" {
			out = append(out, strings.TrimSpace(s))
		}
	}
	return out
}

// resolveSkillKey finds a skill by id first, then by name.
func (m *RiotMcp) resolveSkillKey(ctx context.Context, key string) (*orm.Skill, error) {
	if skill, err := m.skillByID(ctx, key); err != nil {
		return nil, err
	} else if skill != nil {
		return skill, nil
	}
	return m.skillByName(ctx, key)
}

// applyAssociations reconciles one side of a skill's associations inside a
// transaction. selfColumn is the column holding the subject skill's id and
// otherColumn the column holding the associated skill's id.
func (m *RiotMcp) applyAssociations(ctx context.Context, tx *orm.Tx, subjectID, selfColumn, otherColumn string, mutation associationMutation) error {
	if mutation.Replace {
		if _, err := tx.SkillDependency.DeleteMany(eq[orm.SkillDependency, string](selfColumn, subjectID)).Exec(ctx); err != nil {
			return err
		}
	}
	for _, key := range mutation.Remove {
		other, err := m.resolveSkillKey(ctx, key)
		if err != nil {
			return err
		}
		if other == nil {
			continue
		}
		if _, err := tx.SkillDependency.DeleteMany(
			eq[orm.SkillDependency, string](selfColumn, subjectID),
			eq[orm.SkillDependency, string](otherColumn, other.Id),
		).Exec(ctx); err != nil {
			return err
		}
	}
	for _, key := range mutation.Add {
		other, err := m.resolveSkillKey(ctx, key)
		if err != nil {
			return err
		}
		if other == nil {
			return fmt.Errorf("skill %q not found", key)
		}
		if other.Id == subjectID {
			continue
		}
		existing, err := tx.SkillDependency.FindFirst(
			eq[orm.SkillDependency, string](selfColumn, subjectID),
			eq[orm.SkillDependency, string](otherColumn, other.Id),
		).Exec(ctx)
		if err != nil {
			return err
		}
		if existing != nil {
			continue
		}
		builder := tx.SkillDependency.Create()
		if selfColumn == "dependentId" {
			builder.SetDependentId(subjectID).SetDependencyId(other.Id)
		} else {
			builder.SetDependentId(other.Id).SetDependencyId(subjectID)
		}
		if _, err := builder.Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

// onMutateSkill creates, updates or deletes a skill and reconciles its
// dependency and dependent associations in a single transaction.
func (m *RiotMcp) onMutateSkill(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	mutation, ok := args["mutation"].(map[string]any)
	if !ok || mutation == nil {
		return errorResult(fmt.Errorf("mutate-skill: mutation is required")), nil
	}

	op := strings.ToLower(strings.TrimSpace(stringValue(mutation["op"])))
	if op == "" {
		return errorResult(fmt.Errorf("mutate-skill: mutation.op is required (create, update or delete)")), nil
	}

	// the target may be named at the top level or inside the mutation.
	id := firstNonEmpty(stringArg(args, "id"), stringValue(mutation["id"]))
	name := firstNonEmpty(stringArg(args, "name"), stringValue(mutation["name"]))
	if id == "" && name == "" {
		return errorResult(fmt.Errorf("mutate-skill: one of id or name is required so the target skill can be pulled back")), nil
	}

	data, _ := mutation["data"].(map[string]any)
	dependencies, err := parseAssociationMutation(mutation["dependencies"])
	if err != nil {
		return errorResult(fmt.Errorf("mutate-skill: dependencies: %w", err)), nil
	}
	dependents, err := parseAssociationMutation(mutation["dependents"])
	if err != nil {
		return errorResult(fmt.Errorf("mutate-skill: dependents: %w", err)), nil
	}

	var result *orm.Skill
	err = m.db.Transaction(ctx, func(tx *orm.Tx) error {
		var subject *orm.Skill
		if id != "" {
			subject, err = tx.Skill.FindFirst(eq[orm.Skill, string]("id", id)).Exec(ctx)
		} else {
			subject, err = tx.Skill.FindFirst(eq[orm.Skill, string]("name", name)).Exec(ctx)
		}
		if err != nil {
			return err
		}

		switch op {
		case "create":
			if subject != nil {
				return fmt.Errorf("skill %q already exists", subject.Name)
			}
			newName := firstNonEmpty(stringValue(data["name"]), name)
			if newName == "" {
				return fmt.Errorf("create requires a name")
			}
			// name, description and content are required columns: reject an
			// incomplete create here rather than surfacing a raw validation error.
			if stringValue(data["description"]) == "" {
				return fmt.Errorf("create requires data.description")
			}
			if stringValue(data["content"]) == "" {
				return fmt.Errorf("create requires data.content")
			}
			builder := tx.Skill.Create().
				SetName(newName).
				SetDescription(stringValue(data["description"])).
				SetUrl(firstNonEmpty(stringValue(data["url"]), "skill://"+newName)).
				SetContent(stringValue(data["content"])).
				SetKeywords(stringValue(data["keywords"]))
			created, err := builder.Exec(ctx)
			if err != nil {
				return err
			}
			subject = created
		case "update":
			if subject == nil {
				return fmt.Errorf("skill %q not found", firstNonEmpty(id, name))
			}
			builder := tx.Skill.Update(unique[orm.Skill]("id", subject.Id))
			if v, ok := data["name"].(string); ok && v != "" {
				builder.SetName(v)
			}
			if v, ok := data["description"].(string); ok && v != "" {
				builder.SetDescription(v)
			}
			if v, ok := data["url"].(string); ok && v != "" {
				builder.SetUrl(v)
			}
			if v, ok := data["content"].(string); ok && v != "" {
				builder.SetContent(v)
			}
			if v, ok := data["keywords"].(string); ok {
				builder.SetKeywords(v)
			}
			updated, err := builder.Exec(ctx)
			if err != nil {
				return err
			}
			subject = updated
		case "delete":
			if subject == nil {
				return fmt.Errorf("skill %q not found", firstNonEmpty(id, name))
			}
			if err := m.deleteSkillAssociations(ctx, tx, subject.Id); err != nil {
				return err
			}
			if _, err := tx.Skill.Delete(unique[orm.Skill]("id", subject.Id)).Exec(ctx); err != nil {
				return err
			}
			result = subject
			return nil
		default:
			return fmt.Errorf("mutation.op must be one of: create, update, delete")
		}

		// dependencies: rows where this skill is the dependent.
		if err := m.applyAssociations(ctx, tx, subject.Id, "dependentId", "dependencyId", dependencies); err != nil {
			return err
		}
		// dependents: rows where this skill is the dependency.
		if err := m.applyAssociations(ctx, tx, subject.Id, "dependencyId", "dependentId", dependents); err != nil {
			return err
		}

		result = subject
		return nil
	})
	if err != nil {
		return errorResult(fmt.Errorf("mutate-skill: %w", err)), nil
	}

	return jsonResult(map[string]any{"op": op, "skill": result})
}

// deleteSkillAssociations removes every row that points at a skill.
func (m *RiotMcp) deleteSkillAssociations(ctx context.Context, tx *orm.Tx, skillID string) error {
	if _, err := tx.ProjectSkill.DeleteMany(eq[orm.ProjectSkill, string]("skillId", skillID)).Exec(ctx); err != nil {
		return err
	}
	if _, err := tx.AgentSkill.DeleteMany(eq[orm.AgentSkill, string]("skillId", skillID)).Exec(ctx); err != nil {
		return err
	}
	if _, err := tx.AppSkill.DeleteMany(eq[orm.AppSkill, string]("skillId", skillID)).Exec(ctx); err != nil {
		return err
	}
	if _, err := tx.SkillReference.DeleteMany(eq[orm.SkillReference, string]("skillId", skillID)).Exec(ctx); err != nil {
		return err
	}
	if _, err := tx.SkillDependency.DeleteMany(eq[orm.SkillDependency, string]("dependencyId", skillID)).Exec(ctx); err != nil {
		return err
	}
	if _, err := tx.SkillDependency.DeleteMany(eq[orm.SkillDependency, string]("dependentId", skillID)).Exec(ctx); err != nil {
		return err
	}
	return nil
}

// stringValue coerces a JSON value to a trimmed string.
func stringValue(v any) string {
	if s, ok := v.(string); ok {
		return strings.TrimSpace(s)
	}
	return ""
}

// firstNonEmpty returns the first non-empty string.
func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

func (m *RiotMcp) onCreateAgent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := stringArg(args, "name")
	description := stringArg(args, "description")
	if name == "" || description == "" {
		return errorResult(fmt.Errorf("create-agent: name and description are required")), nil
	}
	if existing, err := m.agentByName(ctx, name); err != nil {
		return errorResult(fmt.Errorf("create-agent: %w", err)), nil
	} else if existing != nil {
		return errorResult(fmt.Errorf("create-agent: agent %q already exists", name)), nil
	}
	allowedTools := stringArg(args, "allowed-tools")
	if allowedTools == "" {
		allowedTools = "[]"
	}
	metadata := stringArg(args, "metadata")
	if metadata == "" {
		metadata = "{}"
	}
	agent, err := m.db.Agent.Create().SetName(name).SetDescription(description).SetAllowedTools(allowedTools).SetMetadata(metadata).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("create-agent: %w", err)), nil
	}
	if err := writeAgentFile(agent.Name, agent.Description, agent.AllowedTools); err != nil {
		return errorResult(fmt.Errorf("create-agent: %w", err)), nil
	}
	return jsonResult(agent)
}

func (m *RiotMcp) onUpdateAgent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := stringArg(args, "name")
	if name == "" {
		return errorResult(fmt.Errorf("update-agent: name is required")), nil
	}
	agent, err := m.agentByName(ctx, name)
	if err != nil {
		return errorResult(fmt.Errorf("update-agent: %w", err)), nil
	}
	if agent == nil {
		return errorResult(fmt.Errorf("update-agent: agent %q not found", name)), nil
	}
	builder := m.db.Agent.Update(unique[orm.Agent]("id", agent.Id))
	if v := stringArg(args, "description"); v != "" {
		builder.SetDescription(v)
	}
	if v := stringArg(args, "allowed-tools"); v != "" {
		builder.SetAllowedTools(v)
	}
	if v := stringArg(args, "metadata"); v != "" {
		builder.SetMetadata(v)
	}
	updated, err := builder.Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("update-agent: %w", err)), nil
	}
	if err := writeAgentFile(updated.Name, updated.Description, updated.AllowedTools); err != nil {
		return errorResult(fmt.Errorf("update-agent: %w", err)), nil
	}
	return jsonResult(updated)
}

func (m *RiotMcp) onDeleteAgent(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := stringArg(args, "name")
	if name == "" {
		return errorResult(fmt.Errorf("delete-agent: name is required")), nil
	}
	agent, err := m.agentByName(ctx, name)
	if err != nil {
		return errorResult(fmt.Errorf("delete-agent: %w", err)), nil
	}
	if agent == nil {
		return errorResult(fmt.Errorf("delete-agent: agent %q not found", name)), nil
	}
	if _, err := m.db.AgentSkill.DeleteMany(eq[orm.AgentSkill, string]("agentId", agent.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-agent: %w", err)), nil
	}
	if _, err := m.db.TrashAgent.DeleteMany(eq[orm.TrashAgent, string]("agentId", agent.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-agent: %w", err)), nil
	}
	if _, err := m.db.Agent.Delete(unique[orm.Agent]("id", agent.Id)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("delete-agent: %w", err)), nil
	}
	if err := removeAgentFile(agent.Name); err != nil {
		return errorResult(fmt.Errorf("delete-agent: %w", err)), nil
	}
	return jsonResult(map[string]any{"deleted": true, "name": agent.Name})
}

func (m *RiotMcp) onAssociateAgentSkill(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	agentName := stringArg(args, "agent")
	skillName := stringArg(args, "skill")
	if agentName == "" || skillName == "" {
		return errorResult(fmt.Errorf("associate-agent-skill: agent and skill are required")), nil
	}
	priority, err := priorityArg(args, "priority")
	if err != nil {
		return errorResult(fmt.Errorf("associate-agent-skill: %w", err)), nil
	}
	agent, err := m.agentByName(ctx, agentName)
	if err != nil {
		return errorResult(fmt.Errorf("associate-agent-skill: %w", err)), nil
	}
	if agent == nil {
		return errorResult(fmt.Errorf("associate-agent-skill: agent %q not found", agentName)), nil
	}
	skill, err := m.skillByName(ctx, skillName)
	if err != nil {
		return errorResult(fmt.Errorf("associate-agent-skill: %w", err)), nil
	}
	if skill == nil {
		return errorResult(fmt.Errorf("associate-agent-skill: skill %q not found", skillName)), nil
	}
	existing, err := m.db.AgentSkill.FindFirst(eq[orm.AgentSkill, string]("agentId", agent.Id), eq[orm.AgentSkill, string]("skillId", skill.Id)).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("associate-agent-skill: %w", err)), nil
	}
	if existing != nil {
		if _, err := m.db.AgentSkill.Update(unique[orm.AgentSkill]("agentId", agent.Id), eq[orm.AgentSkill, string]("skillId", skill.Id)).SetPriority(priority).Exec(ctx); err != nil {
			return errorResult(fmt.Errorf("associate-agent-skill: %w", err)), nil
		}
	} else {
		if _, err := m.db.AgentSkill.Create().SetAgentId(agent.Id).SetSkillId(skill.Id).SetPriority(priority).Exec(ctx); err != nil {
			return errorResult(fmt.Errorf("associate-agent-skill: %w", err)), nil
		}
	}
	return jsonResult(map[string]any{"associated": true, "agent": agentName, "skill": skillName, "priority": string(priority)})
}

func (m *RiotMcp) onAssociateAppSkill(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	appName := stringArg(args, "app")
	skillName := stringArg(args, "skill")
	if appName == "" || skillName == "" {
		return errorResult(fmt.Errorf("associate-app-skill: app and skill are required")), nil
	}
	priority, err := priorityArg(args, "priority")
	if err != nil {
		return errorResult(fmt.Errorf("associate-app-skill: %w", err)), nil
	}
	app, err := m.db.App.FindFirst(eq[orm.App, string]("name", appName)).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("associate-app-skill: %w", err)), nil
	}
	if app == nil {
		return errorResult(fmt.Errorf("associate-app-skill: app %q not found", appName)), nil
	}
	skill, err := m.skillByName(ctx, skillName)
	if err != nil {
		return errorResult(fmt.Errorf("associate-app-skill: %w", err)), nil
	}
	if skill == nil {
		return errorResult(fmt.Errorf("associate-app-skill: skill %q not found", skillName)), nil
	}
	existing, err := m.db.AppSkill.FindFirst(eq[orm.AppSkill, string]("app_id", app.Name), eq[orm.AppSkill, string]("skillId", skill.Id)).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("associate-app-skill: %w", err)), nil
	}
	if existing != nil {
		if _, err := m.db.AppSkill.Update(unique[orm.AppSkill]("app_id", app.Name), eq[orm.AppSkill, string]("skillId", skill.Id)).SetPriority(priority).Exec(ctx); err != nil {
			return errorResult(fmt.Errorf("associate-app-skill: %w", err)), nil
		}
	} else {
		if _, err := m.db.AppSkill.Create().SetApp_id(app.Name).SetSkillId(skill.Id).SetPriority(priority).Exec(ctx); err != nil {
			return errorResult(fmt.Errorf("associate-app-skill: %w", err)), nil
		}
	}
	return jsonResult(map[string]any{"associated": true, "app": appName, "skill": skillName, "priority": string(priority)})
}

func (m *RiotMcp) onAssociateAgentProject(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	agentName := stringArg(args, "agent")
	projectName := stringArg(args, "project")
	if agentName == "" || projectName == "" {
		return errorResult(fmt.Errorf("associate-agent-project: agent and project are required")), nil
	}
	agent, err := m.agentByName(ctx, agentName)
	if err != nil {
		return errorResult(fmt.Errorf("associate-agent-project: %w", err)), nil
	}
	if agent == nil {
		return errorResult(fmt.Errorf("associate-agent-project: agent %q not found", agentName)), nil
	}
	project, err := m.projectByName(ctx, projectName)
	if err != nil {
		return errorResult(fmt.Errorf("associate-agent-project: %w", err)), nil
	}
	if project == nil {
		return errorResult(fmt.Errorf("associate-agent-project: project %q not found", projectName)), nil
	}

	var count int
	sel := fmt.Sprintf(`SELECT COUNT(*) FROM "_ProjectAgents" WHERE "agentId" = %s AND "projectName" = %s`, m.ph(1), m.ph(2))
	if err := m.db.Raw().QueryRowContext(ctx, sel, agent.Id, projectName).Scan(&count); err != nil {
		return errorResult(fmt.Errorf("associate-agent-project: %w", err)), nil
	}
	if count == 0 {
		ins := fmt.Sprintf(`INSERT INTO "_ProjectAgents" ("agentId", "projectName") VALUES (%s, %s)`, m.ph(1), m.ph(2))
		if _, err := m.db.Raw().ExecContext(ctx, ins, agent.Id, projectName); err != nil {
			return errorResult(fmt.Errorf("associate-agent-project: %w", err)), nil
		}
	}
	return jsonResult(map[string]any{"associated": true, "agent": agentName, "project": projectName})
}

func (m *RiotMcp) onAddReference(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := stringArg(args, "name")
	description := stringArg(args, "description")
	if name == "" || description == "" {
		return errorResult(fmt.Errorf("add-reference: name and description are required")), nil
	}
	accessType, err := accessTypeArg(stringArg(args, "accessType"))
	if err != nil {
		return errorResult(fmt.Errorf("add-reference: %w", err)), nil
	}
	resourceType, err := resourceTypeArg(stringArg(args, "resourceType"))
	if err != nil {
		return errorResult(fmt.Errorf("add-reference: %w", err)), nil
	}
	if existing, err := m.db.Reference.FindFirst(eq[orm.Reference, string]("name", name)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("add-reference: %w", err)), nil
	} else if existing != nil {
		return errorResult(fmt.Errorf("add-reference: reference %q already exists", name)), nil
	}
	builder := m.db.Reference.Create().SetName(name).SetDescription(description).SetAccessType(accessType).SetResourceType(resourceType)
	if v := stringArg(args, "content"); v != "" {
		builder.SetContent(v)
	}
	record, err := builder.Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("add-reference: %w", err)), nil
	}
	return jsonResult(record)
}

func (m *RiotMcp) onDropReference(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	name := stringArg(args, "name")
	if name == "" {
		return errorResult(fmt.Errorf("drop-reference: name is required")), nil
	}
	ref, err := m.db.Reference.FindFirst(eq[orm.Reference, string]("name", name)).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("drop-reference: %w", err)), nil
	}
	if ref == nil {
		return errorResult(fmt.Errorf("drop-reference: reference %q not found", name)), nil
	}
	if _, err := m.db.SkillReference.DeleteMany(eq[orm.SkillReference, string]("referenceName", name)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("drop-reference: %w", err)), nil
	}
	if _, err := m.db.Reference.Delete(unique[orm.Reference]("name", name)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("drop-reference: %w", err)), nil
	}
	return jsonResult(map[string]any{"dropped": true, "name": name})
}

func (m *RiotMcp) onAssociateSkillReference(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	skillName := stringArg(args, "skill")
	referenceName := stringArg(args, "reference")
	if skillName == "" || referenceName == "" {
		return errorResult(fmt.Errorf("associate-skill-reference: skill and reference are required")), nil
	}
	skill, err := m.skillByName(ctx, skillName)
	if err != nil {
		return errorResult(fmt.Errorf("associate-skill-reference: %w", err)), nil
	}
	if skill == nil {
		return errorResult(fmt.Errorf("associate-skill-reference: skill %q not found", skillName)), nil
	}
	ref, err := m.db.Reference.FindFirst(eq[orm.Reference, string]("name", referenceName)).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("associate-skill-reference: %w", err)), nil
	}
	if ref == nil {
		return errorResult(fmt.Errorf("associate-skill-reference: reference %q not found", referenceName)), nil
	}
	existing, err := m.db.SkillReference.FindFirst(eq[orm.SkillReference, string]("skillId", skill.Id), eq[orm.SkillReference, string]("referenceName", referenceName)).Exec(ctx)
	if err != nil {
		return errorResult(fmt.Errorf("associate-skill-reference: %w", err)), nil
	}
	if existing == nil {
		if _, err := m.db.SkillReference.Create().SetSkillId(skill.Id).SetReferenceName(referenceName).Exec(ctx); err != nil {
			return errorResult(fmt.Errorf("associate-skill-reference: %w", err)), nil
		}
	}
	return jsonResult(map[string]any{"associated": true, "skill": skillName, "reference": referenceName})
}

func (m *RiotMcp) onDropSkillReference(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	args := request.GetArguments()
	skillName := stringArg(args, "skill")
	referenceName := stringArg(args, "reference")
	if skillName == "" || referenceName == "" {
		return errorResult(fmt.Errorf("drop-skill-reference: skill and reference are required")), nil
	}
	skill, err := m.skillByName(ctx, skillName)
	if err != nil {
		return errorResult(fmt.Errorf("drop-skill-reference: %w", err)), nil
	}
	if skill == nil {
		return errorResult(fmt.Errorf("drop-skill-reference: skill %q not found", skillName)), nil
	}
	if _, err := m.db.SkillReference.DeleteMany(eq[orm.SkillReference, string]("skillId", skill.Id), eq[orm.SkillReference, string]("referenceName", referenceName)).Exec(ctx); err != nil {
		return errorResult(fmt.Errorf("drop-skill-reference: %w", err)), nil
	}
	return jsonResult(map[string]any{"detached": true, "skill": skillName, "reference": referenceName})
}

// jsonMarshal is a tiny alias so handlers read naturally.
func jsonMarshal(v any) (string, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
