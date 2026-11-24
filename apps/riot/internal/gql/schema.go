package gql

import (
	"context"
	"fmt"
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

// unique builds a unique-field predicate, used by Update and Delete.
func unique[M any](column string, value any) orm.UniquePredicate[M] {
	return orm.UniquePredicate[M]{
		Data: orm.PredicateData{Column: column, Operator: "=", Value: value},
	}
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
			"keywords":    &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"createdAt":   &graphql.Field{Type: dateTimeScalar},
			"updatedAt":   &graphql.Field{Type: dateTimeScalar},
		},
	})

	// Attach the self-referential dependency relations once both types exist.
	// These bridge through the SkillDependency join table and return Skills.
	skillType.AddFieldConfig("dependencies", &graphql.Field{Type: graphql.NewList(skillType), Resolve: r.skillDependencies})
	skillType.AddFieldConfig("dependents", &graphql.Field{Type: graphql.NewList(skillType), Resolve: r.skillDependents})

	// Reference is the standalone catalogue of resources a skill can draw on;
	// SkillReference is the join that attaches one to a skill. The enums are
	// exposed as strings, matching how StyleToken's enums are served.
	referenceType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Reference",
		Fields: graphql.Fields{
			"name":         &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"description":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"accessType":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"resourceType": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"content":      &graphql.Field{Type: graphql.String},
		},
	})

	skillReferenceType := graphql.NewObject(graphql.ObjectConfig{
		Name: "SkillReference",
		Fields: graphql.Fields{
			"skillId":       &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"referenceName": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"reference":     &graphql.Field{Type: referenceType, Resolve: r.skillReferenceReference},
		},
	})

	skillType.AddFieldConfig("references", &graphql.Field{Type: graphql.NewList(skillReferenceType), Resolve: r.skillReferences})

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

	// AgentSkill is the binding row, not the skill: it carries the priority the
	// agent holds the skill at, which a bare Skill record cannot express.
	agentSkillType := graphql.NewObject(graphql.ObjectConfig{
		Name: "AgentSkill",
		Fields: graphql.Fields{
			"agentId":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"skillId":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"priority": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"skill":    &graphql.Field{Type: skillType, Resolve: r.agentSkillSkill},
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
			"skills":       &graphql.Field{Type: graphql.NewList(agentSkillType), Resolve: r.agentSkills},
		},
	})

	appType := graphql.NewObject(graphql.ObjectConfig{
		Name: "App",
		Fields: graphql.Fields{
			"type":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"language":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"name":      &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"location":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"projectId": &graphql.Field{Type: graphql.String},
			"skills":    &graphql.Field{Type: graphql.NewList(skillType), Resolve: r.appSkills},
			"agents":    &graphql.Field{Type: graphql.NewList(agentType), Resolve: r.appAgents},
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
			"location":   &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
			"skillsPath": &graphql.Field{Type: graphql.String},
			"agentsPath": &graphql.Field{Type: graphql.String},
			"skills":     &graphql.Field{Type: graphql.NewList(skillType), Resolve: r.projectSkills},
			"agents":     &graphql.Field{Type: graphql.NewList(agentType), Resolve: r.projectAgents},
			"apps":       &graphql.Field{Type: graphql.NewList(appType), Resolve: r.projectApps},
		},
	})

	// projects is the reverse of Project.skills: the projects a skill is
	// assigned to, bridged through the ProjectSkill join table.
	skillType.AddFieldConfig("projects", &graphql.Field{Type: graphql.NewList(projectType), Resolve: r.skillProjects})

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
					"name":     &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"language": &graphql.ArgumentConfig{Type: graphql.String},
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
			"Reference": &graphql.Field{
				Type: referenceType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.reference,
			},
			"References": &graphql.Field{
				Type:    graphql.NewList(referenceType),
				Resolve: r.references,
			},
		},
	})

	skillInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "SkillInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"id":           &graphql.InputObjectFieldConfig{Type: graphql.String},
			"name":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"description":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"url":          &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"content":      &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"keywords":     &graphql.InputObjectFieldConfig{Type: graphql.String},
			"dependencies": &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.NewNonNull(graphql.String))},
			"references":   &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.NewNonNull(graphql.String))},
			"projects":     &graphql.InputObjectFieldConfig{Type: graphql.NewList(graphql.NewNonNull(graphql.String))},
		},
	})

	projectInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "ProjectInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":       &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"location":   &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"skillsPath": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"agentsPath": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"isNew":      &graphql.InputObjectFieldConfig{Type: graphql.Boolean},
		},
	})

	// ReferenceInput is the write shape for a reference. `name` is the unique
	// identity, so a create with an existing name is rejected rather than
	// silently overwriting the stored resource.
	referenceInputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "ReferenceInput",
		Fields: graphql.InputObjectConfigFieldMap{
			"name":         &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"description":  &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"accessType":   &graphql.InputObjectFieldConfig{Type: graphql.String},
			"resourceType": &graphql.InputObjectFieldConfig{Type: graphql.NewNonNull(graphql.String)},
			"content":      &graphql.InputObjectFieldConfig{Type: graphql.String},
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"SaveSkill": &graphql.Field{
				Type: skillType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(skillInputType)},
				},
				Resolve: r.saveSkill,
			},
			"SaveReference": &graphql.Field{
				Type: referenceType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(referenceInputType)},
				},
				Resolve: r.saveReference,
			},
			"SaveProject": &graphql.Field{
				Type: projectType,
				Args: graphql.FieldConfigArgument{
					"input": &graphql.ArgumentConfig{Type: graphql.NewNonNull(projectInputType)},
				},
				Resolve: r.saveProject,
			},
			"BindAgentSkills": &graphql.Field{
				Type: graphql.NewList(skillType),
				Args: graphql.FieldConfigArgument{
					"agentId":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"skillIds": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String)))},
					"priority": &graphql.ArgumentConfig{Type: graphql.String},
				},
				Resolve: r.bindAgentSkills,
			},
			"UnbindAgentSkill": &graphql.Field{
				Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(skillType))),
				Args: graphql.FieldConfigArgument{
					"agentId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"skillId": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
				},
				Resolve: r.unbindAgentSkill,
			},
			"AddSkillDependencies": &graphql.Field{
				Type: graphql.NewList(skillType),
				Args: graphql.FieldConfigArgument{
					"skillId":       &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.String)},
					"dependencyIds": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(graphql.String)))},
				},
				Resolve: r.addSkillDependencies,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{Query: query, Mutation: mutation})
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
	preds := []orm.PredicateOf[orm.App]{eq[orm.App, string]("name", p.Args["name"].(string))}
	if language, ok := p.Args["language"].(string); ok && language != "" {
		preds = append(preds, eq[orm.App, string]("language", language))
	}
	return r.db.App.FindFirst(preds...).Exec(r.ctx(p))
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

// saveSkill creates or updates a skill and replaces its dependency set. The
// whole write runs in one transaction so a failure part-way through cannot
// leave a skill with half of its dependencies attached.
func (r *resolver) saveSkill(p graphql.ResolveParams) (interface{}, error) {
	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("SaveSkill: input is required")
	}

	name, _ := input["name"].(string)
	description, _ := input["description"].(string)
	url, _ := input["url"].(string)
	content, _ := input["content"].(string)
	keywords, _ := input["keywords"].(string)
	id, _ := input["id"].(string)

	// a skill created from the browser form may omit its url; the catalogue's
	// convention is a sandbox url of the form skill://<project>/<name>, and
	// every skill currently lives under the "all" project.
	if url == "" && name != "" {
		url = fmt.Sprintf("skill://all/%s", name)
	}

	dependencyIDs := make([]string, 0)
	if raw, ok := input["dependencies"].([]interface{}); ok {
		for _, item := range raw {
			if dep, ok := item.(string); ok && dep != "" {
				dependencyIDs = append(dependencyIDs, dep)
			}
		}
	}

	// references REPLACES the skill's reference set, exactly like dependencies:
	// the reference records themselves are never deleted, only the join rows.
	referenceNames := make([]string, 0)
	if raw, ok := input["references"].([]interface{}); ok {
		for _, item := range raw {
			if name, ok := item.(string); ok && name != "" {
				referenceNames = append(referenceNames, name)
			}
		}
	}

	// projects REPLACES the skill's project assignments the same way: the join
	// rows are rewritten and the Project records themselves are never deleted.
	// A nil list leaves the assignments untouched, so a caller that does not
	// know about projects cannot silently unassign every one of them.
	projectNames := make([]string, 0)
	replaceProjects := false
	if raw, ok := input["projects"].([]interface{}); ok {
		replaceProjects = true
		for _, item := range raw {
			if name, ok := item.(string); ok && name != "" {
				projectNames = append(projectNames, name)
			}
		}
	}

	ctx := r.ctx(p)

	var saved *orm.Skill
	err := r.db.Transaction(ctx, func(tx *orm.Tx) error {
		var err error
		if id == "" {
			saved, err = tx.Skill.Create().
				SetName(name).
				SetDescription(description).
				SetUrl(url).
				SetContent(content).
				SetKeywords(keywords).
				Exec(ctx)
		} else {
			saved, err = tx.Skill.Update(unique[orm.Skill]("id", id)).
				SetName(name).
				SetDescription(description).
				SetUrl(url).
				SetContent(content).
				SetKeywords(keywords).
				Exec(ctx)
		}
		if err != nil {
			return err
		}

		// replace the dependency set: detach everything, then attach what was
		// asked for. A skill may not depend on itself.
		if _, err := tx.SkillDependency.DeleteMany(eq[orm.SkillDependency, string]("dependentId", saved.Id)).Exec(ctx); err != nil {
			return err
		}
		for _, dependencyID := range dependencyIDs {
			if dependencyID == saved.Id {
				continue
			}
			if _, err := tx.SkillDependency.Create().
				SetDependentId(saved.Id).
				SetDependencyId(dependencyID).
				Exec(ctx); err != nil {
				return err
			}
		}

		// replace the reference set the same way: detach every join row, then
		// attach the named references. The Reference records survive.
		if _, err := tx.SkillReference.DeleteMany(eq[orm.SkillReference, string]("skillId", saved.Id)).Exec(ctx); err != nil {
			return err
		}
		for _, referenceName := range referenceNames {
			if _, err := tx.SkillReference.Create().
				SetSkillId(saved.Id).
				SetReferenceName(referenceName).
				Exec(ctx); err != nil {
				return err
			}
		}

		// replace the project assignments only when the caller supplied a list.
		if replaceProjects {
			if _, err := tx.ProjectSkill.DeleteMany(eq[orm.ProjectSkill, string]("skillId", saved.Id)).Exec(ctx); err != nil {
				return err
			}
			for _, projectName := range projectNames {
				if _, err := tx.ProjectSkill.Create().
					SetProjectName(projectName).
					SetSkillId(saved.Id).
					Exec(ctx); err != nil {
					return err
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return saved, nil
}

func (r *resolver) user(p graphql.ResolveParams) (interface{}, error) {
	return r.db.User.FindFirst(eq[orm.User, string]("id", p.Args["id"].(string))).Exec(r.ctx(p))
}

func (r *resolver) users(p graphql.ResolveParams) (interface{}, error) {
	return r.db.User.FindMany().Exec(r.ctx(p))
}

// bindAgentSkills attaches each skill to the agent at the given priority. A
// skill already bound is re-prioritised rather than duplicated, so the call is
// idempotent. It returns the agent's full skill list afterwards.
func (r *resolver) bindAgentSkills(p graphql.ResolveParams) (interface{}, error) {
	agentID, _ := p.Args["agentId"].(string)
	priority, _ := p.Args["priority"].(string)
	if agentID == "" {
		return nil, fmt.Errorf("BindAgentSkills: agentId is required")
	}
	if priority == "" {
		priority = string(orm.Priority_NORMAL)
	}
	if !orm.PriorityType(priority).IsValid() {
		return nil, fmt.Errorf("BindAgentSkills: unknown priority %q", priority)
	}

	skillIDs := make([]string, 0)
	if raw, ok := p.Args["skillIds"].([]interface{}); ok {
		for _, item := range raw {
			if id, ok := item.(string); ok && id != "" {
				skillIDs = append(skillIDs, id)
			}
		}
	}

	ctx := r.ctx(p)
	err := r.db.Transaction(ctx, func(tx *orm.Tx) error {
		for _, skillID := range skillIDs {
			existing, err := tx.AgentSkill.FindFirst(
				eq[orm.AgentSkill, string]("agentId", agentID),
				eq[orm.AgentSkill, string]("skillId", skillID),
			).Exec(ctx)
			if err != nil {
				return err
			}
			if existing != nil {
				if _, err := tx.AgentSkill.Update(unique[orm.AgentSkill]("agentId", agentID), eq[orm.AgentSkill, string]("skillId", skillID)).
					SetPriority(orm.PriorityType(priority)).
					Exec(ctx); err != nil {
					return err
				}
				continue
			}
			if _, err := tx.AgentSkill.Create().
				SetAgentId(agentID).
				SetSkillId(skillID).
				SetPriority(orm.PriorityType(priority)).
				Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return r.skillsForAgent(ctx, agentID)
}

// unbindAgentSkill detaches one skill from one agent. A pair that is not bound is a no-op, so the call is idempotent.
func (r *resolver) unbindAgentSkill(p graphql.ResolveParams) (interface{}, error) {
	agentID, _ := p.Args["agentId"].(string)
	skillID, _ := p.Args["skillId"].(string)
	if agentID == "" {
		return nil, fmt.Errorf("UnbindAgentSkill: agentId is required")
	}
	if skillID == "" {
		return nil, fmt.Errorf("UnbindAgentSkill: skillId is required")
	}

	ctx := r.ctx(p)
	if _, err := r.db.AgentSkill.DeleteMany(eq[orm.AgentSkill, string]("agentId", agentID), eq[orm.AgentSkill, string]("skillId", skillID)).Exec(ctx); err != nil {
		return nil, err
	}
	return r.skillsForAgent(ctx, agentID)
}

func (r *resolver) reference(p graphql.ResolveParams) (interface{}, error) {
	return r.db.Reference.FindFirst(eq[orm.Reference, string]("name", p.Args["name"].(string))).Exec(r.ctx(p))
}

// addSkillDependencies attaches each dependency to the skill, leaving the
// skill's existing dependencies alone. A pair already linked is skipped, so the
// call is idempotent. It returns the skill's full dependency list afterwards.
func (r *resolver) addSkillDependencies(p graphql.ResolveParams) (interface{}, error) {
	skillID, _ := p.Args["skillId"].(string)
	if skillID == "" {
		return nil, fmt.Errorf("AddSkillDependencies: skillId is required")
	}

	dependencyIDs := make([]string, 0)
	if raw, ok := p.Args["dependencyIds"].([]interface{}); ok {
		for _, item := range raw {
			if id, ok := item.(string); ok && id != "" && id != skillID {
				dependencyIDs = append(dependencyIDs, id)
			}
		}
	}

	ctx := r.ctx(p)
	err := r.db.Transaction(ctx, func(tx *orm.Tx) error {
		for _, dependencyID := range dependencyIDs {
			existing, err := tx.SkillDependency.FindFirst(
				eq[orm.SkillDependency, string]("dependentId", skillID),
				eq[orm.SkillDependency, string]("dependencyId", dependencyID),
				eq[orm.SkillDependency, string]("agent", ""),
			).Exec(ctx)
			if err != nil {
				return err
			}
			if existing != nil {
				continue
			}
			if _, err := tx.SkillDependency.Create().
				SetDependentId(skillID).
				SetDependencyId(dependencyID).
				SetAgent("").
				Exec(ctx); err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	skill, err := r.db.Skill.FindFirst(eq[orm.Skill, string]("id", skillID)).Exec(ctx)
	if err != nil {
		return nil, err
	}
	if skill == nil {
		return []*orm.Skill{}, nil
	}
	return r.skillDependencies(graphql.ResolveParams{Source: skill, Context: ctx})
}

func (r *resolver) references(p graphql.ResolveParams) (interface{}, error) {
	return r.db.Reference.FindMany().Exec(r.ctx(p))
}

// saveProject creates a project, or updates the one named by `name`. The name is the primary key, so an update never renames it. @see .agents.md "### project identity"
func (r *resolver) saveProject(p graphql.ResolveParams) (interface{}, error) {
	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("SaveProject: input is required")
	}

	name, _ := input["name"].(string)
	location, _ := input["location"].(string)
	skillsPath, _ := input["skillsPath"].(string)
	agentsPath, _ := input["agentsPath"].(string)
	isNew, _ := input["isNew"].(bool)

	if name == "" {
		return nil, fmt.Errorf("SaveProject: name is required")
	}

	ctx := r.ctx(p)

	if isNew {
		return r.db.Project.Create().SetName(name).SetLocation(location).SetSkillsPath(skillsPath).SetAgentsPath(agentsPath).Exec(ctx)
	}

	return r.db.Project.Update(unique[orm.Project]("name", name)).SetLocation(location).SetSkillsPath(skillsPath).SetAgentsPath(agentsPath).Exec(ctx)
}

// saveReference creates a reference. The name is the unique identity, so an
// existing name is an error rather than an overwrite: the browser's create
// modal must not silently replace a resource another skill already draws on.
func (r *resolver) saveReference(p graphql.ResolveParams) (interface{}, error) {
	input, ok := p.Args["input"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("SaveReference: input is required")
	}

	name, _ := input["name"].(string)
	description, _ := input["description"].(string)
	accessType, _ := input["accessType"].(string)
	resourceType, _ := input["resourceType"].(string)
	content, _ := input["content"].(string)

	if name == "" || description == "" || resourceType == "" {
		return nil, fmt.Errorf("SaveReference: name, description and resourceType are required")
	}
	if accessType == "" {
		accessType = string(orm.URI_STATIC)
	}
	if !orm.URIType(accessType).IsValid() {
		return nil, fmt.Errorf("SaveReference: invalid accessType %q", accessType)
	}
	if !orm.ReferenceTypeType(resourceType).IsValid() {
		return nil, fmt.Errorf("SaveReference: invalid resourceType %q", resourceType)
	}

	ctx := r.ctx(p)

	existing, err := r.db.Reference.FindFirst(eq[orm.Reference, string]("name", name)).Exec(ctx)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, fmt.Errorf("SaveReference: a reference named %q already exists", name)
	}

	return r.db.Reference.Create().
		SetName(name).
		SetDescription(description).
		SetAccessType(orm.URIType(accessType)).
		SetResourceType(orm.ReferenceTypeType(resourceType)).
		SetContent(content).
		Exec(ctx)
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
	rows, err := r.db.AppSkill.FindMany(
		eq[orm.AppSkill, string]("app_language", string(app.Language)),
		eq[orm.AppSkill, string]("app_name", app.Name),
	).Exec(r.ctx(p))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.SkillId)
	}
	return r.skillsByIDs(p, ids)
}

func (r *resolver) appAgents(p graphql.ResolveParams) (interface{}, error) {
	app, ok := p.Source.(*orm.App)
	if !ok || app == nil {
		return []*orm.Agent{}, nil
	}
	rows, err := r.db.AppAgent.FindMany(
		eq[orm.AppAgent, string]("app_language", string(app.Language)),
		eq[orm.AppAgent, string]("app_name", app.Name),
	).Exec(r.ctx(p))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.AgentId)
	}
	return r.agentsByIDs(p, ids)
}

func (r *resolver) appStyles(p graphql.ResolveParams) (interface{}, error) {
	app, ok := p.Source.(*orm.App)
	if !ok || app == nil {
		return []*orm.StyleToken{}, nil
	}
	return r.db.StyleToken.FindMany(
		eq[orm.StyleToken, string]("app_language", string(app.Language)),
		eq[orm.StyleToken, string]("app_name", app.Name),
	).Exec(r.ctx(p))
}

// agentSkills returns the agent's binding rows, so each skill arrives with the
// priority the agent holds it at.
func (r *resolver) agentSkills(p graphql.ResolveParams) (interface{}, error) {
	agent, ok := p.Source.(*orm.Agent)
	if !ok || agent == nil {
		return []*orm.AgentSkill{}, nil
	}
	return r.db.AgentSkill.FindMany(eq[orm.AgentSkill, string]("agentId", agent.Id)).Exec(r.ctx(p))
}

// agentSkillSkill loads the skill a binding row points at.
func (r *resolver) agentSkillSkill(p graphql.ResolveParams) (interface{}, error) {
	row, ok := p.Source.(*orm.AgentSkill)
	if !ok || row == nil {
		return nil, nil
	}
	return r.db.Skill.FindFirst(eq[orm.Skill, string]("id", row.SkillId)).Exec(r.ctx(p))
}

// skillsForAgent loads the skills bound to an agent, in binding order.
func (r *resolver) skillsForAgent(ctx context.Context, agentID string) (interface{}, error) {
	rows, err := r.db.AgentSkill.FindMany(eq[orm.AgentSkill, string]("agentId", agentID)).Exec(ctx)
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.SkillId)
	}
	return r.skillsByIDs(graphql.ResolveParams{Context: ctx}, ids)
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

// skillDependencies returns the skills this skill depends on, bridged through
// the SkillDependency join table (rows where this skill is the dependent).
func (r *resolver) skillDependencies(p graphql.ResolveParams) (interface{}, error) {
	skill, ok := p.Source.(*orm.Skill)
	if !ok || skill == nil {
		return []*orm.Skill{}, nil
	}
	rows, err := r.db.SkillDependency.FindMany(eq[orm.SkillDependency, string]("dependentId", skill.Id)).Exec(r.ctx(p))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.DependencyId)
	}
	return r.skillsByIDs(p, ids)
}

// skillDependents returns the skills that depend on this skill, bridged through
// the SkillDependency join table (rows where this skill is the dependency).
func (r *resolver) skillDependents(p graphql.ResolveParams) (interface{}, error) {
	skill, ok := p.Source.(*orm.Skill)
	if !ok || skill == nil {
		return []*orm.Skill{}, nil
	}
	rows, err := r.db.SkillDependency.FindMany(eq[orm.SkillDependency, string]("dependencyId", skill.Id)).Exec(r.ctx(p))
	if err != nil {
		return nil, err
	}
	ids := make([]string, 0, len(rows))
	for _, row := range rows {
		ids = append(ids, row.DependentId)
	}
	return r.skillsByIDs(p, ids)
}

// skillProjects returns the projects this skill is assigned to, bridged
// through the ProjectSkill join table.
func (r *resolver) skillProjects(p graphql.ResolveParams) (interface{}, error) {
	skill, ok := p.Source.(*orm.Skill)
	if !ok || skill == nil {
		return []*orm.Project{}, nil
	}
	rows, err := r.db.ProjectSkill.FindMany(eq[orm.ProjectSkill, string]("skillId", skill.Id)).Exec(r.ctx(p))
	if err != nil {
		return nil, err
	}
	names := make([]string, 0, len(rows))
	for _, row := range rows {
		names = append(names, row.ProjectName)
	}
	if len(names) == 0 {
		return []*orm.Project{}, nil
	}
	return r.db.Project.FindMany(orm.Field[orm.Project, string]{Column: "name"}.In(names)).Exec(r.ctx(p))
}

// skillReferences returns the join rows attaching references to this skill.
// The reference records themselves are loaded by the SkillReference.reference
// field resolver, so a skill's references are always the current catalogue
// entries rather than a copy taken when the join was written.
func (r *resolver) skillReferences(p graphql.ResolveParams) (interface{}, error) {
	skill, ok := p.Source.(*orm.Skill)
	if !ok || skill == nil {
		return []*orm.SkillReference{}, nil
	}
	return r.db.SkillReference.FindMany(eq[orm.SkillReference, string]("skillId", skill.Id)).Exec(r.ctx(p))
}

// skillReferenceReference loads the reference a join row points at.
func (r *resolver) skillReferenceReference(p graphql.ResolveParams) (interface{}, error) {
	row, ok := p.Source.(*orm.SkillReference)
	if !ok || row == nil {
		return nil, nil
	}
	return r.db.Reference.FindFirst(eq[orm.Reference, string]("name", row.ReferenceName)).Exec(r.ctx(p))
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
