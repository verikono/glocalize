package gql

import (
	"context"
	"time"

	"github.com/graphql-go/graphql"

	orm "github.com/verikono/glocalize/apps/riot/internal/orm"
)

// dateTimeScalar serialises time.Time as RFC 3339.
var dateTimeScalar = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "DateTime",
	Description: "RFC 3339 timestamp",
	Serialize: func(value interface{}) interface{} {
		switch t := value.(type) {
		case time.Time:
			return t.Format(time.RFC3339)
		case *time.Time:
			if t == nil {
				return nil
			}
			return t.Format(time.RFC3339)
		}
		return value
	},
})

// resolver carries the orm client the field resolvers read from.
type resolver struct {
	db *orm.DB
}

func (r *resolver) ctx(p graphql.ResolveParams) context.Context {
	if p.Context == nil {
		return context.Background()
	}
	return p.Context
}

// eq builds a field-equality predicate against a column.
func eq[M any, T any](column string, value T) orm.Predicate[M] {
	return orm.Field[M, T]{Column: column}.EQ(value)
}

// buildSchema assembles the executable read-only schema. Object types are built
// locally so relation resolvers can close over the resolver instance.
func buildSchema(db *orm.DB) (graphql.Schema, error) {
	r := &resolver{db: db}

	skillType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Skill",
		Fields: graphql.Fields{
			"id":          &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"name":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"description": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"url":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"content":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"createdAt":   &graphql.Field{Type: dateTimeScalar},
			"updatedAt":   &graphql.Field{Type: dateTimeScalar},
		},
	})

	styleTokenType := graphql.NewObject(graphql.ObjectConfig{
		Name: "StyleToken",
		Fields: graphql.Fields{
			"style":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"value":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"mode":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"type":     &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"token":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"variable": &graphql.Field{Type: graphql.String},
		},
	})

	agentType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Agent",
		Fields: graphql.Fields{
			"id":           &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"name":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"description":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"content":      &graphql.Field{Type: graphql.String},
			"allowedTools": &graphql.Field{Type: graphql.String},
			"metadata":     &graphql.Field{Type: graphql.String},
			"createdAt":    &graphql.Field{Type: dateTimeScalar},
			"updatedAt":    &graphql.Field{Type: dateTimeScalar},
			"skills":       &graphql.Field{Type: graphql.NewList(skillType), Resolve: r.agentSkills},
		},
	})

	appType := graphql.NewObject(graphql.ObjectConfig{
		Name: "App",
		Fields: graphql.Fields{
			"name":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"projectId": &graphql.Field{Type: graphql.String},
			"skills":    &graphql.Field{Type: graphql.NewList(skillType), Resolve: r.appSkills},
			"styles":    &graphql.Field{Type: graphql.NewList(styleTokenType), Resolve: r.appStyles},
		},
	})

	userType := graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.NewNonNull(graphql.ID)},
			"name":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"preferred": &graphql.Field{Type: graphql.NewList(skillType), Resolve: r.userPreferred},
			"blacklist": &graphql.Field{Type: graphql.NewList(skillType), Resolve: r.userBlacklist},
		},
	})

	projectType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Project",
		Fields: graphql.Fields{
			"name":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"skillsPath": &graphql.Field{Type: graphql.String},
			"agentsPath": &graphql.Field{Type: graphql.String},
			"skills":     &graphql.Field{Type: graphql.NewList(skillType), Resolve: r.projectSkills},
			"agents":     &graphql.Field{Type: graphql.NewList(agentType), Resolve: r.projectAgents},
			"apps":       &graphql.Field{Type: graphql.NewList(appType), Resolve: r.projectApps},
		},
	})

	documentType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Document",
		Fields: graphql.Fields{
			"id":        &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"kind":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"body":      &graphql.Field{Type: graphql.String},
			"updatedAt": &graphql.Field{Type: graphql.NewNonNull(graphql.Int)},
		},
	})

	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"Document": &graphql.Field{
				Type: documentType,
				Args: graphql.FieldConfigArgument{
					"id":   &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"kind": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					preds := []orm.PredicateOf[orm.Document]{eq[orm.Document, string]("id", p.Args["id"].(string))}
					if kind, ok := p.Args["kind"].(string); ok && kind != "" {
						preds = append(preds, eq[orm.Document, string]("kind", kind))
					}
					return r.db.Document.FindFirst(preds...).Exec(r.ctx(p))
				},
			},
			"Documents": &graphql.Field{
				Type:    graphql.NewList(documentType),
				Resolve: r.documents,
			},
			"Agent": &graphql.Field{
				Type: agentType,
				Args: graphql.FieldConfigArgument{
					"id":   &graphql.ArgumentConfig{Type: graphql.String},
					"name": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: r.agent,
			},
			"Agents": &graphql.Field{
				Type:    graphql.NewList(agentType),
				Resolve: r.agents,
			},
			"App": &graphql.Field{
				Type: appType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.app,
			},
			"Apps": &graphql.Field{
				Type:    graphql.NewList(appType),
				Resolve: r.apps,
			},
			"Skill": &graphql.Field{
				Type: skillType,
				Args: graphql.FieldConfigArgument{
					"id":   &graphql.ArgumentConfig{Type: graphql.String},
					"name": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: r.skill,
			},
			"Skills": &graphql.Field{
				Type:    graphql.NewList(skillType),
				Resolve: r.skills,
			},
			"User": &graphql.Field{
				Type: userType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.user,
			},
			"Users": &graphql.Field{
				Type:    graphql.NewList(userType),
				Resolve: r.users,
			},
			"Project": &graphql.Field{
				Type: projectType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.project,
			},
			"Projects": &graphql.Field{
				Type:    graphql.NewList(projectType),
				Resolve: r.projects,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{Query: query})
}

// Entrypoint resolvers.

func (r *resolver) documents(p graphql.ResolveParams) (interface{}, error) {
	return r.db.Document.FindMany().Exec(r.ctx(p))
}

func (r *resolver) agent(p graphql.ResolveParams) (interface{}, error) {
	if id, ok := p.Args["id"].(string); ok && id != "" {
		return r.db.Agent.FindFirst(eq[orm.Agent, string]("id", id)).Exec(r.ctx(p))
	}
	if name, ok := p.Args["name"].(string); ok && name != "" {
		return r.db.Agent.FindFirst(eq[orm.Agent, string]("name", name)).Exec(r.ctx(p))
	}
	return r.db.Agent.FindFirst().Exec(r.ctx(p))
}

func (r *resolver) agents(p graphql.ResolveParams) (interface{}, error) {
	return r.db.Agent.FindMany().Exec(r.ctx(p))
}

func (r *resolver) app(p graphql.ResolveParams) (interface{}, error) {
	return r.db.App.FindFirst(eq[orm.App, string]("name", p.Args["name"].(string))).Exec(r.ctx(p))
}

func (r *resolver) apps(p graphql.ResolveParams) (interface{}, error) {
	return r.db.App.FindMany().Exec(r.ctx(p))
}

func (r *resolver) skill(p graphql.ResolveParams) (interface{}, error) {
	if id, ok := p.Args["id"].(string); ok && id != "" {
		return r.db.Skill.FindFirst(eq[orm.Skill, string]("id", id)).Exec(r.ctx(p))
	}
	if name, ok := p.Args["name"].(string); ok && name != "" {
		return r.db.Skill.FindFirst(eq[orm.Skill, string]("name", name)).Exec(r.ctx(p))
	}
	return r.db.Skill.FindFirst().Exec(r.ctx(p))
}

func (r *resolver) skills(p graphql.ResolveParams) (interface{}, error) {
	return r.db.Skill.FindMany().Exec(r.ctx(p))
}

func (r *resolver) user(p graphql.ResolveParams) (interface{}, error) {
	return r.db.User.FindFirst(eq[orm.User, string]("id", p.Args["id"].(string))).Exec(r.ctx(p))
}

func (r *resolver) users(p graphql.ResolveParams) (interface{}, error) {
	return r.db.User.FindMany().Exec(r.ctx(p))
}

func (r *resolver) project(p graphql.ResolveParams) (interface{}, error) {
	return r.db.Project.FindFirst(eq[orm.Project, string]("name", p.Args["name"].(string))).Exec(r.ctx(p))
}

func (r *resolver) projects(p graphql.ResolveParams) (interface{}, error) {
	return r.db.Project.FindMany().Exec(r.ctx(p))
}

// Relation resolvers.

func (r *resolver) projectSkills(p graphql.ResolveParams) (interface{}, error) {
	project, ok := p.Source.(*orm.Project)
	if !ok || project == nil {
		return []*orm.Skill{}, nil
	}
	rows, err := r.db.ProjectSkill.FindMany(eq[orm.ProjectSkill, string]("projectName", project.Name)).Exec(r.ctx(p))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.SkillId)
	}
	return r.skillsByIDs(p, ids)
}

func (r *resolver) projectAgents(p graphql.ResolveParams) (interface{}, error) {
	project, ok := p.Source.(*orm.Project)
	if !ok || project == nil {
		return []*orm.Agent{}, nil
	}
	ids, err := r.joinColumnValues(p, "_ProjectAgents", "projectName", project.Name, "agentId")
	if err != nil {
		return nil, err
	}
	return r.agentsByIDs(p, ids)
}

func (r *resolver) projectApps(p graphql.ResolveParams) (interface{}, error) {
	project, ok := p.Source.(*orm.Project)
	if !ok || project == nil {
		return []*orm.App{}, nil
	}
	return r.db.App.FindMany(eq[orm.App, string]("project_id", project.Name)).Exec(r.ctx(p))
}

func (r *resolver) appSkills(p graphql.ResolveParams) (interface{}, error) {
	app, ok := p.Source.(*orm.App)
	if !ok || app == nil {
		return []*orm.Skill{}, nil
	}
	rows, err := r.db.AppSkill.FindMany(eq[orm.AppSkill, string]("app_id", app.Name)).Exec(r.ctx(p))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.SkillId)
	}
	return r.skillsByIDs(p, ids)
}

func (r *resolver) appStyles(p graphql.ResolveParams) (interface{}, error) {
	app, ok := p.Source.(*orm.App)
	if !ok || app == nil {
		return []*orm.StyleToken{}, nil
	}
	return r.db.StyleToken.FindMany(eq[orm.StyleToken, string]("app_id", app.Name)).Exec(r.ctx(p))
}

func (r *resolver) agentSkills(p graphql.ResolveParams) (interface{}, error) {
	agent, ok := p.Source.(*orm.Agent)
	if !ok || agent == nil {
		return []*orm.Skill{}, nil
	}
	rows, err := r.db.AgentSkill.FindMany(eq[orm.AgentSkill, string]("agentId", agent.Id)).Exec(r.ctx(p))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.SkillId)
	}
	return r.skillsByIDs(p, ids)
}

func (r *resolver) userPreferred(p graphql.ResolveParams) (interface{}, error) {
	user, ok := p.Source.(*orm.User)
	if !ok || user == nil {
		return []*orm.Skill{}, nil
	}
	ids, err := r.joinColumnValues(p, "_UserSkills", "userId", user.Id, "skillId")
	if err != nil {
		return nil, err
	}
	return r.skillsByIDs(p, ids)
}

func (r *resolver) userBlacklist(p graphql.ResolveParams) (interface{}, error) {
	user, ok := p.Source.(*orm.User)
	if !ok || user == nil {
		return []*orm.Skill{}, nil
	}
	ids, err := r.joinColumnValues(p, "_UserSkillBlacklist", "userId", user.Id, "skillId")
	if err != nil {
		return nil, err
	}
	return r.skillsByIDs(p, ids)
}

// helpers

// skillsByIDs loads skills by id, returning an empty slice for no ids.
func (r *resolver) skillsByIDs(p graphql.ResolveParams, ids []string) (interface{}, error) {
	if len(ids) == 0 {
		return []*orm.Skill{}, nil
	}
	return r.db.Skill.FindMany(orm.Field[orm.Skill, string]{Column: "id"}.In(ids)).Exec(r.ctx(p))
}

func (r *resolver) agentsByIDs(p graphql.ResolveParams, ids []string) (interface{}, error) {
	if len(ids) == 0 {
		return []*orm.Agent{}, nil
	}
	return r.db.Agent.FindMany(orm.Field[orm.Agent, string]{Column: "id"}.In(ids)).Exec(r.ctx(p))
}

// joinColumnValues reads an implicit many-to-many join table and returns the
// values of outColumn for rows whose inColumn equals key. The whole table is
// scanned and filtered in Go, which keeps the query dialect-neutral across
// sqlite and postgres.
func (r *resolver) joinColumnValues(p graphql.ResolveParams, table, inColumn, key, outColumn string) ([]string, error) {
	query := `SELECT "` + inColumn + `", "` + outColumn + `" FROM "` + table + `"`
	rows, err := r.db.Raw().QueryContext(r.ctx(p), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []string
	for rows.Next() {
		var in, val string
		if err := rows.Scan(&in, &val); err != nil {
			return nil, err
		}
		if in == key {
			out = append(out, val)
		}
	}
	return out, rows.Err()
}
