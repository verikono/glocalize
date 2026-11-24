package appAgent

import (
	"context"
	"github.com/verikono/glocalize/apps/riot/internal/orm"
)

type Select = riot.AppAgentSelect
type Omit = riot.AppAgentOmit
type QueryBuilder = riot.AppAgentQueryBuilder
type CreateBuilder = riot.AppAgentCreateBuilder
type Upsert = riot.AppAgentUpsert
type ConflictBuilder[B any] = riot.AppAgentConflictBuilder[B]

func Query() *QueryBuilder {
	return &QueryBuilder{}
}

func Record(assignments ...riot.FieldAssignmentOf[riot.AppAgent]) riot.RecordInput {
	raw := make([]riot.FieldAssignment, len(assignments))
	for i, a := range assignments {
		raw[i] = riot.FieldAssignment(a)
	}
	return riot.RecordInput{Assignments: raw}
}

func And(preds ...riot.PredicateOf[riot.AppAgent]) riot.PredicateOf[riot.AppAgent] {
	return riot.And(preds...)
}

func Or(preds ...riot.PredicateOf[riot.AppAgent]) riot.PredicateOf[riot.AppAgent] {
	return riot.Or(preds...)
}

func Not(pred riot.PredicateOf[riot.AppAgent]) riot.PredicateOf[riot.AppAgent] {
	return riot.Not(pred)
}

var App_language = riot.Field[riot.AppAgent, riot.AppLanguageType]{Column: "app_language"}

var App_name = riot.StringField[riot.AppAgent]{Column: "app_name"}

var AgentId = riot.StringField[riot.AppAgent]{Column: "agentId"}

type app_language_App_name_AgentId struct {
	riot.CompositeUniqueConstraint[riot.AppAgent]
	Column string
}

var App_language_App_name_AgentId = app_language_App_name_AgentId{
	CompositeUniqueConstraint: riot.CompositeUniqueConstraint[riot.AppAgent]{
		Name: "App_language_App_name_AgentId",
		Columns: []string{
			"app_language",
			"app_name",
			"agentId",
		},
	},
	Column: "App_language_App_name_AgentId",
}

func (f app_language_App_name_AgentId) EQ(app_language riot.AppLanguageType, app_name string, agentId string) riot.UniquePredicate[riot.AppAgent] {
	return riot.UniquePredicate[riot.AppAgent]{
		Data: riot.PredicateData{
			Column:   "App_language_App_name_AgentId",
			Operator: "AND",
			Value: map[string]any{
				"app_language": app_language,
				"app_name":     app_name,
				"agentId":      agentId,
			},
			IsLogical: true,
			Children: []riot.PredicateData{
				{Column: "app_language", Operator: "=", Value: app_language},
				{Column: "app_name", Operator: "=", Value: app_name},
				{Column: "agentId", Operator: "=", Value: agentId},
			},
		},
	}
}

type CreateInput = riot.AppAgentCreate
type Create = riot.AppAgentCreate

type CreateArgs = riot.AppAgentCreateArgs
type CreateManyArgs = riot.AppAgentCreateManyArgs
type CreateManyAndReturnArgs = riot.AppAgentCreateManyAndReturnArgs

type CreateQuery = riot.AppAgentCreateQuery
type CreateHook = func(context.Context, *CreateArgs, CreateQuery) (*riot.AppAgent, error)

type CreateManyQuery = riot.AppAgentCreateManyQuery
type CreateManyHook = func(context.Context, *CreateManyArgs, CreateManyQuery) (int64, error)

type CreateManyAndReturnQuery = riot.AppAgentCreateManyAndReturnQuery
type CreateManyAndReturnHook = func(context.Context, *CreateManyAndReturnArgs, CreateManyAndReturnQuery) ([]*riot.AppAgent, error)

type FindUniqueArgs = riot.AppAgentFindUniqueArgs
type FindUniqueQuery = riot.AppAgentFindUniqueQuery
type FindUniqueHook = func(context.Context, *FindUniqueArgs, FindUniqueQuery) (*riot.AppAgent, error)

type FindFirstArgs = riot.AppAgentFindFirstArgs
type FindFirstQuery = riot.AppAgentFindFirstQuery
type FindFirstHook = func(context.Context, *FindFirstArgs, FindFirstQuery) (*riot.AppAgent, error)

type FindManyArgs = riot.AppAgentFindManyArgs
type FindManyQuery = riot.AppAgentFindManyQuery
type FindManyHook = func(context.Context, *FindManyArgs, FindManyQuery) ([]*riot.AppAgent, error)

type CountArgs = riot.AppAgentCountArgs
type CountQuery = riot.AppAgentCountQuery
type CountHook = func(context.Context, *CountArgs, CountQuery) (int64, error)

// type Update = riot.AppAgentUpdate

type Extension = riot.AppAgentExtension

// ConflictUpdate creates a custom ConflictAction for AppAgent upsert conflicts in the hooks.
func ConflictUpdate(fn func(u *riot.AppAgentUpsert)) *riot.ConflictAction {
	return riot.AppAgentConflictUpdate(fn)
}
